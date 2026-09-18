import { Agent, fetch as undiciFetch } from "undici";
import {
  type CallRequest,
  type CallResponse,
  IngextHttpError,
  IngextRpcError,
} from "./errors.js";

export interface IngextClientOptions {
  url: string;
  token?: string;
  debug?: boolean;
  /**
   * Target a tenant account through a provider (grid manager) proxy. When set,
   * "?gridaccount=<account>" is appended to every `api/ds` request URL, exactly
   * as `IngextClient.SetGridAccount` does in the Go client. Client-wide, not
   * per call.
   *
   * Note that failures differ through the proxy: a tenant "ERROR" verdict
   * arrives as "EXCEPTION", with the message text preserved.
   */
  gridAccount?: string;
  /**
   * If true, TLS certificate verification is disabled. Mirrors the Go
   * client's `InsecureSkipVerify: true` default. Off by default in the TS
   * port — set to true when talking to development clusters with self-signed
   * certs.
   */
  insecure?: boolean;
  /** Per-request timeout in milliseconds. Default 600_000 (10 minutes). */
  timeoutMs?: number;
  /**
   * Optional injected fetch implementation (primarily for tests). When
   * unset, undici's fetch is used so the `insecure` flag can take effect.
   */
  fetch?: typeof globalThis.fetch;
  /** Logger for debug output. Defaults to `console`. */
  logger?: Pick<Console, "debug" | "error">;
}

export class IngextClient {
  private url: string;
  private token: string;
  private debugFlag: boolean;
  private timeoutMs: number;
  private logger: Pick<Console, "debug" | "error">;
  private fetchImpl: typeof globalThis.fetch;
  private dispatcher?: Agent;
  private gridAccount: string;

  constructor(opts: IngextClientOptions) {
    this.url = opts.url.replace(/\/+$/, "");
    this.token = opts.token ?? "";
    this.debugFlag = opts.debug ?? false;
    this.timeoutMs = opts.timeoutMs ?? 600_000;
    this.logger = opts.logger ?? console;
    this.gridAccount = opts.gridAccount ?? "";

    if (opts.fetch) {
      this.fetchImpl = opts.fetch;
    } else {
      if (opts.insecure) {
        this.dispatcher = new Agent({
          connect: { rejectUnauthorized: false },
        });
      }
      const dispatcher = this.dispatcher;
      this.fetchImpl = ((input, init) => {
        const merged = dispatcher
          ? { ...(init as Record<string, unknown>), dispatcher }
          : init;
        return undiciFetch(input as Parameters<typeof undiciFetch>[0], merged as never) as unknown as ReturnType<typeof globalThis.fetch>;
      }) as typeof globalThis.fetch;
    }
  }

  setToken(token: string): void {
    this.token = token;
  }

  setDebug(flag: boolean): void {
    this.debugFlag = flag;
  }

  /**
   * Route every subsequent call to a tenant account through the provider proxy.
   * Pass "" to go back to the site's own account.
   */
  setGridAccount(account: string): void {
    this.gridAccount = account;
  }

  getGridAccount(): string {
    return this.gridAccount;
  }

  getUrl(): string {
    return this.url;
  }

  /**
   * Low-level RPC call. Mirrors `IngextClient.GenericCall` in the Go client:
   * POSTs `{function, kargs}` to `<url>/<prefix>/<functionName>`, decodes the
   * `CallResponse` envelope, and throws `IngextRpcError` on non-OK verdicts.
   */
  async call<R = unknown>(
    prefix: string,
    functionName: string,
    kargs: unknown,
  ): Promise<R> {
    const body: CallRequest = {
      function: functionName,
      kargs: kargs ?? {},
    };

    let fullUrl = prefix
      ? `${this.url}/${prefix}/${functionName}`
      : `${this.url}/${functionName}`;

    if (this.gridAccount) {
      fullUrl = `${fullUrl}?gridaccount=${encodeURIComponent(this.gridAccount)}`;
    }

    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };
    if (this.token) {
      headers["Authorization"] = `Bearer ${this.token}`;
    }

    if (this.debugFlag) {
      this.logger.debug(`[ingext] POST ${fullUrl}`);
      this.logger.debug(`[ingext] body: ${JSON.stringify(body)}`);
    }

    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort(), this.timeoutMs);

    let resp: Response;
    try {
      resp = await this.fetchImpl(fullUrl, {
        method: "POST",
        headers,
        body: JSON.stringify(body),
        signal: controller.signal,
      });
    } finally {
      clearTimeout(timer);
    }

    if (!resp.ok) {
      throw new IngextHttpError(fullUrl, resp.status, resp.statusText);
    }

    const payload = (await resp.json()) as CallResponse<R>;

    if (this.debugFlag) {
      this.logger.debug(`[ingext] response: ${JSON.stringify(payload)}`);
    }

    if (payload.verdict !== "OK") {
      throw new IngextRpcError(
        prefix,
        functionName,
        payload.verdict,
        payload.error,
        payload.exception,
      );
    }

    return payload.response as R;
  }

  async close(): Promise<void> {
    if (this.dispatcher) {
      await this.dispatcher.close();
    }
  }
}

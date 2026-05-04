export type Verdict = "OK" | "ERROR" | "EXCEPTION";

export interface Attachment {
  name?: string;
  path?: string;
  id?: string;
  type?: string;
  hash?: string;
}

export interface CallRequest<T = unknown> {
  function: string;
  kargs: T;
  attachments?: Attachment[];
}

export interface CallResponse<T = unknown> {
  verdict: Verdict;
  response?: T;
  attachments?: Attachment[];
  error?: string;
  exception?: string;
}

export class IngextHttpError extends Error {
  readonly status: number;
  readonly statusText: string;
  readonly url: string;

  constructor(url: string, status: number, statusText: string) {
    super(`HTTP error from ${url}: ${status} ${statusText}`);
    this.name = "IngextHttpError";
    this.url = url;
    this.status = status;
    this.statusText = statusText;
  }
}

export class IngextRpcError extends Error {
  readonly verdict: Verdict;
  readonly prefix: string;
  readonly functionName: string;
  readonly rpcError?: string;
  readonly rpcException?: string;

  constructor(
    prefix: string,
    functionName: string,
    verdict: Verdict,
    rpcError?: string,
    rpcException?: string,
  ) {
    const detail = verdict === "EXCEPTION" ? rpcException : rpcError;
    super(`RPC ${verdict} from ${prefix}/${functionName}: ${detail ?? "(no detail)"}`);
    this.name = "IngextRpcError";
    this.verdict = verdict;
    this.prefix = prefix;
    this.functionName = functionName;
    this.rpcError = rpcError;
    this.rpcException = rpcException;
  }
}

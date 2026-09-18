import { describe, expect, it, vi } from "vitest";
import { IngextClient } from "../src/client.js";
import { IngextHttpError, IngextRpcError } from "../src/errors.js";

function mockFetch(impl: (url: string, init: RequestInit) => Response | Promise<Response>) {
  return vi.fn(async (input: string | URL | Request, init?: RequestInit) => {
    const url = typeof input === "string" ? input : input.toString();
    return impl(url, init ?? {});
  }) as unknown as typeof globalThis.fetch;
}

describe("IngextClient.call", () => {
  it("posts CallRequest envelope to <url>/<prefix>/<function> with bearer token", async () => {
    const seen: { url: string; init: RequestInit } = { url: "", init: {} };
    const fetchImpl = mockFetch((url, init) => {
      seen.url = url;
      seen.init = init;
      return new Response(JSON.stringify({ verdict: "OK", response: { ok: 1 } }), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    });

    const client = new IngextClient({
      url: "https://ingext.example.com",
      token: "tok123",
      fetch: fetchImpl,
    });

    const res = await client.call<{ ok: number }>("api/auth", "userList", { foo: "bar" });

    expect(res).toEqual({ ok: 1 });
    expect(seen.url).toBe("https://ingext.example.com/api/auth/userList");
    expect(seen.init.method).toBe("POST");
    const headers = seen.init.headers as Record<string, string>;
    expect(headers["Content-Type"]).toBe("application/json");
    expect(headers["Authorization"]).toBe("Bearer tok123");
    expect(JSON.parse(seen.init.body as string)).toEqual({
      function: "userList",
      kargs: { foo: "bar" },
    });
  });

  it("appends ?gridaccount= when a grid account is set", async () => {
    let url = "";
    const fetchImpl = mockFetch((u) => {
      url = u;
      return new Response(JSON.stringify({ verdict: "OK", response: {} }), { status: 200 });
    });

    const client = new IngextClient({
      url: "https://develop.app.ingext.io",
      token: "tok",
      gridAccount: "jet",
      fetch: fetchImpl,
    });

    await client.call("api/ds", "platform_list_configs", {});
    expect(url).toBe("https://develop.app.ingext.io/api/ds/platform_list_configs?gridaccount=jet");
  });

  it("setGridAccount switches and clears the tenant", async () => {
    const urls: string[] = [];
    const fetchImpl = mockFetch((u) => {
      urls.push(u);
      return new Response(JSON.stringify({ verdict: "OK", response: {} }), { status: 200 });
    });

    const client = new IngextClient({ url: "https://develop.app.ingext.io", fetch: fetchImpl });

    await client.call("api/ds", "platform_list_configs", {});
    client.setGridAccount("titan");
    expect(client.getGridAccount()).toBe("titan");
    await client.call("api/ds", "platform_list_configs", {});
    client.setGridAccount("");
    await client.call("api/ds", "platform_list_configs", {});

    expect(urls[0]).toBe("https://develop.app.ingext.io/api/ds/platform_list_configs");
    expect(urls[1]).toBe("https://develop.app.ingext.io/api/ds/platform_list_configs?gridaccount=titan");
    expect(urls[2]).toBe("https://develop.app.ingext.io/api/ds/platform_list_configs");
  });

  it("url-encodes a grid account name", async () => {
    let url = "";
    const fetchImpl = mockFetch((u) => {
      url = u;
      return new Response(JSON.stringify({ verdict: "OK", response: {} }), { status: 200 });
    });

    const client = new IngextClient({
      url: "https://develop.app.ingext.io",
      gridAccount: "acc name/1",
      fetch: fetchImpl,
    });

    await client.call("api/ds", "platform_list_configs", {});
    expect(url).toBe("https://develop.app.ingext.io/api/ds/platform_list_configs?gridaccount=acc%20name%2F1");
  });

  it("strips trailing slashes from base url", async () => {
    let url = "";
    const fetchImpl = mockFetch((u) => {
      url = u;
      return new Response(JSON.stringify({ verdict: "OK", response: {} }), { status: 200 });
    });

    const client = new IngextClient({ url: "https://x.example/", fetch: fetchImpl });
    await client.call("api/ds", "ping", null);
    expect(url).toBe("https://x.example/api/ds/ping");
  });

  it("uses empty kargs object when input is null/undefined", async () => {
    let body = "";
    const fetchImpl = mockFetch((_, init) => {
      body = init.body as string;
      return new Response(JSON.stringify({ verdict: "OK", response: {} }), { status: 200 });
    });
    const client = new IngextClient({ url: "https://x.example", fetch: fetchImpl });
    await client.call("api/ds", "ping", null);
    expect(JSON.parse(body)).toEqual({ function: "ping", kargs: {} });
  });

  it("throws IngextRpcError on verdict ERROR", async () => {
    const fetchImpl = mockFetch(() =>
      new Response(JSON.stringify({ verdict: "ERROR", error: "no such user" }), { status: 200 }),
    );
    const client = new IngextClient({ url: "https://x.example", fetch: fetchImpl });
    await expect(client.call("api/auth", "getUser", { username: "x" })).rejects.toMatchObject({
      name: "IngextRpcError",
      verdict: "ERROR",
      rpcError: "no such user",
      prefix: "api/auth",
      functionName: "getUser",
    });
  });

  it("throws IngextRpcError on verdict EXCEPTION", async () => {
    const fetchImpl = mockFetch(() =>
      new Response(JSON.stringify({ verdict: "EXCEPTION", exception: "panic in handler" }), {
        status: 200,
      }),
    );
    const client = new IngextClient({ url: "https://x.example", fetch: fetchImpl });
    await expect(client.call("api/ds", "boom", null)).rejects.toBeInstanceOf(IngextRpcError);
  });

  it("throws IngextHttpError on non-200 response", async () => {
    const fetchImpl = mockFetch(
      () => new Response("nope", { status: 502, statusText: "Bad Gateway" }),
    );
    const client = new IngextClient({ url: "https://x.example", fetch: fetchImpl });
    await expect(client.call("api/ds", "x", null)).rejects.toMatchObject({
      name: "IngextHttpError",
      status: 502,
    });
    await expect(client.call("api/ds", "x", null)).rejects.toBeInstanceOf(IngextHttpError);
  });

  it("calls function without prefix when prefix is empty", async () => {
    let url = "";
    const fetchImpl = mockFetch((u) => {
      url = u;
      return new Response(JSON.stringify({ verdict: "OK", response: {} }), { status: 200 });
    });
    const client = new IngextClient({ url: "https://x.example", fetch: fetchImpl });
    await client.call("", "rootFn", null);
    expect(url).toBe("https://x.example/rootFn");
  });

  it("omits Authorization when no token is set", async () => {
    let headers: Record<string, string> = {};
    const fetchImpl = mockFetch((_, init) => {
      headers = init.headers as Record<string, string>;
      return new Response(JSON.stringify({ verdict: "OK", response: {} }), { status: 200 });
    });
    const client = new IngextClient({ url: "https://x.example", fetch: fetchImpl });
    await client.call("api/ds", "x", null);
    expect(headers["Authorization"]).toBeUndefined();
  });
});

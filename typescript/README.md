# ingext-api

TypeScript client for the Ingext platform API. Mirrors the public `client/` +
`api/` packages of the Go module `github.com/SecurityDo/ingext_api`.

- ESM only, requires Node ≥ 18 (uses the global `fetch` via `undici`).
- Bearer-token auth — bring your own token (no login flow).
- One typed service per backend area (auth, datalake, search, eventwatch, fpl,
  syslog, grid, collector, resource, notification, application, repo, platform).

## Install

```bash
npm install ingext-api
```

## Quick start

```ts
import { Ingext } from "ingext-api";

const ingext = new Ingext({
  url: "https://my-ingext.example.com",
  token: process.env.INGEXT_TOKEN!,
  // insecure: true,   // skip TLS verification (dev clusters with self-signed certs)
  // debug: true,      // log every request/response to console
});

const users = await ingext.auth.listUser();
const lakes = await ingext.datalake.listDatalake();
const result = await ingext.search.kqlSearch("MyTable | where status == 200 | take 10");
```

## Constructor options

| Option       | Default     | Description                                                              |
| ------------ | ----------- | ------------------------------------------------------------------------ |
| `url`        | _required_  | Base URL of the Ingext site, e.g. `https://demo.example.com`.            |
| `token`      | `""`        | Bearer token. Obtain via the Ingext UI / `auth add-token` Go CLI.        |
| `debug`      | `false`     | Log request/response bodies to `console.debug`.                          |
| `insecure`   | `false`     | Disable TLS cert verification (dev only).                                |
| `timeoutMs`  | `600000`    | Per-request timeout (default matches the Go client's 600s).              |
| `fetch`      | undici      | Inject a custom fetch (used by tests).                                   |
| `logger`     | `console`   | Replace the debug logger.                                                |

## One example per service area

```ts
// auth
await ingext.auth.addUser({ user: { username: "a@x.com", roles: ["admin"], oauthProvider: "", oauthFlag: false } });
await ingext.auth.addToken("ci-bot", "CI bot", "analyst");

// datalake / schema
await ingext.datalake.addDatalake("", true, "");
await ingext.datalake.addDatalakeIndex("managed", "events", "ingext default");
await ingext.datalake.addSchema("my-schema", "...", JSON.stringify(schema));

// search
const r = await ingext.search.kqlSearch("MyTable | take 10");
const v = await ingext.search.kqlValidate("MyTable | where x == 1");

// eventwatch
const now = Date.now();
const summary  = await ingext.eventwatch.summarySearch("",  now - 3_600_000, now);
const timeline = await ingext.eventwatch.timelineSearch("", now - 3_600_000, now);
const rules    = await ingext.eventwatch.ruleSearch("");

// eventwatch rules (eventwatch_bucket_dao)
const ruleList = await ingext.eventwatch.listRule();
const rule     = await ingext.eventwatch.getRule("my_rule");
await ingext.eventwatch.addRule(rule!);
await ingext.eventwatch.updateRule(rule!);
await ingext.eventwatch.toggleRule("my_rule");
await ingext.eventwatch.deleteRule("my_rule");

// fpl
const taskID = await ingext.fpl.runReport({ reportName: "my-report" });
const status = await ingext.fpl.getTaskByID(taskID);
const data   = await ingext.fpl.getResultsByID(taskID);

// syslog
const cfg = await ingext.syslog.get();
await ingext.syslog.register(["tcp", "udp"]);

// grid (SaaS accounts)
const accounts = await ingext.grid.listAccount();

// collector
const collectors = await ingext.collector.collectorList();

// resource
const office = await ingext.resource.search("office365User", "_all_");

// notification
const epID = await ingext.notification.addEmail("ops", "alert", ["a@x"], []);

// application
const tplID = await ingext.application.addAppTemplate(yamlContent);
await ingext.application.installAppTemplate({ config: { application: "tpl", instance: "i1" } });

// repo
const repos = await ingext.repo.listRepos();
await ingext.repo.importRepoProcessors(repos[0]!.id, ["proc.fpl"]);

// platform (sources, sinks, routers, processors, integrations, …)
const sources = await ingext.platform.listDataSource();
const role = await ingext.platform.getPodRole();
```

## Errors

Calls throw on non-OK responses:

- `IngextHttpError` for non-200 HTTP responses.
- `IngextRpcError` when `verdict` is `ERROR` or `EXCEPTION`. Carries
  `verdict`, `prefix`, `functionName`, `rpcError`, `rpcException`.

```ts
import { IngextRpcError } from "ingext-api";

try {
  await ingext.auth.getUser("nope@x.com");
} catch (e) {
  if (e instanceof IngextRpcError && e.verdict === "ERROR") {
    console.error("rpc error:", e.rpcError);
  }
}
```

## Low-level escape hatch

For RPC functions not yet wrapped by a typed service:

```ts
import { IngextClient } from "ingext-api";

const c = new IngextClient({ url, token });
const result = await c.call<MyResponse>("api/ds", "my_function", { x: 1 });
```

## Wire format (for cross-checking with the Go client)

Every call is `POST <baseUrl>/<prefix>/<function>` with body:

```json
{ "function": "<function>", "kargs": { ... } }
```

The server responds with:

```json
{ "verdict": "OK" | "ERROR" | "EXCEPTION", "response": { ... }, "error": "...", "exception": "..." }
```

This matches the Go `fsb.CallRequest` / `fsb.CallResponse` envelope exactly.

## Caveats

- `int64` values larger than `2^53` lose precision when JSON-decoded into
  TypeScript `number`. The Go server emits these as JSON numbers, so this
  caveat applies to both clients in practice.
- `time.Time` values arrive as ISO-8601 strings (Go's default JSON encoding).
  The TS types reflect this — they are typed as `string`, not `Date`.
- The package is Node-only today (`undici` for `insecure: true` support).
  Browser support would be a small `package.json` tweak — code uses only
  standards-based `fetch`.

import { describe, expect, it, vi } from "vitest";
import { Ingext } from "../src/index.js";

interface Capture {
  url: string;
  body: unknown;
}

function makeIngext(handler: (capture: Capture) => unknown) {
  const cap: Capture = { url: "", body: null };
  const fetchImpl = vi.fn(async (input: string | URL | Request, init?: RequestInit) => {
    cap.url = typeof input === "string" ? input : input.toString();
    cap.body = JSON.parse((init?.body as string) ?? "{}");
    const response = handler(cap);
    return new Response(JSON.stringify({ verdict: "OK", response }), {
      status: 200,
      headers: { "Content-Type": "application/json" },
    });
  }) as unknown as typeof globalThis.fetch;

  const ingext = new Ingext({ url: "https://x.example", token: "t", fetch: fetchImpl });
  return { ingext, cap };
}

describe("AuthService", () => {
  it("listUser hits api/auth/userList and unwraps users", async () => {
    const { ingext, cap } = makeIngext(() => ({
      users: [{ username: "a", roles: ["admin"], oauthProvider: "", oauthFlag: false }],
    }));
    const users = await ingext.auth.listUser();
    expect(cap.url).toBe("https://x.example/api/auth/userList");
    expect(users).toHaveLength(1);
    expect(users[0]?.username).toBe("a");
  });

  it("addToken sends action=add and unwraps the token", async () => {
    const { ingext, cap } = makeIngext(() => ({ token: "tok-xyz" }));
    const tok = await ingext.auth.addToken("ci", "ci bot", "analyst");
    expect(cap.url).toBe("https://x.example/api/auth/api_token");
    expect(cap.body).toMatchObject({
      function: "api_token",
      kargs: {
        action: "add",
        args: { entry: { name: "ci", description: "ci bot", roles: ["analyst"] } },
      },
    });
    expect(tok).toBe("tok-xyz");
  });
});

describe("DatalakeService", () => {
  it("listDatalake calls ingext_datalake_dao with action=list", async () => {
    const { ingext, cap } = makeIngext(() => ({
      entries: [{ name: "managed", managed: true }],
    }));
    const lakes = await ingext.datalake.listDatalake();
    expect(cap.url).toBe("https://x.example/api/ds/ingext_datalake_dao");
    expect(cap.body).toMatchObject({ kargs: { action: "list" } });
    expect(lakes[0]?.name).toBe("managed");
  });

  it("addDatalake forces name='managed' when managed=true", async () => {
    const { ingext, cap } = makeIngext(() => ({}));
    await ingext.datalake.addDatalake("ignored", true, "int-1");
    expect(cap.body).toMatchObject({
      kargs: {
        action: "add",
        args: { entry: { name: "managed", managed: true, integrationID: "int-1" } },
      },
    });
  });
});

describe("SearchService", () => {
  it("kqlSearch posts kql_search with the kql payload", async () => {
    const { ingext, cap } = makeIngext(() => ({ total: 0, data: { Tables: [] } }));
    await ingext.search.kqlSearch("MyTable | take 10");
    expect(cap.url).toBe("https://x.example/api/ds/kql_search");
    expect(cap.body).toMatchObject({
      function: "kql_search",
      kargs: { kql: "MyTable | take 10" },
    });
  });
});

describe("SyslogService", () => {
  it("register translates port strings into the SyslogPortRequest shape", async () => {
    const { ingext, cap } = makeIngext(() => ({ config: { domain: "x" } }));
    await ingext.syslog.register(["udp", "tls-rfc6587"]);
    expect(cap.body).toMatchObject({
      kargs: {
        syslog_udp: true,
        syslog_tcp: false,
        syslog_tls: false,
        tls_rfc6587: true,
      },
    });
  });
});

describe("FPLService", () => {
  it("runReport returns the taskID number", async () => {
    const { ingext, cap } = makeIngext(() => ({ taskID: 42 }));
    const id = await ingext.fpl.runReport({ reportName: "foo" });
    expect(cap.url).toBe("https://x.example/api/ds/run_fplv2_report");
    expect(id).toBe(42);
  });
});

describe("GridService", () => {
  it("listAccount uses the api/grid prefix", async () => {
    const { ingext, cap } = makeIngext(() => ({ accounts: [], entries: [] }));
    await ingext.grid.listAccount();
    expect(cap.url).toBe("https://x.example/api/grid/get_grid_accounts");
  });
});

describe("CollectorService", () => {
  it("collectorList returns the entries array", async () => {
    const { ingext, cap } = makeIngext(() => ({
      entries: [{ id: "c1", name: "c1", local: false, ebpf: false, description: "", createdOn: 0, token: "", lastPoll: 0 }],
    }));
    const collectors = await ingext.collector.collectorList();
    expect(cap.url).toBe("https://x.example/api/ds/collector_list");
    expect(collectors).toHaveLength(1);
  });
});

describe("NotificationService", () => {
  it("addEmail calls platform_notification_endpoint_dao with action=add", async () => {
    const { ingext, cap } = makeIngext(() => ({ id: "ep-1" }));
    const id = await ingext.notification.addEmail("ops", "alert", ["a@x"], ["b@x"]);
    expect(id).toBe("ep-1");
    expect(cap.body).toMatchObject({
      function: "platform_notification_endpoint_dao",
      kargs: {
        action: "add",
        args: {
          entry: {
            name: "ops",
            integration: "Email",
            action: "alert",
            email: { to: ["a@x"], cc: ["b@x"] },
          },
        },
      },
    });
  });
});

describe("RepoService", () => {
  it("importRepoProcessors posts the paths array", async () => {
    const { ingext, cap } = makeIngext(() => ({}));
    await ingext.repo.importRepoProcessors("repo-1", ["a.fpl", "b.fpl"]);
    expect(cap.url).toBe("https://x.example/api/ds/platform_import_processors");
    expect(cap.body).toMatchObject({
      kargs: { id: "repo-1", paths: ["a.fpl", "b.fpl"] },
    });
  });
});

describe("PlatformService", () => {
  it("getPodRole returns the role/arn pair", async () => {
    const { ingext, cap } = makeIngext(() => ({ role: "ingest-role", arn: "arn:aws:iam::1:role/x" }));
    const r = await ingext.platform.getPodRole();
    expect(cap.url).toBe("https://x.example/api/ds/get_pod_role");
    expect(r).toEqual({ role: "ingest-role", arn: "arn:aws:iam::1:role/x" });
  });

  it("listAssumedRole calls platform_instancerole_dao with action=list", async () => {
    const { ingext, cap } = makeIngext(() => ({ entries: [] }));
    await ingext.platform.listAssumedRole();
    expect(cap.body).toMatchObject({
      function: "platform_instancerole_dao",
      kargs: { action: "list" },
    });
  });

  it("addSimpleRouter falls back router=processor when routerName is empty", async () => {
    const { ingext, cap } = makeIngext(() => ({ id: "r-1" }));
    const id = await ingext.platform.addSimpleRouter("proc-1", "");
    expect(id).toBe("r-1");
    expect(cap.body).toMatchObject({
      kargs: { processor: "proc-1", router: "proc-1", pipe: "proc-1" },
    });
  });

  it("updatePipeProcessor passes the names through unchanged", async () => {
    const { ingext, cap } = makeIngext(() => ({}));
    await ingext.platform.updatePipeProcessor({
      routerName: "main",
      pipeName: "p1",
      processorName: "filter",
    });
    expect(cap.body).toMatchObject({
      kargs: { routerName: "main", pipeName: "p1", processorName: "filter" },
    });
  });
});

describe("EventWatchService", () => {
  it("ruleSearch posts to eventwatch_bucket_search with default options", async () => {
    const { ingext, cap } = makeIngext(() => ({}));
    await ingext.eventwatch.ruleSearch("attack");
    expect(cap.url).toBe("https://x.example/api/ds/eventwatch_bucket_search");
    expect(cap.body).toMatchObject({
      kargs: {
        options: {
          searchStr: "attack",
          fetchLimit: 100,
          sortField: "name",
          sortOrder: "asc",
        },
      },
    });
  });
});

describe("ResourceService", () => {
  it("search posts to resource_search with the resource type and customer", async () => {
    const { ingext, cap } = makeIngext(() => ({ took: 1, query: {} }));
    await ingext.resource.search("office365User", "_all_");
    expect(cap.url).toBe("https://x.example/api/ds/resource_search");
    expect(cap.body).toMatchObject({
      kargs: { resource: "office365User", customer: "_all_" },
    });
  });
});

describe("ApplicationService", () => {
  it("addAppTemplate returns the new id", async () => {
    const { ingext, cap } = makeIngext(() => ({ id: "tpl-1" }));
    const id = await ingext.application.addAppTemplate("apiVersion: ...");
    expect(id).toBe("tpl-1");
    expect(cap.body).toMatchObject({
      function: "platform_application_template_dao",
      kargs: { action: "add" },
    });
  });
});

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

  it("lakeSearch posts lake_search with the console's kargs shape", async () => {
    const { ingext, cap } = makeIngext(() => ({ took: 12, total: 4, filtered: 172 }));
    const resp = await ingext.search.lakeSearch({
      index: "managed-Office365",
      query: "71.178.173.2",
      rangeFrom: 1787766069700,
      rangeTo: 1787769669701,
      facets: [{ title: "Username", field: "@fields.UserId" }, { field: "@fields.ClientIP", size: 5 }],
      mustFilters: [{ field: "@source", terms: ["Audit.AzureActiveDirectory"] }],
      fetchLimit: 1000,
    });
    expect(cap.url).toBe("https://x.example/api/ds/lake_search");
    expect(cap.body).toMatchObject({
      function: "lake_search",
      kargs: {
        index: "managed-Office365",
        options: {
          searchStr: "71.178.173.2",
          range_from: 1787766069700,
          range_to: 1787769669701,
          fetchLimit: 1000,
          fetchOffset: 0,
          // No server-side default: an empty sort field fails the query parser.
          sortField: "@timestamp",
          sortOrder: "desc",
          facets: {
            // The platform keys the histogram by the date facet's name.
            dateFacets: [{ name: "dateHistogram" }],
            // A facet of size 0 comes back with no buckets at all.
            facets: [
              { title: "Username", field: "@fields.UserId", size: 20 },
              { title: "@fields.ClientIP", field: "@fields.ClientIP", size: 5 },
            ],
            mustFilters: [{ field: "@source", terms: ["Audit.AzureActiveDirectory"] }],
            mustNotFilters: [],
          },
        },
      },
    });
    // dataType, partition and dayIndex are declared by the endpoint but never
    // read, so they stay off the wire.
    expect(Object.keys((cap.body as { kargs: object }).kargs).sort()).toEqual(["index", "options"]);
    expect(resp.filtered).toBe(172);
  });

  it("lakeSearch always sends the facets object: the endpoint walks it without a nil check", async () => {
    const { ingext, cap } = makeIngext(() => ({}));
    await ingext.search.lakeSearch({ index: "managed-Office365", rangeFrom: 1, rangeTo: 2 });
    const options = (cap.body as { kargs: { options: { facets: unknown; fetchLimit: number } } })
      .kargs.options;
    expect(options.facets).toEqual({
      dateFacets: [{ name: "dateHistogram" }],
      facets: [],
      mustFilters: [],
      mustNotFilters: [],
    });
    // A fetchLimit of 0 reaches the platform and panics its search node.
    expect(options.fetchLimit).toBe(100);
  });

  it("lakeSearch rejects the requests that fail unhelpfully server-side", async () => {
    const { ingext } = makeIngext(() => ({}));
    await expect(ingext.search.lakeSearch({ rangeFrom: 0, rangeTo: 0 })).rejects.toThrow(
      /epoch milliseconds/,
    );
    await expect(ingext.search.lakeSearch({ rangeFrom: 200, rangeTo: 100 })).rejects.toThrow(
      /range is empty/,
    );
    await expect(
      ingext.search.lakeSearch({ rangeFrom: 1, rangeTo: 2, fetchOffset: 4000, fetchLimit: 1001 }),
    ).rejects.toThrow(/pagination limit is 5000/);
    await expect(
      ingext.search.lakeSearch({ rangeFrom: 1, rangeTo: 2, fetchOffset: 4950 }),
    ).rejects.toThrow(/pagination limit is 5000/);
    await expect(
      ingext.search.lakeSearch({ rangeFrom: 1, rangeTo: 2, facets: [{ title: "Username" }] }),
    ).rejects.toThrow(/facet 0: field is required/);
    await expect(
      ingext.search.lakeSearch({ rangeFrom: 1, rangeTo: 2, mustFilters: [{ field: "@source", terms: [] }] }),
    ).rejects.toThrow(/at least one term/);
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

  it("listRule calls eventwatch_bucket_dao with action=list and unwraps entries", async () => {
    const { ingext, cap } = makeIngext(() => ({
      entries: [{ name: "rule-a", group: "g1", disabled: false }],
      tags: ["t1"],
      groups: ["g1"],
    }));
    const res = await ingext.eventwatch.listRule();
    expect(cap.url).toBe("https://x.example/api/ds/eventwatch_bucket_dao");
    expect(cap.body).toMatchObject({ kargs: { action: "list" } });
    expect(res.entries[0]?.name).toBe("rule-a");
    expect(res.tags).toEqual(["t1"]);
    expect(res.groups).toEqual(["g1"]);
  });

  it("getRule calls action=get with the id and unwraps entry", async () => {
    const { ingext, cap } = makeIngext(() => ({ entry: { name: "rule-a", group: "g1" } }));
    const entry = await ingext.eventwatch.getRule("rule-a");
    expect(cap.url).toBe("https://x.example/api/ds/eventwatch_bucket_dao");
    expect(cap.body).toMatchObject({ kargs: { action: "get", args: { id: "rule-a" } } });
    expect(entry?.name).toBe("rule-a");
  });

  it("getRule returns null when no entry is present", async () => {
    const { ingext } = makeIngext(() => ({ entry: null }));
    const entry = await ingext.eventwatch.getRule("missing");
    expect(entry).toBeNull();
  });

  it("addRule calls action=add with the entry", async () => {
    const { ingext, cap } = makeIngext(() => ({}));
    await ingext.eventwatch.addRule({ name: "rule-a", group: "g1" } as never);
    expect(cap.body).toMatchObject({
      kargs: { action: "add", args: { entry: { name: "rule-a", group: "g1" } } },
    });
  });

  it("updateRule calls action=update with the entry", async () => {
    const { ingext, cap } = makeIngext(() => ({}));
    await ingext.eventwatch.updateRule({ name: "rule-a", group: "g1" } as never);
    expect(cap.body).toMatchObject({
      kargs: { action: "update", args: { entry: { name: "rule-a" } } },
    });
  });

  it("toggleRule calls action=toggle with the id", async () => {
    const { ingext, cap } = makeIngext(() => ({}));
    await ingext.eventwatch.toggleRule("rule-a");
    expect(cap.body).toMatchObject({ kargs: { action: "toggle", args: { id: "rule-a" } } });
  });

  it("deleteRule calls action=delete with the id", async () => {
    const { ingext, cap } = makeIngext(() => ({}));
    await ingext.eventwatch.deleteRule("rule-a");
    expect(cap.body).toMatchObject({ kargs: { action: "delete", args: { id: "rule-a" } } });
  });

  it("listBehaviorFilter calls behavior_filter_dao with action=list and no args", async () => {
    const { ingext, cap } = makeIngext(() => ({
      entries: [{ behaviorRule: "*", name: "scanner" }],
    }));
    const entries = await ingext.eventwatch.listBehaviorFilter();
    expect(cap.url).toBe("https://x.example/api/ds/behavior_filter_dao");
    expect(cap.body).toMatchObject({ kargs: { action: "list" } });
    expect((cap.body as { kargs: { args?: unknown } }).kargs.args).toBeUndefined();
    expect(entries).toHaveLength(1);
    expect(entries[0]?.name).toBe("scanner");
  });

  it("getBehaviorFilter joins the rule and the name into the id", async () => {
    const { ingext, cap } = makeIngext(() => ({ entry: { behaviorRule: "r1", name: "f1" } }));
    const entry = await ingext.eventwatch.getBehaviorFilter("r1", "f1");
    expect(cap.body).toMatchObject({ kargs: { action: "get", args: { id: "r1/f1" } } });
    expect(entry?.name).toBe("f1");
  });

  it("rejects behavior filter keys the server cannot split", async () => {
    const { ingext } = makeIngext(() => ({}));
    await expect(ingext.eventwatch.getBehaviorFilter("", "f1")).rejects.toThrow(
      /behaviorRule is required/,
    );
    await expect(ingext.eventwatch.getBehaviorFilter("r1", "")).rejects.toThrow(
      /name is required/,
    );
    await expect(ingext.eventwatch.deleteBehaviorFilter("a/b", "f1")).rejects.toThrow(
      /must not contain a slash/,
    );
  });

  it("getBehaviorFilter returns null when the entry is absent", async () => {
    const { ingext } = makeIngext(() => ({}));
    const entry = await ingext.eventwatch.getBehaviorFilter("r1", "missing");
    expect(entry).toBeNull();
  });

  it("addBehaviorFilter identifies the filter through the entry alone", async () => {
    const { ingext, cap } = makeIngext(() => ({}));
    await ingext.eventwatch.addBehaviorFilter({
      behaviorRule: "r1",
      name: "f1",
      action: "discard",
    } as never);
    expect(cap.body).toMatchObject({
      kargs: { action: "add", args: { entry: { behaviorRule: "r1", name: "f1" } } },
    });
    const args = (cap.body as { kargs: { args: Record<string, unknown> } }).kargs.args;
    expect(args.id).toBeUndefined();
  });

  it("updateBehaviorFilter sends action=update with the entry", async () => {
    const { ingext, cap } = makeIngext(() => ({}));
    await ingext.eventwatch.updateBehaviorFilter({ behaviorRule: "r1", name: "f1" } as never);
    expect(cap.body).toMatchObject({
      kargs: { action: "update", args: { entry: { behaviorRule: "r1", name: "f1" } } },
    });
  });

  it("toggleBehaviorFilter and deleteBehaviorFilter key on the composite id", async () => {
    const toggle = makeIngext(() => ({}));
    await toggle.ingext.eventwatch.toggleBehaviorFilter("*", "scanner");
    expect(toggle.cap.body).toMatchObject({
      kargs: { action: "toggle", args: { id: "*/scanner" } },
    });

    const del = makeIngext(() => ({}));
    await del.ingext.eventwatch.deleteBehaviorFilter("*", "scanner");
    expect(del.cap.body).toMatchObject({
      kargs: { action: "delete", args: { id: "*/scanner" } },
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

  it("dumpDelete posts the customer and returns the purged count", async () => {
    const { ingext, cap } = makeIngext(() => ({ deleted: 3 }));
    const deleted = await ingext.resource.dumpDelete("office365");
    expect(cap.url).toBe("https://x.example/api/ds/ingext_resource_dump_delete");
    expect(cap.body).toMatchObject({
      function: "ingext_resource_dump_delete",
      kargs: { customer: "office365" },
    });
    expect(deleted).toBe(3);
  });

  it("dumpDelete reports an unknown customer as 0, not an error", async () => {
    const { ingext } = makeIngext(() => ({ deleted: 0 }));
    await expect(ingext.resource.dumpDelete("never-collected")).resolves.toBe(0);
  });

  it("dumpDelete rejects an empty customer without a request", async () => {
    const { ingext, cap } = makeIngext(() => {
      throw new Error("should not be called");
    });
    await expect(ingext.resource.dumpDelete("")).rejects.toThrow(/customer is required/);
    expect(cap.url).toBe("");
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

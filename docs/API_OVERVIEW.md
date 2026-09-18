= Ingext API overview

* All APIs are HTTP POST => https://$domain/prefix/$function
**prefixs:  "api/auth" (user/role related), "api/grid" (grid management), or "api/ds" (other apis)
**The function call request format: fsb.CallRequest 
[source, golang]
type CallRequest struct {
   Function    string        `json:"function"`
   Kargs       json.RawMessage        `json:"kargs"`
}
** The response format: fsb.CallResponse
***if the verdict is "OK", the response field contains the result.
***if the verdict is "ERROR", the error field show the error information.
[source, golang]
type CallResponse struct {
   Verdict string `json:"verdict"`
   Response    json.RawMessage        `json:"response"`
   Error       string        `json:"error,omitempty"`
}
* Any 'api/ds' function can be routed to a tenant through a grid manager by adding a "gridaccount" query parameter: https://$grid/api/ds/$function?gridaccount=$account
** In the Go client this is client-wide, not per-call: IngextClient.SetGridAccount("$account").
** Success payloads are identical whether the call went direct or through the grid.
** Failures are not: a tenant "ERROR" verdict arrives as "EXCEPTION" through the grid. The message text is preserved.
* For typical CRUD operations, the dao APIs (with a suffix "_dao") has the following format. 
** action values: get, add, delete, update, list, toggle
[source, golang] 
type GenericDaoRequestArgs[T any] struct {
  Id    string `json:"id,omitempty"`
  Entry *T      `json:"entry,omitempty"`
  Flag  bool   `json:"flag,omitempty"`  // optional for toggle action
}
type GenericDaoRequest[T any] struct {
	Action string                    `json:"action"`
	Args   *GenericDaoRequestArgs[T] `json:"args,omitempty"`
}

== authorization related APIs: (with prefix 'api/auth')

* userAdd
* userList
* getUser
* userDelete
* api_token (dao-style: action = add | delete | list)
* setUserSitePolicy (kargs: username, policyName)

== grid APIs: (with prefix 'api/grid')

* get_grid_accounts — list all SaaS accounts
* add_saas_account (kargs: name, region, cluster, siteURL, token, displayName)
* delete_saas_account (kargs: name)

== collector APIs: (with prefix 'api/ds')

* collector_list — returns entries of CollectorForWeb
* collector_import (kargs: collector) — import a CollectorT, keeping its token, to migrate a collector between sites
* get_system_status (kargs: collector, cargs) — get collector status

== eventwatch / overview APIs: (with prefix 'api/ds')

* behavior_summary_search (kargs: options with SimpleSearchOption) — summary search
* fsm_behavior_search (kargs: options with SimpleSearchOption) — timeline search
* eventwatch_bucket_search (kargs: options with SimpleSearchOption) — rule search
* eventwatch_bucket_delete_group (kargs: group) — delete all rules in a group
* eventwatch_bucket_dao (dao-style: action = list | get | add | update | delete | toggle) — EventWatchBucket CRUD, keyed by rule name
* eventwatch_rule_test (kargs: bucket, input) — run one rule against one event, without deploying the rule or storing what it produces
** bucket is a whole EventWatchBucket, deployed or not, and input is the event as a JSON object. Neither is validated: a missing bucket, or one whose name is empty, is answered with a panic ("runtime error: invalid memory address or nil pointer dereference") rather than an error, and an input that is not an object (a JSON string holding an event, an array) is accepted and then ignored, leaving hit false with nothing to show for it.
** The response is {hit, input | output, signals, behaviorEvent}. The event is echoed back under input, or under output on a hit that produced a behavior event, where it has been null in every response seen so far.
** hit only means the event selector matched. A rule with no fields or behavior rule configured hits with signals and behaviorEvent null, and a rule whose isProcessor is set hits the same way.
** signals are JSON documents carried as strings: {signal, ts, count, key, valueMap}, where ts is in seconds while behaviorEvent.timestamp is in milliseconds.
* behavior_filter_dao (dao-style: action = list | get | add | update | delete | toggle) — BehaviorEventFilterT CRUD
** A behavior filter is keyed by (behaviorRule, name), not by name alone: together they form the etcd key acc_$account/fsm/filters/$behaviorRule/$name. The pair travels in the single "id" arg joined by a slash — "$behaviorRule/$name".
** The server splits id on the *first* slash, so a filter name may contain slashes but a behavior rule may not. An id with no slash at all is rejected as "invalid filter name".
** get / delete / toggle take id. add / update take entry, which carries its own name and behaviorRule, and ignore id; list takes no args and returns every filter of the account.
** A behaviorRule of "*" applies the filter to every behavior rule.
** get reports a missing filter as an error ("behavior filter not found"), not as an empty entry.
** add and update both require a non-empty filters list — a filter with no entries is rejected even when matchAll is set. Regex values must compile, or the entry is rejected as invalid.
** update matches on the entry's (behaviorRule, name); changing either field does not rename the filter, it fails as "behavior filter does not exist".

== FPL report APIs: (with prefix 'api/ds')

* run_fplv2_report (kargs: RunFPLV2Report) — submit an FPL v2 report, returns taskID
* get_fplv2_task (kargs: id, fpl) — get task status by ID
* get_fplv2_result (kargs: id, fpl) — get task results by ID

== platform topology APIs: (with prefix 'api/ds')

* platform_list_configs (no kargs) — the whole topology of the account in one call: sources, sinks, routers, pipes, channels, connections, integrations, errors and errorStates.
** This is the *only* inventory call for routers: platform_router_dao has no "list" action and answers one with "unknown router action:list".
* platform_router_dao (dao-style: action = get | add | delete) — RouterConfig CRUD by id. "get" returns the router and its pipes.
* platform_add_simple_router (kargs: processor, router, pipe) — create a router with one pipe carrying one processor. The pipe is created with priority 0, no tags and no sinks.
* platform_router_add_pipe (kargs: routerID, pipeConfig) — add a pipe to an existing router, leaving its other pipes alone. This is how a second consumer (a behavior pipe beside an application's main pipe) is attached without reinstalling the application template.
** pipeConfig is a StreamPipeConfig: {name, routerID, matchAll, selector, processorNames, sinkIDs, priority, tags}.
** **processorNames is a list, but a pipe carries exactly one processor.** The extra entries are not a supported chain; build a chain as a second pipe that the first hands on to by returning "abort".
** **priority and tags are stored but were missing from both client structs** until 2026-09. A pipe written by a client that does not know about them comes back priority 0 and untagged, which is not what any app-installed pipe looks like: the templates give a main pipe 1000 and a behavior pipe 2000, plus "application"/"appInstance" tags. Read the pipe back after writing it.
** The processor named in processorNames is not validated at write time. A pipe naming a processor that was never deployed is accepted and fails when events reach it.
* platform_router_delete_pipe (kargs: routerID, pipeID) — remove one pipe. The router's other pipes keep running, which makes this the rollback for platform_router_add_pipe.
* platform_router_update_pipes (kargs: routerID, pipeIDs) — replace the set of pipes on a router.
* platform_pipe_update (kargs: pipe) / platform_pipe_update_processor (kargs: routerName, pipeName, processorName) — edit a pipe, or just swap its processor by name.
* platform_event_tail (kargs: id, status, limit) — the most recent events **one pipe** ended with, per status. id is a pipe id ("pipe_xxxx"), status is one of pass | abort | drop | error, limit defaults small.
** This is the live-debugging call. It answers the two questions metrics cannot: is this pipe seeing traffic at all, and what did it do with it — without waiting for the next vendor event.
** Entries are the whole event as the pipe saw it, as JSON strings, newest first. Each status has its own buffer, so a pipe that is busy on "drop" and empty on "pass" tells you the processor is running and rejecting, not that the pipe is starved.
** Reading the status of two adjacent pipes is how a multi-pipe router is verified: the upstream pipe should show its events under "abort" (handed on) and the downstream pipe the same events under whatever it returned.
* platform_processor_tail (kargs: pipeID, processorName, workerIndex, limit) — the trace log lines one processor emitted, for printf/console.log debugging.
* platform_datasink_dao (dao-style: action = get | add | delete | list) — DataSinkConfig CRUD.
** A behavior-signal sink is type "redis" with redis.redis.queue = "queue:BehaviorSummary:EventQueue" and no datalake block; the datalake writer is the same type with queue "queue:LVDBService:JobQueue" plus redis.datalake / redis.datalakeIndex.
* platform_datasource_dao (dao-style) — DataSourceConfig CRUD. platform_source_set_router (kargs: dataSourceID, routerID) connects one to a router.
* platform_source_reload (kargs: dataSourceID) — restart one data source: the manager stops whatever it is running and starts it again.
** For a PLUGIN source this is what re-forks the plugin, and so what makes a newly published binary take effect: the fork re-resolves the pinned tag's digest and re-downloads when it has changed. Publishing to the registry is not a deploy on its own — without a reload the source keeps running the binary it started with until platform-0 restarts.
** The name is platform_source_reload. platform_datasource_reload answers "unknown function".
** Siblings on the same request shape: platform_source_enable / platform_source_disable, platform_source_enable_job / platform_source_disable_job, platform_source_set_config, platform_source_set_receiver.
* platform_processor_dao (dao-style: action = get | add | update | delete | list) — FPLScript CRUD.
** "get" reports a name it does not hold as an **error** ("processor not found: <name>"), not as an empty entry.

== resource APIs: (with prefix 'api/ds')

* resource_search (kargs: resource, customer, options) — search resources by type
* ingext_resource_dump_delete (kargs: customer) — purge every resource dump of one customer, across all resource types, and return {deleted} — the number of (resourceType, customer) dump folders removed
** This is the same cleanup that deleting the owning plugin data source performs; the endpoint exists for the dumps of data sources removed while that purge silently matched nothing.
** The customer is matched literally — "_all_" is a customer name here, not a wildcard. For a plugin data source the customer segment is the plugin name the dump was written under.
** An unknown customer is not an error: it deletes nothing and returns 0, which is the only way to tell a purge from a no-op.
** Requires data/manage, not data/write — it destroys collected data with no undo. A partial failure returns ERROR with no count even though some dumps were removed; re-running it is safe.

== datalake APIs: (with prefix 'api/ds')

* list_data_tables — list every queryable table; returns streamTables (datalake indexes; the name is the KQL table identifier, query with kql_search) and resourceTables (vendor entity tables, query with resource_search). Takes no kargs. The "default" and "AzureAudit" indexes are excluded.
* ingext_datalake_dao (dao-style: action = list | add)
* ingext_datalake_index_list (kargs: lake) / ingext_datalake_index_add (kargs: entry) / ingext_datalake_index_delete (kargs: lake, index)
* ingext_datalake_schema_dao (dao-style: action = list | add | update | delete)

== search APIs: (with prefix 'api/ds')

* kql_search (kargs: kql, index, rangeFrom, rangeTo) — run a KQL query against the datalake
* kql_validate (kargs: kql) — parse a KQL query without executing it
* lake_search (kargs: index, options) — Elastic-style facet search over one datalake index: a Lucene query over a time range, must / must-not term filters, and a term count per facet field
** index is the "<datalake>-<index>" name, e.g. "managed-Office365" for the "Office365" stream table of list_data_tables. An empty index searches "default". The endpoint also declares dataType, partition and dayIndex but never reads them.
** options: searchStr (lucene, empty matches everything), range_from / range_to (epoch ms), fetchLimit / fetchOffset, sortField / sortOrder, facets{facets, mustFilters, mustNotFilters, dateFacets}
** options.facets and its four arrays must always be sent: the endpoint walks them without a nil check. The filters live inside facets — the mustFilters at the top of options are not read. Terms of one filter are OR'ed, separate filters are AND'ed.
** Two fields have no server-side default and fail unhelpfully when omitted: sortField (an empty one fails the query parser — use "@timestamp") and a facet's size (0 returns no buckets at all rather than all of them — use e.g. 20). The facet title is ignored: the response keys aggregations by field.
** fetchOffset+fetchLimit may not exceed 5000. There is no facet-only search: fetchLimit 0 does not mean "counts without hits", it panics the search node ("index out of range [0] with length 0") because the top-hits heap is built with capacity 0 and then indexed.
** The response is a LakeSearchResponse: hits, aggregations, total (matched), filtered (read and discarded), took, cost. aggregations is keyed by facet field, and is not guaranteed to hold exactly the requested set — a console response for this endpoint carried facets no one asked for — plus the date histogram, a fixed set of slots spanning the range, empty ones included. Facet buckets key on a string, histogram buckets on an epoch-millisecond number.
** The histogram aggregation is keyed by the name of the first entry in options.facets.dateFacets, so send [{"name":"dateHistogram"}] to get the key the console uses — with no date facet the histogram comes back under the empty string "".

== SentinelOne APIs: (with prefix 'api/ds')

* investigate_sentinelone_alert (kargs: alertId, integrationName, options) — investigate one SentinelOne Unified Alert end to end: the alert, the threats correlated to it, the endpoint activity around the detection, and the affected endpoint's agent record
** alertId is the alert's own id, not the threat id. integrationName is optional when exactly one SentinelOneEvents integration is configured (the type is "SentinelOneEvents", not "SentinelOne").
** options: includeThreats (default true), includeActivities (default true), includeEndpoint (default true), includeRelatedAlerts (default false), activityWindowMinutes (half-width centered on the detection time, default 60), maxActivities (default 200). Omit a flag to keep its default; an explicit false is honored. Out-of-range integers are clamped, not rejected.
** Only the alert lookup is fatal — a successful response can still report failed or partial sections in collectionStatus. Read associations.confidence before trusting the correlated threats.

== syslog APIs: (with prefix 'api/ds')

* ingext_syslog_get_config — get current syslog config
* ingext_syslog_register_config (kargs: syslogUDP, syslogTCP, syslogTLS, tlsRfc6587)
* ingext_syslog_update_config (kargs: syslogUDP, syslogTCP, syslogTLS, tlsRfc6587)
* ingext_syslog_delete_config — delete syslog config

== notification APIs: (with prefix 'api/ds')

* platform_notification_endpoint_dao (dao-style: action = list | get | add | update | delete)
** An endpoint pairs an integration ("Email", "Slack") with the name of an fpl_action script in its action field, plus the per-integration config that script reads as its `config` argument. Use platform_list_actions to find the action names, and match the action's integration to the endpoint's.
** Endpoints are keyed by **name**: get, update and delete all take the name in args.id. "add" answers "{}" with **no id field**, so there is no id to keep.
** list — list all notification endpoints (EndpointConfig), under "entries"
** get — one endpoint by name (kargs: args.id), under "entry". A name the dao does not hold is an **error**, "export not found: <name>", not an empty result.
** add — add an endpoint (kargs: args.entry = {name, integration, action, email{to, cc} | slack{channel, channels, integrationName}}). Not an upsert: a name already stored is refused with "duplicate endpoint".
** update — replace an endpoint (kargs: args.id = name, args.entry). Two traps: it takes the whole entry, not a patch, so a field left out is stored empty (read the endpoint with get and modify that); and it **is an upsert** — a name the dao does not hold is created rather than refused, so a typo in the name silently adds a second endpoint instead of editing the one meant.
** delete — delete a notification endpoint (kargs: args.id = name). Deleting a name the dao does not hold fails with "unknown export".
** Every action but list needs args: an unknown action value fails with "kargs.args field missing" rather than naming the action.

* platform_list_actions (no kargs) — the FPL action scripts the platform knows about, under "actions"
** Entries are FPLScript (the platform_processor_dao type) of type "fpl_action", with the actionConfig field that dao never fills in: actionConfig.target is the subsystem that invokes the action and actionConfig.integration the endpoint kind it expects.
** The ones a notification endpoint can name are those with target "Platform Notification"; filter further on integration to match the endpoint. The endpoint's name says nothing about the target, so filter rather than treat everything it returns as a notification action.
** The order actions come back in is not stable between calls.
** Each entry carries its whole scriptText, so the response is large relative to what a name-and-integration lookup needs.

### Billing usage (metering)

The account measures its own usage once a day from its cluster's VictoriaMetrics
and records it in the shared Postgres. Quantities and evidence only: ingext
stores no prices, no SKUs and nothing from Stripe.

* metering_daily_list (kargs: from, to, tenantKey?, includeOpen?) — the daily ledger for [from, to] inclusive, YYYY-MM-DD UTC
** A meter that was not measured is **null, and null is not zero**. Treating it as 0 turns "this tenant could not be measured" into "this tenant used nothing".
** An elapsed day with no row is returned with state "missing" rather than omitted, so a short array cannot be mistaken for a short month. state is closed | open | missing | in_progress.
** storeAvailable false means the ledger was unreachable; days will be empty and that emptiness says NOTHING about usage. Never sum such a response.
** Only state closed with status final is billable. status incomplete is written when a day aged out of the metrics retention window without resolving — it records a permanent gap and is never billable.
** Eight meters: eventwatch_bytes, processed_bytes, deleted_bytes, input_bytes, platform_datalake_bytes, lake_ingress_bytes, lake_search_bytes, lake_realtime_search_bytes.
** platform_datalake_bytes and lake_ingress_bytes measure different things and WILL disagree — the first is what the datalake sink emitted, the second what the lake actually ingested including data arriving by other paths (measured ~12x apart on a live tenant). Neither is a broken copy of the other.
** paidUsers carries one entry per integrated provider. A provider that is not integrated produces NO entry, which is different from an entry of 0.
* metering_attempts (kargs: from, to, tenantKey?, limit?) — every collection try, including failures and skips
** This is what separates "this tenant-day has no data" from "nobody ever looked". Grep errorCode to find every tenant-day sharing one cause.
* metering_collect_day (kargs: date, tenantKey?) — collect and close one past UTC day on demand
** Cannot rewrite a day that already closed; the ledger refuses it. closed false is a normal outcome meaning a meter did not resolve, and nothing was written.
* metering_sink_classification (kargs: none) — how the collector currently sorts datasinks into meters
** processed_bytes is the one meter whose value depends on a judgement about the topology rather than a metric label, so check this first when that number looks wrong.
** ambiguous lists sinks whose meter differs between the platform's resident and job execution modes. They are counted as eventwatch, which under-counts by at most that sink rather than double-billing it.
* grid_metering_daily_list (api/grid; kargs: from, to, accounts?, tenantKey?, includeOpen?) — the ledger for every tenant of a provider, in one call
** Provider-level: issue against a provider site, never a tenant. Scoped by the caller's allowed sites, and each per-tenant call is independently policy-checked.
** accounts can only narrow the caller's scope. Anything out of scope comes back in outOfScope rather than being silently dropped.
** A tenant that could not be reached appears with an error and no days. summary.failed non-zero means the document is INCOMPLETE and any provider total from it is a lower bound.

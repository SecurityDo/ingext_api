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

* platform_notification_endpoint_dao (dao-style: action = list | add | delete)
** list — list all notification endpoints (EndpointConfig)
** add — add an email notification endpoint (kargs: name, integration, action, email{to, cc})
** delete — delete a notification endpoint (kargs: id)


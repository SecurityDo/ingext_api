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
* get_system_status (kargs: collector, cargs) — get collector status

== eventwatch / overview APIs: (with prefix 'api/ds')

* behavior_summary_search (kargs: options with SimpleSearchOption) — summary search
* fsm_behavior_search (kargs: options with SimpleSearchOption) — timeline search
* eventwatch_bucket_search (kargs: options with SimpleSearchOption) — rule search

== FPL report APIs: (with prefix 'api/ds')

* run_fplv2_report (kargs: RunFPLV2Report) — submit an FPL v2 report, returns taskID
* get_fplv2_task (kargs: id, fpl) — get task status by ID
* get_fplv2_result (kargs: id, fpl) — get task results by ID

== resource APIs: (with prefix 'api/ds')

* resource_search (kargs: resource, customer, options) — search resources by type

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


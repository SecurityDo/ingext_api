// Core
export { IngextClient, type IngextClientOptions } from "./client.js";
export {
  IngextHttpError,
  IngextRpcError,
  type Attachment,
  type CallRequest,
  type CallResponse,
  type Verdict,
} from "./errors.js";
export { Ingext } from "./ingext.js";

// Services
export { ApplicationService } from "./services/application.js";
export { AuthService } from "./services/auth.js";
export { CollectorService } from "./services/collector.js";
export { DatalakeService } from "./services/datalake.js";
export { EventWatchService, type EventWatchRuleListResponse } from "./services/eventwatch.js";
export { FPLService } from "./services/fpl.js";
export { GridService } from "./services/grid.js";
export { NotificationService } from "./services/notification.js";
export { PlatformService } from "./services/platform.js";
export { RepoService } from "./services/repo.js";
export { ResourceService } from "./services/resource.js";
export { SearchService } from "./services/search.js";
export { SyslogService, type SyslogPort } from "./services/syslog.js";

// Service-level request / response types
export type {
  GetAppInstanceRequest,
  GetAppInstanceResponse,
  InstallAppInstanceRequest,
  ListAppInstanceResponse,
  ListAppTemplateResponse,
  UnInstallAppInstanceRequest,
} from "./services/application.js";
export type { AddUserRequest } from "./services/auth.js";
export type {
  GitRepoContentRequest,
  GitRepoContentResponse,
  ImportRepoObjectsRequest,
  RepositoryContent,
} from "./services/repo.js";
export type {
  AddChannelResponse,
  AddDataSinkResponse,
  AddDataSourceResponse,
  AddIntegrationResponse,
  AddRouterResponse,
  ComponentMetricReq,
  EventTailReq,
  EventTailResponse,
  FPLProcessorTestRequest,
  FPLProcessorTestResult,
  FPLProcessorValidateRequest,
  FPLProcessorValidateResult,
  GetComponentStateResponse,
  ListComponentErrorResponse,
  ListConfigsResponse,
  PipeInfo,
  PipeProcessorUpdateReq,
  PipeUpdateReq,
  PlatformMetric,
  PlatformMetricsResponse,
  PluginTailReq,
  PluginTailResponse,
  PodRoleResponse,
  ProcessorMetricReq,
  ProcessorPipes,
  ProcessorPipesReq,
  ProcessorPipesResponse,
  ProcessorTailReq,
  ProcessorTailResponse,
  RouterAddPipeReq,
  RouterAddPipeResponse,
  RouterDeletePipeReq,
  RouterEntryResponse,
  RouterUpdatePipesReq,
  SetComponentTagsReq,
  SourceSetRouterReq,
} from "./services/platform.js";

// Domain model types
export * from "./types/application.js";
export * from "./types/auth.js";
export * from "./types/collector.js";
export * from "./types/dao.js";
export * from "./types/datalake.js";
export * from "./types/elasticsearch.js";
export * from "./types/eventwatch.js";
export * from "./types/fpl.js";
export * from "./types/grid.js";
export * from "./types/integration.js";
export * from "./types/kql.js";
export * from "./types/notification.js";
export * from "./types/platform.js";
export * from "./types/pod_role.js";
export * from "./types/repo.js";
export * from "./types/search.js";
export * from "./types/site_credentials.js";
export * from "./types/syslog.js";

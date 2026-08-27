import type { IngextClient } from "../client.js";
import type {
  GenericDAORequest,
  GenericDaoAddResponse,
} from "../types/dao.js";
import type {
  ChannelConfig,
  ComponentErrorState,
  ComponentInfo,
  DataSinkConfig,
  DataSourceConfig,
  FPLScript,
  LogEvent,
  PluginNotification,
  RouterConfig,
  RouterInput,
  StreamPipeConfig,
} from "../types/platform.js";
import type {
  Integration,
  Tag,
} from "../types/integration.js";
import type { InstanceRole } from "../types/pod_role.js";

const DS = "api/ds";

export interface ListConfigsResponse {
  sources: DataSourceConfig[];
  sinks: DataSinkConfig[];
  routers: RouterConfig[];
  pipes: StreamPipeConfig[];
  channels: ChannelConfig[];
  connections: RouterInput[];
  integrations: Integration[];
  errors: PluginNotification[];
  errorStates: ComponentErrorState[];
}

export interface AddDataSourceResponse {
  id: string;
  secret?: unknown;
  url?: string;
}

export interface AddDataSinkResponse {
  id: string;
}

export interface AddRouterResponse {
  id: string;
}

export interface AddChannelResponse {
  id: string;
}

export interface AddIntegrationResponse {
  id: string;
}

export interface RouterEntryResponse {
  entry: RouterConfig | null;
  pipes: StreamPipeConfig[];
}

export interface SourceSetRouterReq {
  routerID: string;
  dataSourceID: string;
}

export interface RouterAddPipeReq {
  routerID: string;
  pipeConfig: StreamPipeConfig;
}

export interface RouterAddPipeResponse {
  id: string;
}

export interface RouterDeletePipeReq {
  routerID: string;
  pipeID: string;
}

export interface RouterUpdatePipesReq {
  routerID: string;
  pipeIDs: string[];
}

export interface PipeUpdateReq {
  routerID: string;
  pipeConfig: StreamPipeConfig;
}

export interface PipeProcessorUpdateReq {
  routerName: string;
  pipeName: string;
  processorName: string;
}

export interface FPLProcessorValidateRequest {
  name: string;
  script: string;
}

export interface FPLProcessorValidateResult {
  console: string;
  error: string;
  ok: boolean;
}

/**
 * Runs one script against one sample document without deploying it. `type`
 * selects the runtime and decides how `source` is read: an `fpl_processor`
 * takes the JSON document envelope of {@link FPLTestDocument} and returns the
 * mutated envelope, while an `fpl_receiver` takes its raw payload and returns
 * `{"docs": [...]}`. An empty or unrecognized `type` falls back to the
 * processor runtime. `name` is not looked up -- the script under test is the
 * one in `script`, deployed or not.
 */
export interface FPLProcessorTestRequest {
  name?: string;
  script: string;
  source: string;
  type?: string;
  tenant?: string;
}

/**
 * The document envelope an `fpl_processor` script is handed as its
 * `main({obj, size})` argument. platform_processor_test carries it as a JSON
 * *string* in the request's `source` field, not as a nested object, and returns
 * the mutated envelope the same way in `FPLProcessorTestResult.newContent`.
 */
export interface FPLTestDocument {
  obj: unknown;
  props: Record<string, unknown>;
  size: number;
  source: string;
}

/**
 * Wraps one event object in the test envelope and renders it as the JSON string
 * the `source` field carries. `size` is the compact byte length of the object,
 * standing in for the size the ingest path would have measured on the wire.
 */
export function buildFplTestSource(obj: unknown): string {
  if (obj === null || typeof obj !== "object" || Array.isArray(obj)) {
    throw new Error(
      "event object must be a JSON object, since the script destructures it as main({obj, size})",
    );
  }
  const encoded = JSON.stringify(obj);
  const doc: FPLTestDocument = {
    obj,
    props: {},
    size: new TextEncoder().encode(encoded).length,
    source: "",
  };
  return JSON.stringify(doc);
}

export interface FPLProcessorTestResult {
  console: string;
  error: string;
  newContent: string;
  status: string;
}

export interface EventTailReq {
  id: string;
  status?: string;
  limit?: number;
}

export interface EventTailResponse {
  entries: string[];
}

export interface ProcessorTailReq {
  pipeID: string;
  processorName?: string;
  workerIndex?: number;
  limit?: number;
}

export interface ProcessorTailResponse {
  entries: LogEvent[];
}

export interface ProcessorPipesReq {
  processorName?: string;
  processorNames?: string[];
}

export interface PipeInfo {
  pipeID: string;
  pipeName: string;
  router: string;
  workerCount: number;
}

export interface ProcessorPipes {
  processorName: string;
  pipes: PipeInfo[];
}

export interface ProcessorPipesResponse {
  pipes?: PipeInfo[];
  entries?: ProcessorPipes[];
}

export interface ComponentMetricReq {
  component: string;
  id: string;
  from: string;
  to: string;
  interval: string;
}

export interface ProcessorMetricReq {
  channel?: string;
  pipeID?: string;
  processor: string;
  from: string;
  to: string;
  interval: string;
}

export interface PlatformMetric {
  tenant?: string;
  name?: string;
  id?: string;
  pipe?: string;
  processor?: string;
  unit?: string;
  slots: number[];
  values: number[];
}

export interface PlatformMetricsResponse {
  metrics: PlatformMetric[];
}

export interface GetComponentStateResponse {
  state: unknown;
}

export interface SetComponentTagsReq {
  id: string;
  tags: Tag[];
}

export interface PluginTailReq {
  id: string;
  limit?: number;
}

export interface PluginTailResponse {
  lines: string[];
}

export interface ListComponentErrorResponse {
  errors: PluginNotification[];
}

export interface PodRoleResponse {
  role: string;
  arn: string;
}

interface IntegrationDAORequest {
  action: string;
  args?: { id?: string; entry?: Integration };
}

interface IntegrationDAOResponse {
  entry?: Integration;
  entries?: Integration[];
}

export class PlatformService {
  constructor(private client: IngextClient) {}

  // ---- Pod identity / assumed roles ----

  async getPodRole(): Promise<PodRoleResponse> {
    return await this.client.call<PodRoleResponse>(DS, "get_pod_role", null);
  }

  async testAssumedRole(roleARN: string, externalID: string): Promise<void> {
    await this.client.call(DS, "platform_instancerole_test", {
      names: [],
      roleARN,
      externalID,
    });
  }

  async addLocalAssumedRole(
    name: string,
    roleARN: string,
    externalID: string,
  ): Promise<string> {
    const req: GenericDAORequest<InstanceRole> = {
      action: "add",
      args: {
        entry: {
          id: "",
          displayName: name,
          description: "",
          externalID,
          roleARN,
          createdOn: "",
          local: true,
        },
      },
    };
    const res = await this.client.call<GenericDaoAddResponse>(
      DS,
      "platform_instancerole_add_local",
      req,
    );
    return res.id;
  }

  async addAssumedRole(
    name: string,
    roleARN: string,
    externalID: string,
  ): Promise<string> {
    const req: GenericDAORequest<InstanceRole> = {
      action: "add",
      args: {
        entry: {
          id: "",
          displayName: name,
          description: "",
          externalID,
          roleARN,
          createdOn: "",
        },
      },
    };
    const res = await this.client.call<GenericDaoAddResponse>(
      DS,
      "platform_instancerole_dao",
      req,
    );
    return res.id;
  }

  async deleteAssumedRole(roleID: string): Promise<void> {
    const req: GenericDAORequest<InstanceRole> = {
      action: "delete",
      args: { id: roleID },
    };
    await this.client.call(DS, "platform_instancerole_dao", req);
  }

  async listAssumedRole(): Promise<InstanceRole[]> {
    const req: GenericDAORequest<InstanceRole> = { action: "list" };
    const res = await this.client.call<{ entries: InstanceRole[] }>(
      DS,
      "platform_instancerole_dao",
      req,
    );
    return res.entries ?? [];
  }

  // ---- Platform info ----

  async listPlugins(): Promise<string[]> {
    const res = await this.client.call<{ plugins: string[] }>(
      DS,
      "platform_list_plugins",
      null,
    );
    return res.plugins ?? [];
  }

  async listConfigs(): Promise<ListConfigsResponse> {
    return await this.client.call<ListConfigsResponse>(
      DS,
      "platform_list_configs",
      null,
    );
  }

  // ---- Data sources ----

  async getDataSource(id: string): Promise<DataSourceConfig | null> {
    const req: GenericDAORequest<DataSourceConfig> = {
      action: "get",
      args: { id },
    };
    const res = await this.client.call<{ entry: DataSourceConfig | null }>(
      DS,
      "platform_datasource_dao",
      req,
    );
    return res.entry ?? null;
  }

  async addDataSource(entry: DataSourceConfig): Promise<AddDataSourceResponse> {
    const req: GenericDAORequest<DataSourceConfig> = {
      action: "add",
      args: { entry },
    };
    return await this.client.call<AddDataSourceResponse>(
      DS,
      "platform_datasource_dao",
      req,
    );
  }

  async deleteDataSource(id: string): Promise<void> {
    const req: GenericDAORequest<DataSourceConfig> = {
      action: "delete",
      args: { id },
    };
    await this.client.call(DS, "platform_datasource_dao", req);
  }

  async listDataSource(): Promise<DataSourceConfig[]> {
    const req: GenericDAORequest<DataSourceConfig> = { action: "list" };
    const res = await this.client.call<{ entries: DataSourceConfig[] }>(
      DS,
      "platform_datasource_dao",
      req,
    );
    return res.entries ?? [];
  }

  // ---- Data sinks ----

  async getDataSink(id: string): Promise<DataSinkConfig | null> {
    const req: GenericDAORequest<DataSinkConfig> = {
      action: "get",
      args: { id },
    };
    const res = await this.client.call<{ entry: DataSinkConfig | null }>(
      DS,
      "platform_datasink_dao",
      req,
    );
    return res.entry ?? null;
  }

  async addDataSink(entry: DataSinkConfig): Promise<AddDataSinkResponse> {
    const req: GenericDAORequest<DataSinkConfig> = {
      action: "add",
      args: { entry },
    };
    return await this.client.call<AddDataSinkResponse>(
      DS,
      "platform_datasink_dao",
      req,
    );
  }

  async deleteDataSink(id: string): Promise<void> {
    const req: GenericDAORequest<DataSinkConfig> = {
      action: "delete",
      args: { id },
    };
    await this.client.call(DS, "platform_datasink_dao", req);
  }

  async listDataSink(): Promise<DataSinkConfig[]> {
    const req: GenericDAORequest<DataSinkConfig> = { action: "list" };
    const res = await this.client.call<{ entries: DataSinkConfig[] }>(
      DS,
      "platform_datasink_dao",
      req,
    );
    return res.entries ?? [];
  }

  // ---- Routers + channels ----

  async getRouter(id: string): Promise<RouterEntryResponse> {
    const req: GenericDAORequest<RouterConfig> = {
      action: "get",
      args: { id },
    };
    return await this.client.call<RouterEntryResponse>(
      DS,
      "platform_router_dao",
      req,
    );
  }

  async addSimpleRouter(processorName: string, routerName: string): Promise<string> {
    const router = routerName === "" ? processorName : routerName;
    const req = {
      processor: processorName,
      router,
      pipe: processorName,
    };
    const res = await this.client.call<AddRouterResponse>(
      DS,
      "platform_add_simple_router",
      req,
    );
    return res.id;
  }

  async addRouter(entry: RouterConfig): Promise<AddRouterResponse> {
    const req: GenericDAORequest<RouterConfig> = {
      action: "add",
      args: { entry },
    };
    return await this.client.call<AddRouterResponse>(
      DS,
      "platform_router_dao",
      req,
    );
  }

  async deleteRouter(id: string): Promise<void> {
    const req: GenericDAORequest<RouterConfig> = {
      action: "delete",
      args: { id },
    };
    await this.client.call(DS, "platform_router_dao", req);
  }

  async getChannel(id: string): Promise<ChannelConfig | null> {
    const req: GenericDAORequest<ChannelConfig> = {
      action: "get",
      args: { id },
    };
    const res = await this.client.call<{ entry: ChannelConfig | null }>(
      DS,
      "platform_channel_dao",
      req,
    );
    return res.entry ?? null;
  }

  async addChannel(entry: ChannelConfig): Promise<AddChannelResponse> {
    const req: GenericDAORequest<ChannelConfig> = {
      action: "add",
      args: { entry },
    };
    return await this.client.call<AddChannelResponse>(
      DS,
      "platform_channel_dao",
      req,
    );
  }

  async deleteChannel(id: string): Promise<void> {
    const req: GenericDAORequest<ChannelConfig> = {
      action: "delete",
      args: { id },
    };
    await this.client.call(DS, "platform_channel_dao", req);
  }

  // ---- Processors ----

  async listProcessors(): Promise<FPLScript[]> {
    const req: GenericDAORequest<FPLScript> = { action: "list" };
    const res = await this.client.call<{ entries: FPLScript[] }>(
      DS,
      "platform_processor_dao",
      req,
    );
    return res.entries ?? [];
  }

  async getProcessor(name: string): Promise<FPLScript | null> {
    const req: GenericDAORequest<FPLScript> = {
      action: "get",
      args: { id: name },
    };
    const res = await this.client.call<{ entry: FPLScript | null }>(
      DS,
      "platform_processor_dao",
      req,
    );
    return res.entry ?? null;
  }

  async addProcessor(entry: FPLScript): Promise<void> {
    const req: GenericDAORequest<FPLScript> = {
      action: "add",
      args: { entry },
    };
    await this.client.call(DS, "platform_processor_dao", req);
  }

  async updateProcessor(entry: FPLScript): Promise<void> {
    const req: GenericDAORequest<FPLScript> = {
      action: "update",
      args: { entry },
    };
    await this.client.call(DS, "platform_processor_dao", req);
  }

  async deleteProcessor(name: string): Promise<void> {
    const req: GenericDAORequest<FPLScript> = {
      action: "delete",
      args: { id: name },
    };
    await this.client.call(DS, "platform_processor_dao", req);
  }

  async validateProcessor(
    req: FPLProcessorValidateRequest,
  ): Promise<FPLProcessorValidateResult> {
    return await this.client.call<FPLProcessorValidateResult>(
      DS,
      "platform_processor_validate",
      req,
    );
  }

  /**
   * Runs a processor script against sample data. An unset `type` is sent as
   * `fpl_processor`, which is the runtime the endpoint falls back to anyway, so
   * that the request on the wire says which one it meant.
   *
   * A script that fails on the document is not a call failure: the message
   * lands in the result's `error` field with an empty `status`, and only
   * transport and endpoint errors -- an empty `source` among them -- throw.
   */
  async testProcessor(
    req: FPLProcessorTestRequest,
  ): Promise<FPLProcessorTestResult> {
    return await this.client.call<FPLProcessorTestResult>(
      DS,
      "platform_processor_test",
      { ...req, type: req.type || "fpl_processor" },
    );
  }

  /**
   * Runs an `fpl_processor` script against one event object, wrapping it in the
   * document envelope the endpoint expects in its `source` field.
   */
  async testProcessorObject(
    script: string,
    obj: unknown,
  ): Promise<FPLProcessorTestResult> {
    return await this.testProcessor({
      script,
      source: buildFplTestSource(obj),
      type: "fpl_processor",
    });
  }

  // ---- Source-router wiring + pipes ----

  async setDataSourceRouter(req: SourceSetRouterReq): Promise<void> {
    await this.client.call(DS, "platform_source_set_router", req);
  }

  async addRouterPipe(req: RouterAddPipeReq): Promise<RouterAddPipeResponse> {
    return await this.client.call<RouterAddPipeResponse>(
      DS,
      "platform_router_add_pipe",
      req,
    );
  }

  async deleteRouterPipe(req: RouterDeletePipeReq): Promise<void> {
    await this.client.call(DS, "platform_router_delete_pipe", req);
  }

  async updateRouterPipes(req: RouterUpdatePipesReq): Promise<void> {
    await this.client.call(DS, "platform_router_update_pipes", req);
  }

  async updatePipe(req: PipeUpdateReq): Promise<void> {
    await this.client.call(DS, "platform_pipe_update", req);
  }

  async updatePipeProcessor(req: PipeProcessorUpdateReq): Promise<void> {
    await this.client.call(DS, "platform_pipe_update_processor", req);
  }

  // ---- Monitoring / tail / metrics ----

  async eventTail(req: EventTailReq): Promise<EventTailResponse> {
    return await this.client.call<EventTailResponse>(DS, "platform_event_tail", req);
  }

  async processorTail(req: ProcessorTailReq): Promise<ProcessorTailResponse> {
    return await this.client.call<ProcessorTailResponse>(
      DS,
      "platform_processor_tail",
      req,
    );
  }

  async processorPipes(req: ProcessorPipesReq): Promise<ProcessorPipesResponse> {
    return await this.client.call<ProcessorPipesResponse>(
      DS,
      "platform_processor_pipes",
      req,
    );
  }

  async componentMetrics(req: ComponentMetricReq): Promise<PlatformMetricsResponse> {
    return await this.client.call<PlatformMetricsResponse>(
      DS,
      "platform_component_metrics",
      req,
    );
  }

  async processorMetrics(req: ProcessorMetricReq): Promise<PlatformMetricsResponse> {
    return await this.client.call<PlatformMetricsResponse>(
      DS,
      "platform_processor_metrics",
      req,
    );
  }

  async getComponentState(id: string): Promise<GetComponentStateResponse> {
    return await this.client.call<GetComponentStateResponse>(
      DS,
      "platform_get_component_state",
      { id },
    );
  }

  async setComponentTags(req: SetComponentTagsReq): Promise<void> {
    await this.client.call(DS, "platform_set_component_tags", req);
  }

  async pluginTail(req: PluginTailReq): Promise<PluginTailResponse> {
    return await this.client.call<PluginTailResponse>(DS, "platform_plugin_tail", req);
  }

  async listComponentErrors(id: string): Promise<ListComponentErrorResponse> {
    return await this.client.call<ListComponentErrorResponse>(
      DS,
      "platform_list_component_errors",
      { id },
    );
  }

  async clearComponentError(id: string): Promise<void> {
    await this.client.call(DS, "platform_clear_component_error", { id });
  }

  async getComponentInfo(id: string): Promise<ComponentInfo> {
    return await this.client.call<ComponentInfo>(DS, "platform_get_component_info", {
      id,
    });
  }

  async sourceReload(id: string): Promise<void> {
    await this.client.call(DS, "platform_source_reload", { dataSourceID: id });
  }

  // ---- Integrations ----

  async listIntegrations(): Promise<Integration[]> {
    const req: IntegrationDAORequest = { action: "list" };
    const res = await this.client.call<IntegrationDAOResponse>(
      DS,
      "platform_integration_dao",
      req,
    );
    return res.entries ?? [];
  }

  async getIntegration(id: string): Promise<Integration | null> {
    const req: IntegrationDAORequest = {
      action: "get",
      args: { id },
    };
    const res = await this.client.call<IntegrationDAOResponse>(
      DS,
      "platform_integration_dao",
      req,
    );
    return res.entry ?? null;
  }

  async addIntegration(entry: Integration): Promise<string> {
    const req: IntegrationDAORequest = {
      action: "add",
      args: { entry },
    };
    const res = await this.client.call<AddIntegrationResponse>(
      DS,
      "platform_integration_dao",
      req,
    );
    return res.id;
  }

  async updateIntegration(entry: Integration): Promise<void> {
    const req: IntegrationDAORequest = {
      action: "update",
      args: { entry },
    };
    await this.client.call(DS, "platform_integration_dao", req);
  }

  async deleteIntegration(id: string): Promise<void> {
    const req: IntegrationDAORequest = {
      action: "delete",
      args: { id },
    };
    await this.client.call(DS, "platform_integration_dao", req);
  }
}

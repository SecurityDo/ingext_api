import type {
  AWSFirehoseConfig,
  AWSKinesisStreamConfig,
  AWSLambdaConfig,
  S3NotificationConfig,
  Tag,
} from "./integration.js";

export interface DataSourceConfig {
  type: string;
  name?: string;
  id?: string;
  format?: string;
  tags?: Tag[];
  compression?: string;
  receiverName?: string;
  receiverInput?: string;
  s3?: S3SourceConfig;
  s3Notification?: S3NotificationSourceConfig;
  kinesis?: KinesisSourceConfig;
  prom?: PromSourceConfig;
  hec?: HecSourceConfig;
  webhook?: WebhookSourceConfig;
  plugin?: PluginSourceConfig;
  syslog?: SyslogSourceConfig;
  secret?: unknown;
  integrationPull?: IntegrationPullSourceConfig;
}

export interface RedisConfig {
  host: string;
  port: number;
  queue: string;
}

export interface RedisSourceConfig {
  redis: RedisConfig;
}

export interface HecSourceConfig {
  ack: boolean;
  url: string;
}

export interface WebhookSourceConfig {
  url: string;
  token: string;
  securityToken: string;
  signatureHeader: string;
}

export interface PluginSourceConfig {
  name?: string;
  id?: string;
  integration?: string;
}

export interface S3SourceConfig {
  plugin?: PluginSourceConfig;
  bucket: string;
  prefix?: string;
  begin?: string;
  end?: string;
}

export interface SyslogSourceConfig {
  path: string;
}

export interface RedisSinkConfig {
  redis: RedisConfig;
}

export interface S3SinkConfig {
  integrationID: string;
  plugin?: PluginSourceConfig;
  bucket: string;
  compression?: string;
  objectPath: string;
}

export interface LavaDBSinkConfig {
  tenant: string;
  index: string;
}

export interface DataSinkConfig {
  type: string;
  name: string;
  id: string;
  groupSupport?: boolean;
  groupLambda?: string;
  redis?: RedisSinkConfig;
  s3?: S3SinkConfig;
  hec?: HecSinkConfig;
  webhook?: WebhookSinkConfig;
  kinesis?: KinesisSinkConfig;
  firehose?: FirehoseSinkConfig;
  lambda?: AWSLambdaSinkConfig;
  dataLake?: DataLakeSinkConfig;
  prom?: PromSinkConfig;
  loki?: LokiSinkConfig;
  secret?: unknown;
  tags?: Tag[];
  packerConfig?: SinkPackerConfig;
  flushCount?: number;
  flushBuffer?: number;
  flushInterval?: number;
}

export interface SinkPackerConfig {
  name?: string;
}

export interface FirehoseSinkConfig {
  integrationID: string;
  plugin?: PluginSourceConfig;
  config?: AWSFirehoseConfig;
}

export interface AWSLambdaSinkConfig {
  integrationID: string;
  plugin?: PluginSourceConfig;
  config?: AWSLambdaConfig;
}

export interface KinesisSinkConfig {
  integrationID: string;
  plugin?: PluginSourceConfig;
  config?: AWSKinesisStreamConfig;
}

export interface HecSinkConfig {
  token: string;
  url: string;
  insecureSkipVerify: boolean;
}

export interface RouterConfig {
  name: string;
  id: string;
  workerCount: number;
  pipeIDs: string[];
}

export interface StreamPipeConfig {
  name: string;
  id: string;
  routerID: string;
  matchAll: boolean;
  selector?: string;
  processorNames: string[];
  channelID?: string;
  sinkIDs?: string[];
}

export interface DataLakeSinkConfig {
  datalake?: string;
  datalakeIndex?: string;
  schemaName: string;
}

export interface PromSinkConfig {
  integrationID: string;
  plugin?: PluginSourceConfig;
}

export interface LokiSinkConfig {
  integrationID: string;
  plugin?: PluginSourceConfig;
}

export interface ChannelConfig {
  name: string;
  id: string;
  sinkIDs: string[];
}

export interface RouterInput {
  routerID: string;
  sourceID: string;
}

export interface PluginNotification {
  severity: string;
  subject: string;
  id?: string;
  source?: string;
  createdOn?: string;
  message: string;
}

export interface ComponentErrorState {
  id: string;
  name: string;
  errors: PluginNotification[];
  alerts: PluginNotification[];
}

export interface LogEvent {
  source?: string;
  timestamp?: string;
  msg?: string;
}

export interface ComponentInfo {
  id: string;
  name: string;
  errors: PluginNotification[];
  alerts: PluginNotification[];
  state: unknown;
  logs: LogEvent[];
}

export interface PromSourceConfig {
  url: string;
  authorization: string;
  username: string;
}

export interface PromSecret {
  token?: string;
  password?: string;
}

export interface S3NotificationSourceConfig {
  plugin?: PluginSourceConfig;
  config?: S3NotificationConfig;
}

export interface KinesisSourceConfig {
  integrationID: string;
  plugin?: PluginSourceConfig;
  config?: AWSKinesisStreamConfig;
}

export interface IntegrationPullSourceConfig {
  plugin?: PluginSourceConfig;
  processor: string;
  pollInterval: number;
}

export interface DelaySinkConfig {
  delay: number;
}

export interface HttpHeader {
  header: string;
  value: string;
}

export interface WebhookSinkConfig {
  url: string;
  headers: HttpHeader[];
  insecureSkipVerify: boolean;
}

export interface ElasticSinkConfig {
  url: string;
  index: string;
  dataType: string;
}

export interface ChangeEvent {
  timestamp: string;
  eventType: string;
  entry: unknown;
  revision: number;
}

export interface FPLScript {
  id: number;
  repository?: string;
  group: string;
  gitPath?: string;
  name: string;
  type: string;
  tags?: Tag[];
  description: string;
  scriptText: string;
  createdOn: string;
  updatedOn: string;
}

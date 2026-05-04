/** Minimal mirror of `golang.org/x/oauth2.Token` JSON shape. */
export interface OAuth2Token {
  access_token: string;
  token_type?: string;
  refresh_token?: string;
  expiry?: string;
  expires_in?: number;
}

export interface OAuth2Endpoint {
  AuthURL?: string;
  DeviceAuthURL?: string;
  TokenURL?: string;
  AuthStyle?: number;
}

/** Minimal mirror of `golang.org/x/oauth2.Config`. */
export interface OAuth2Config {
  ClientID?: string;
  ClientSecret?: string;
  Endpoint?: OAuth2Endpoint;
  RedirectURL?: string;
  Scopes?: string[];
}

export const PLUGIN_SERVICE_ENABLE = "enable";
export const PLUGIN_SERVICE_DISABLE = "disable";
export const PLUGIN_SERVICE_STOP = "stop";
export const PLUGIN_SERVICE_UPDATE = "update";
export const PLUGIN_SERVICE_INIT = "init";

export const PLUGIN_EVENT_SAVESTATE = "saveState";
export const PLUGIN_EVENT_EVENTDUMP = "eventDump";
export const PLUGIN_EVENT_TRANSLATIONDUMP = "translationDump";
export const PLUGIN_EVENT_RESOURCEDUMP = "resourceDump";
export const PLUGIN_EVENT_ELASTICDUMP = "elasticDump";
export const PLUGIN_EVENT_NOTIFICATION = "notification";

export const INTEGRATION_CATO = "Cato";
export const INTEGRATION_RESTAPI = "RESTAPI";
export const INTEGRATION_SALESFORCE = "Salesforce";
export const INTEGRATION_MIMECAST = "Mimecast";
export const INTEGRATION_PROOFPOINT = "Proofpoint";
export const INTEGRATION_FALCON = "Falcon";
export const INTEGRATION_OKTA = "Okta";
export const INTEGRATION_SOPHOS = "Sophos";
export const INTEGRATION_DARKTRACE = "Darktrace";
export const INTEGRATION_SENTINELONE = "SentinelOne";
export const INTEGRATION_LDAP = "LDAP";
export const INTEGRATION_DATALAKE = "DataLake";

export const INTEGRATION_AWS_KINESIS_STREAM = "KinesisStream";
export const INTEGRATION_AWS_S3BUCKET = "S3Bucket";
export const INTEGRATION_AWS_GUARDDUTY = "AWSGuardDuty";
export const INTEGRATION_AWS_LAMBDA = "AWSLambda";
export const INTEGRATION_AWS_FIREHOSE = "FirehoseStream";
export const INTEGRATION_AWS_EC2AUDIT = "AWSEC2Audit";
export const INTEGRATION_AWS_S3NOTIFICATION = "S3Notification";

export const INTEGRATION_AZURE_BLOBSTORAGE = "AzureBlobStorage";
export const INTEGRATION_PROM_PUSH = "PROMPush";
export const INTEGRATION_AWS_API = "AWSAPI";
export const INTEGRATION_AWS_USER = "AWSUser";

export interface Integration {
  id: string;
  name: string;
  integration: string;
  description: string;
  createdOn: string;
  updatedOn: string;
  secret?: unknown;
  config: unknown;
  inUse?: boolean;
  plugin?: boolean;
  tags?: Tag[];
}

export interface Tag {
  name: string;
  value: string;
}

export interface ClientSecret {
  clientSecret: string;
  privateKey: string;
}

export interface Office365AuditConfig {
  tenantID: string;
  clientId: string;
  thumbprint: string;
}

export interface Office365AuditState {
  lastPoll: string;
  streamMap: Record<string, number>;
}

export interface Office365ResourceWatchConfig {
  tenantID: string;
  clientId: string;
}

export interface Office365ResourceWatchSecret {
  clientSecret: string;
}

export interface AzureAuditConfig {
  tenantID: string;
  clientId: string;
}

export interface AzureAuditSecret {
  clientSecret: string;
}

export interface SlackConfig {
  token: string;
}

export interface PagerDutyConfig {
  url: string;
}

export interface PagerDutySecret {
  integrationKey: string;
  token: string;
}

export interface CatoConfig {
  roleARN: string;
  prefix: string;
  begin: string;
  end: string;
  bucket: string;
  region: string;
  mode: string;
  accessKey: string;
  role?: string;
}

export interface AzureBlobStorageConfig {
  storageAccount: string;
  container: string;
  authenticationMethod: string;
  tenantID?: string;
  clientID?: string;
}

export interface AzureBlobStorageSecret {
  connectionString?: string;
  clientSecret?: string;
}

export interface S3BucketConfig {
  prefix: string;
  begin: string;
  end: string;
  bucket: string;
  region: string;
  mode: string;
  accessKey: string;
  user?: string;
  role?: string;
}

export interface AWSUserSecret {
  secretKey: string;
}

export interface AWSAuthInfoCopy {
  secretKey?: string;
  accessKey?: string;
  roleObj?: import("./pod_role.js").InstanceRole;
}

export interface GSuiteSecret {
  clientSecret: string;
  token: OAuth2Token;
}

export interface GSuiteConfig {
  config: OAuth2Config;
}

export interface GSuiteState {
  lastPoll: string;
}

export interface DuoConfig {
  integrationKey: string;
  apiHostname: string;
}

export interface DuoSecret {
  secretKey: string;
}

export interface DuoState {
  stateMap: Record<string, number>;
}

export interface MSDefenderATPConfig {
  clientId: string;
  tenantId: string;
}

export interface MSDefenderATPSecret {
  clientSecret: string;
}

export interface OktaConfig {
  domain: string;
}

export interface OktaSecret {
  token: string;
}

export interface SophosEDRConfig {
  clientID: string;
}

export interface SophosEDRSecret {
  clientSecret: string;
}

export interface FalconConfig {
  clientID: string;
  baseURL: string;
}

export interface FalconSecret {
  clientSecret: string;
}

export interface FalconState {
  streamMap: Record<string, number>;
}

export interface MimecastConfig {
  applicationID: string;
  applicationKey: string;
  accessKey: string;
  baseURL: string;
}

export interface MimecastSecret {
  secretKey: string;
}

export interface MimecastState {
  lastPoll: string;
  apiTokenMap: Record<string, string>;
  apiPollMap: Record<string, string>;
}

export interface ProofPointState {
  lastPoll: string;
}

export interface ProofPointConfig {
  principal: string;
}

export interface ProofPointSecret {
  secret: string;
}

export interface EventhubConfig {
  eventhubs: AzureEventhubConfig[];
}

export interface AzureEventhubConfig {
  endpoint: string;
  description: string;
}

export interface EventhubState {
  lastPoll: string;
}

export interface AWSAuditConfig {
  accessKey: string;
  cloudTrails: AWSCloudTrail[];
  cloudWatches: AWSCloudWatch[];
}

export interface AWSAuditSecret {
  secretKey: string;
}

export interface AWSCloudTrail {
  name: string;
  region: string;
  queueURL?: string;
  bucket?: string;
}

export interface AWSCloudWatch {
  region: string;
  logGroups: string[];
}

export interface HECConfig {
  ack: boolean;
  endpoint: string;
  props: unknown;
}

export interface HECSecret {
  token: string;
}

export interface BitdefenderAPIConfig {
  url: string;
  companyID: string;
}

export interface BitdefenderAPISecret {
  apiKey: string;
}

export interface SalesforceConfig {
  baseURL: string;
  clientID: string;
}

export interface SalesforceSecret {
  clientSecret: string;
}

export interface SQSConfig {
  url: string;
  region: string;
  mode: string;
  accessKey: string;
  role?: string;
}

export interface SQSSecret {
  secretKey: string;
}

export interface S3NotificationConfig {
  queueARN: string;
  queueURL: string;
  region: string;
  accessKey?: string;
  user?: string;
  role?: string;
}

export interface S3NotificationSecret {
  secretKey: string;
}

export interface ServiceNowConfig {
  url: string;
  username: string;
  scope: string;
  incidentTable: string;
}

export interface ServiceNowSecret {
  password: string;
}

export interface AWSCloudWatchConfig {
  accessKey: string;
  cloudWatches: AWSCloudWatch[];
}

export interface AWSCloudWatchSecret {
  secretKey: string;
}

export interface AWSCloudTrailConfig {
  accessKey: string;
  cloudTrails: AWSCloudTrail[];
}

export interface AWSCloudTrailSecret {
  secretKey: string;
}

export interface AWSGuardDutyConfig {
  regions: string[];
  accessKey?: string;
  user?: string;
  role?: string;
}

export interface AWSGuardDutySecret {
  secretKey: string;
}

export interface AWSKinesisStreamConfig {
  region: string;
  name: string;
  arn: string;
  mode: string;
  accessKey?: string;
  user?: string;
  role?: string;
}

export interface AWSKinesisStreamSecret {
  secretKey: string;
}

export interface AWSLambdaConfig {
  region: string;
  functionName: string;
  accessKey?: string;
  role?: string;
  user?: string;
}

export interface AWSLambdaSecret {
  secretKey: string;
}

export interface AWSFirehoseConfig {
  region: string;
  name: string;
  accessKey?: string;
  role?: string;
  user?: string;
}

export interface AWSFirehoseSecret {
  secretKey: string;
}

export interface DataLakeConfig {
  managed: boolean;
  storageIntegration: string;
  index: string;
  schemaName: string;
  schema: string;
}

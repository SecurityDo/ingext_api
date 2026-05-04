import type { DataSinkConfig, DataSourceConfig } from "./platform.js";

export interface TypeMeta {
  kind?: string;
  apiVersion?: string;
}

export interface ObjectMeta {
  name?: string;
  uid?: string;
  resourceVersion?: string;
  creationTimestamp?: string;
  labels?: Record<string, string>;
  annotations?: Record<string, string>;
}

export interface ApplicationTemplateConfig {
  name: string;
  displayName?: string;
  description?: string;
  icon?: string;
  category?: string;
  resourceGroups?: string[];
  parameters?: InputParameter[];
  output?: InputParameter[];
}

export interface ApplicationConfigSpec {
  displayName?: string;
  description?: string;
  adminConsent?: boolean;
  category?: string;
  icon?: string;
  resourceGroups?: string[];
  parameters?: InputParameter[];
  output?: InputParameter[];
}

export interface ApplicationConfigResource extends TypeMeta, ObjectMeta {
  metadata?: ObjectMeta;
  spec?: ApplicationConfigSpec;
}

export interface ApplicationTemplate {
  application?: ApplicationConfigResource;
  routers?: Record<string, RouterResource>;
  processors?: Record<string, ProcessorResource>;
  dataSources?: Record<string, DataSourceResource>;
  dataSinks?: Record<string, DataSinkResource>;
  pipes?: Record<string, PipeResource>;
  integrations?: Record<string, IntegrationResource>;
  awsRoles?: Record<string, AWSRoleResource>;
  collectors?: Record<string, CollectorResource>;
}

export interface AppTemplateEntry {
  config: ApplicationTemplateConfig;
  content: string;
}

export interface CollectorSpec {
  name?: string;
  description?: string;
  output?: InputParameter[];
}

export interface CollectorResource extends TypeMeta, ObjectMeta {
  metadata?: ObjectMeta;
  spec?: CollectorSpec;
}

export interface AWSRoleResource extends TypeMeta, ObjectMeta {
  metadata?: ObjectMeta;
  spec?: AWSRoleSpec;
}

export interface AWSRoleSpec {
  iamRoleARN?: string;
  externalID?: string;
}

export interface IntegrationSpec {
  type?: string;
  description?: string;
  adminConsent?: string;
  config?: unknown;
  secret?: unknown;
  output?: InputParameter[];
}

export interface IntegrationResource extends TypeMeta, ObjectMeta {
  metadata?: ObjectMeta;
  spec?: IntegrationSpec;
}

export interface GitObject {
  repo?: string;
  path?: string;
  branch?: string;
}

export interface ProcessorSpec {
  type?: string;
  local?: string;
  remote?: GitObject;
}

export interface ProcessorResource extends TypeMeta, ObjectMeta {
  metadata?: ObjectMeta;
  spec?: ProcessorSpec;
}

export interface DataSinkSpec {
  type?: string;
  config?: DataSinkConfig;
  packer?: string;
}

export interface DataSinkResource extends TypeMeta, ObjectMeta {
  metadata?: ObjectMeta;
  spec?: DataSinkSpec;
}

export interface DataSourceSpec {
  type?: string;
  format?: string;
  config?: DataSourceConfig;
  router?: string;
  receiver?: DataSourceReceiver;
  output?: InputParameter[];
}

export interface DataSourceResource extends TypeMeta, ObjectMeta {
  metadata?: ObjectMeta;
  spec?: DataSourceSpec;
}

export interface DataSourceReceiver {
  processor?: string;
  input?: string;
}

export interface PipeSpec {
  name?: string;
  processors?: string[];
  sinks?: string[];
  routerSelector?: LabelSelector;
  router?: string;
  priority?: number;
}

export interface RouterSpec {
  threadCount?: number;
}

export interface PipeResource extends TypeMeta, ObjectMeta {
  metadata?: ObjectMeta;
  spec?: PipeSpec;
  raw?: string;
}

export interface RouterResource extends TypeMeta, ObjectMeta {
  metadata?: ObjectMeta;
  spec?: RouterSpec;
  raw?: string;
}

export interface ValueReference {
  kind?: string;
  resource?: string;
  name?: string;
}

export interface EnumValue {
  value?: string;
  label?: string;
}

export interface InputParameter {
  name?: string;
  description?: string;
  defaultValue?: string;
  dataType?: string;
  sensitive?: boolean;
  optional?: boolean;
  value?: string;
  isList?: boolean;
  enums?: EnumValue[];
  valueRef?: ValueReference;
}

export interface InstallAppInstanceRequest {
  config: InstanceConfig;
}

export interface InstanceState {
  application?: string;
  appTemplate: ApplicationTemplate;
  config?: InstanceConfig;
  state?: string;
  errorMsg?: string;
}

export interface ResourceMeta {
  id?: string;
  resource?: string;
}

export interface AppInstanceAction {
  action?: string;
  application?: string;
  instance?: string;
  status?: string;
  errorMsg?: string;
  resources?: ResourceMeta[];
  args?: unknown;
}

export interface AppInstanceOutput {
  application?: string;
  instance?: string;
  resource?: string;
  kind?: string;
  name?: string;
  description?: string;
  value?: string;
  valueRef?: ValueReference;
}

export interface LabelSelector {
  matchLabels?: Record<string, string>;
  matchExpressions?: LabelSelectorRequirement[];
}

export interface LabelSelectorRequirement {
  key: string;
  operator: string;
  values?: string[];
}

export interface GenericResource extends TypeMeta, ObjectMeta {
  metadata?: ObjectMeta;
  spec?: Record<string, unknown>;
  raw?: string;
}

export interface InstanceConfig {
  application?: string;
  instance?: string;
  displayName?: string;
  inputParameters?: InputParameter[];
}

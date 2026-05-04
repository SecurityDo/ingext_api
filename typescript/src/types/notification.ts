export interface EndpointConfig {
  name: string;
  integration: string;
  action: string;
  email?: EndpointEmailConfig;
  slack?: EndpointSlackConfig;
}

export interface EndpointEmailConfig {
  to?: string[];
  cc?: string[];
}

export interface EndpointSlackConfig {
  channel?: string;
  channels?: string[];
  integrationName?: string;
}

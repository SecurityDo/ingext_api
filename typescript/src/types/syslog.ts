export interface SyslogConfig {
  domain: string;
  port_begin: number;
  port_end: number;
  syslog_tls?: boolean;
  tls_rfc6587?: boolean;
  syslog_udp?: boolean;
  syslog_tcp?: boolean;
  syslog_tls_port?: number;
  tls_rfc6587_port?: number;
  syslog_udp_port?: number;
  syslog_tcp_port?: number;
}

export interface SyslogPortRequest {
  syslog_tls: boolean;
  tls_rfc6587: boolean;
  syslog_udp: boolean;
  syslog_tcp: boolean;
}

export interface GetSyslogConfigResponse {
  config?: SyslogConfig;
}

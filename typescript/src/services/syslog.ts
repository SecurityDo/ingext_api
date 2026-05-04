import type { IngextClient } from "../client.js";
import type {
  GetSyslogConfigResponse,
  SyslogPortRequest,
} from "../types/syslog.js";

const DS = "api/ds";

export type SyslogPort = "udp" | "tcp" | "tls" | "tls-rfc6587";

function buildPortRequest(ports: SyslogPort[] | string[]): SyslogPortRequest {
  const req: SyslogPortRequest = {
    syslog_tls: false,
    tls_rfc6587: false,
    syslog_udp: false,
    syslog_tcp: false,
  };
  for (const p of ports) {
    switch (p) {
      case "udp":
        req.syslog_udp = true;
        break;
      case "tcp":
        req.syslog_tcp = true;
        break;
      case "tls":
        req.syslog_tls = true;
        break;
      case "tls-rfc6587":
        req.tls_rfc6587 = true;
        break;
    }
  }
  return req;
}

export class SyslogService {
  constructor(private client: IngextClient) {}

  async register(ports: SyslogPort[] | string[]): Promise<GetSyslogConfigResponse> {
    return await this.client.call<GetSyslogConfigResponse>(
      DS,
      "ingext_syslog_register_config",
      buildPortRequest(ports),
    );
  }

  async update(ports: SyslogPort[] | string[]): Promise<GetSyslogConfigResponse> {
    return await this.client.call<GetSyslogConfigResponse>(
      DS,
      "ingext_syslog_update_config",
      buildPortRequest(ports),
    );
  }

  async get(): Promise<GetSyslogConfigResponse> {
    return await this.client.call<GetSyslogConfigResponse>(
      DS,
      "ingext_syslog_get_config",
      null,
    );
  }

  async delete(): Promise<void> {
    await this.client.call(DS, "ingext_syslog_delete_config", null);
  }
}

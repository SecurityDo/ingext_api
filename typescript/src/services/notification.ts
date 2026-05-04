import type { IngextClient } from "../client.js";
import type { GenericDAORequest, GenericDaoAddResponse, GenericDaoListResponse } from "../types/dao.js";
import type { EndpointConfig } from "../types/notification.js";

const DS = "api/ds";

export class NotificationService {
  constructor(private client: IngextClient) {}

  async addEmail(
    name: string,
    action: string,
    to: string[],
    cc: string[],
  ): Promise<string> {
    const req: GenericDAORequest<EndpointConfig> = {
      action: "add",
      args: {
        entry: {
          name,
          integration: "Email",
          action,
          email: { to, cc },
        },
      },
    };
    const res = await this.client.call<GenericDaoAddResponse>(
      DS,
      "platform_notification_endpoint_dao",
      req,
    );
    return res.id;
  }

  async delete(name: string): Promise<void> {
    const req: GenericDAORequest<EndpointConfig> = {
      action: "delete",
      args: { id: name },
    };
    await this.client.call(DS, "platform_notification_endpoint_dao", req);
  }

  async list(): Promise<EndpointConfig[]> {
    const req: GenericDAORequest<EndpointConfig> = { action: "list" };
    const res = await this.client.call<GenericDaoListResponse<EndpointConfig>>(
      DS,
      "platform_notification_endpoint_dao",
      req,
    );
    return res.entries ?? [];
  }
}

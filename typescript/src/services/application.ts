import type { IngextClient } from "../client.js";
import type {
  AppTemplateEntry,
  ApplicationTemplateConfig,
  InputParameter,
  InstanceConfig,
  InstanceState,
} from "../types/application.js";
import type { GenericDAORequest, GenericDaoAddResponse } from "../types/dao.js";

const DS = "api/ds";

export interface ListAppTemplateResponse {
  entries: ApplicationTemplateConfig[];
}

export interface ListAppInstanceResponse {
  entries: InstanceState[];
}

export interface InstallAppInstanceRequest {
  config: InstanceConfig;
}

export interface UnInstallAppInstanceRequest {
  application?: string;
  instance?: string;
}

export interface GetAppInstanceRequest {
  application?: string;
  instance?: string;
}

export interface GetAppInstanceResponse {
  outputs?: InputParameter[];
}

export class ApplicationService {
  constructor(private client: IngextClient) {}

  async listAppTemplates(): Promise<ListAppTemplateResponse> {
    return await this.client.call<ListAppTemplateResponse>(
      DS,
      "platform_list_application_template",
      null,
    );
  }

  async listAppInstances(): Promise<ListAppInstanceResponse> {
    return await this.client.call<ListAppInstanceResponse>(
      DS,
      "platform_list_application_template",
      null,
    );
  }

  async installAppTemplate(req: InstallAppInstanceRequest): Promise<void> {
    await this.client.call(DS, "platform_install_application_instance", req);
  }

  async uninstallAppTemplate(req: UnInstallAppInstanceRequest): Promise<void> {
    await this.client.call(DS, "platform_uninstall_application_instance", req);
  }

  async addAppTemplate(content: string): Promise<string> {
    const req: GenericDAORequest<AppTemplateEntry> = {
      action: "add",
      args: {
        entry: { config: { name: "" }, content },
      },
    };
    const res = await this.client.call<GenericDaoAddResponse>(
      DS,
      "platform_application_template_dao",
      req,
    );
    return res.id;
  }

  async deleteAppTemplate(name: string): Promise<void> {
    const req: GenericDAORequest<AppTemplateEntry> = {
      action: "delete",
      args: { id: name },
    };
    await this.client.call(DS, "platform_application_template_dao", req);
  }

  async updateAppTemplate(name: string, content: string): Promise<void> {
    const req: GenericDAORequest<AppTemplateEntry> = {
      action: "update",
      args: {
        id: name,
        entry: { config: { name: "" }, content },
      },
    };
    await this.client.call(DS, "platform_application_template_dao", req);
  }

  async getAppInstance(
    app: string,
    instance: string,
  ): Promise<GetAppInstanceResponse> {
    const req: GetAppInstanceRequest = { application: app, instance };
    return await this.client.call<GetAppInstanceResponse>(
      DS,
      "platform_get_application_instance",
      req,
    );
  }
}

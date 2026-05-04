import type { IngextClient } from "../client.js";
import type {
  Datalake,
  DatalakeIndex,
  DatalakeIndexAddRequest,
  DatalakeIndexDeleteRequest,
  DatalakeIndexListRequest,
  DatalakeIndexListResponse,
  SchemaEntry,
} from "../types/datalake.js";
import type { GenericDAORequest, GenericDaoListResponse } from "../types/dao.js";

const DS = "api/ds";

export class DatalakeService {
  constructor(private client: IngextClient) {}

  async listDatalake(): Promise<Datalake[]> {
    const req: GenericDAORequest<Datalake> = { action: "list" };
    const res = await this.client.call<GenericDaoListResponse<Datalake>>(
      DS,
      "ingext_datalake_dao",
      req,
    );
    return res.entries ?? [];
  }

  async addDatalake(name: string, managed: boolean, integrationID: string): Promise<void> {
    const req: GenericDAORequest<Datalake> = {
      action: "add",
      args: {
        entry: {
          name: managed ? "managed" : name,
          managed,
          integrationID,
        },
      },
    };
    await this.client.call(DS, "ingext_datalake_dao", req);
  }

  async addDatalakeIndex(datalake: string, index: string, schema: string): Promise<void> {
    const req: DatalakeIndexAddRequest = {
      entry: {
        datalake,
        datalakeIndex: index,
        schemaName: schema,
      },
    };
    await this.client.call(DS, "ingext_datalake_index_add", req);
  }

  async deleteDatalakeIndex(datalake: string, index: string): Promise<void> {
    const req: DatalakeIndexDeleteRequest = { lake: datalake, index };
    await this.client.call(DS, "ingext_datalake_index_delete", req);
  }

  async listDatalakeIndex(datalake: string): Promise<DatalakeIndex[]> {
    const req: DatalakeIndexListRequest = { lake: datalake };
    const res = await this.client.call<DatalakeIndexListResponse>(
      DS,
      "ingext_datalake_index_list",
      req,
    );
    return res.entries ?? [];
  }

  async listSchema(): Promise<SchemaEntry[]> {
    const req: GenericDAORequest<SchemaEntry> = { action: "list" };
    const res = await this.client.call<GenericDaoListResponse<SchemaEntry>>(
      DS,
      "ingext_datalake_schema_dao",
      req,
    );
    return res.entries ?? [];
  }

  async addSchema(name: string, description: string, content: string): Promise<void> {
    const req: GenericDAORequest<SchemaEntry> = {
      action: "add",
      args: {
        entry: {
          name,
          description,
          content,
          createdAt: "",
        },
      },
    };
    await this.client.call(DS, "ingext_datalake_schema_dao", req);
  }

  async updateSchema(name: string, description: string, content: string): Promise<void> {
    const req: GenericDAORequest<SchemaEntry> = {
      action: "update",
      args: {
        entry: {
          name,
          description,
          content,
          createdAt: "",
        },
      },
    };
    await this.client.call(DS, "ingext_datalake_schema_dao", req);
  }

  async deleteSchema(name: string): Promise<void> {
    const req: GenericDAORequest<SchemaEntry> = {
      action: "delete",
      args: { id: name },
    };
    await this.client.call(DS, "ingext_datalake_schema_dao", req);
  }
}

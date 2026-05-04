import type { IngextClient } from "../client.js";
import type {
  GridAddSaasAccountRequest,
  GridDeleteSaasAccountRequest,
  ListFluencyAccountsResponse,
} from "../types/grid.js";

const GRID = "api/grid";

export class GridService {
  constructor(private client: IngextClient) {}

  async listAccount(): Promise<ListFluencyAccountsResponse> {
    return await this.client.call<ListFluencyAccountsResponse>(
      GRID,
      "get_grid_accounts",
      null,
    );
  }

  async addSaasAccount(req: GridAddSaasAccountRequest): Promise<void> {
    await this.client.call(GRID, "add_saas_account", req);
  }

  async deleteSaasAccount(req: GridDeleteSaasAccountRequest): Promise<void> {
    await this.client.call(GRID, "delete_saas_account", req);
  }
}

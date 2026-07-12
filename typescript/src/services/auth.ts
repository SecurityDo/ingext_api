import type { IngextClient } from "../client.js";
import type {
  ApiPolicy,
  ApiTokenEntry,
  Role,
  SitePolicy,
  UserEntry,
} from "../types/auth.js";
import type { GenericDAORequest } from "../types/dao.js";

const AUTH = "api/auth";

export interface AddUserRequest {
  user: UserEntry;
}

export class AuthService {
  constructor(private client: IngextClient) {}

  async addUser(req: AddUserRequest): Promise<void> {
    await this.client.call("api/auth", "userAdd", req);
  }

  async listUser(): Promise<UserEntry[]> {
    const res = await this.client.call<{ users: UserEntry[] }>(
      "api/auth",
      "userList",
      null,
    );
    return res.users ?? [];
  }

  async getUser(username: string): Promise<UserEntry> {
    const res = await this.client.call<{ user: UserEntry }>(
      "api/auth",
      "getUser",
      { username },
    );
    return res.user;
  }

  async deleteUser(username: string): Promise<void> {
    await this.client.call("api/auth", "userDelete", { username });
  }

  async addToken(name: string, description: string, role: string): Promise<string> {
    const res = await this.client.call<{ token: string }>("api/auth", "api_token", {
      action: "add",
      args: {
        entry: {
          name,
          description,
          roles: [role],
        },
      },
    });
    return res.token;
  }

  async deleteToken(name: string): Promise<void> {
    await this.client.call("api/auth", "api_token", {
      action: "delete",
      args: { id: name },
    });
  }

  async listToken(): Promise<ApiTokenEntry[]> {
    const res = await this.client.call<{ entries: ApiTokenEntry[] }>(
      "api/auth",
      "api_token",
      { action: "list" },
    );
    return res.entries ?? [];
  }

  async setUserSitePolicy(username: string, sitePolicy: string): Promise<void> {
    await this.client.call("api/auth", "setUserSitePolicy", {
      username,
      policyName: sitePolicy,
    });
  }

  // --- Role DAO (backed by the roleDao endpoint under api/auth) ---

  /** Create a new RBAC role. A role must reference at least one data or API policy. */
  async addRole(entry: Role): Promise<void> {
    const req: GenericDAORequest<Role> = {
      action: "create",
      args: { entry },
    };
    await this.client.call(AUTH, "roleDao", req);
  }

  /** List all RBAC roles. */
  async listRole(): Promise<Role[]> {
    const req: GenericDAORequest<Role> = { action: "list" };
    const res = await this.client.call<{ entries: Role[] }>(AUTH, "roleDao", req);
    return res.entries ?? [];
  }

  /** Delete the RBAC role identified by name. */
  async deleteRole(name: string): Promise<void> {
    const req: GenericDAORequest<Role> = {
      action: "delete",
      args: { id: name },
    };
    await this.client.call(AUTH, "roleDao", req);
  }

  // --- API policy DAO (backed by the apiPolicyDao endpoint under api/auth) ---

  /** Create a new API policy. A policy must define at least one resource. */
  async addApiPolicy(entry: ApiPolicy): Promise<void> {
    const req: GenericDAORequest<ApiPolicy> = {
      action: "create",
      args: { entry },
    };
    await this.client.call(AUTH, "apiPolicyDao", req);
  }

  /** List all API policies. */
  async listApiPolicy(): Promise<ApiPolicy[]> {
    const req: GenericDAORequest<ApiPolicy> = { action: "list" };
    const res = await this.client.call<{ entries: ApiPolicy[] }>(
      AUTH,
      "apiPolicyDao",
      req,
    );
    return res.entries ?? [];
  }

  /** Delete the API policy identified by name. */
  async deleteApiPolicy(name: string): Promise<void> {
    const req: GenericDAORequest<ApiPolicy> = {
      action: "delete",
      args: { id: name },
    };
    await this.client.call(AUTH, "apiPolicyDao", req);
  }

  // --- Site policy DAO (backed by the sitePolicyDao endpoint under api/auth) ---

  /** Create a new site policy. A policy must reference at least one user or token. */
  async addSitePolicy(entry: SitePolicy): Promise<void> {
    const req: GenericDAORequest<SitePolicy> = {
      action: "create",
      args: { entry },
    };
    await this.client.call(AUTH, "sitePolicyDao", req);
  }

  /** List all site policies. */
  async listSitePolicy(): Promise<SitePolicy[]> {
    const req: GenericDAORequest<SitePolicy> = { action: "list" };
    const res = await this.client.call<{ entries: SitePolicy[] }>(
      AUTH,
      "sitePolicyDao",
      req,
    );
    return res.entries ?? [];
  }

  /** Delete the site policy identified by name. */
  async deleteSitePolicy(name: string): Promise<void> {
    const req: GenericDAORequest<SitePolicy> = {
      action: "delete",
      args: { id: name },
    };
    await this.client.call(AUTH, "sitePolicyDao", req);
  }
}

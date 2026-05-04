import type { IngextClient } from "../client.js";
import type { ApiTokenEntry, UserEntry } from "../types/auth.js";

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
}

export interface UserEntry {
  username?: string;
  email?: string;
  firstName?: string;
  lastName?: string;
  organization?: string;
  roles: string[];
  oauthProvider: string;
  oauthFlag: boolean;
  dataPolicies?: string[];
  APIPolicies?: string[];
}

/**
 * An RBAC role that bundles a set of data and API policies.
 * Backed by the roleDao endpoint under api/auth.
 */
export interface Role {
  name: string;
  description?: string;
  dataPolicies: string[];
  apiPolicies: string[];
}

/** Grants a set of privileges (e.g. "read", "write") on a resource. */
export interface ResourcePolicy {
  resource: string;
  privileges: string[];
}

/**
 * A named collection of resource privileges.
 * Backed by the apiPolicyDao endpoint under api/auth.
 */
export interface ApiPolicy {
  name: string;
  description?: string;
  resources: ResourcePolicy[];
}

/**
 * Scopes which accounts a set of users and tokens may access.
 * Backed by the sitePolicyDao endpoint under api/auth.
 */
export interface SitePolicy {
  name: string;
  description?: string;
  accounts: string[];
  users: string[];
  tokens: string[];
}

export interface ApiTokenEntry {
  id?: string;
  name?: string;
  description?: string;
  disabled: boolean;
  token?: string;
  roles: string[];
  dataPolicies: string[];
  APIPolicies: string[];
  lastAccess: number;
  createdOn: string;
  updatedOn: string;
}

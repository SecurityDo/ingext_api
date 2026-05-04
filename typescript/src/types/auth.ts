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

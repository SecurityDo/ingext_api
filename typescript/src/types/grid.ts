export interface FluencyAccount {
  disabled: boolean;
  multiTenant: boolean;
  name: string;
  region: string;
  cluster: string;
  vpc: string;
  description: string;
  subnet: string;
  subnetID: string;
  serverNode: string;
  workerNodes: string[];
  token: string;
  URL: string;
  oauth2: boolean;
  createdOn: string;
  updatedOn: string;
}

export interface ListFluencyAccountsResponse {
  accounts: FluencyAccount[];
  entries: string[];
}

export interface GridAddSaasAccountRequest {
  name: string;
  region: string;
  cluster: string;
  displayName: string;
  siteURL: string;
  token: string;
}

export interface GridDeleteSaasAccountRequest {
  name: string;
}

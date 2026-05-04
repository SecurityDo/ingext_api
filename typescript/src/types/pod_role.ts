export interface InstanceRole {
  id: string;
  displayName: string;
  description: string;
  externalID: string;
  roleARN: string;
  createdOn: string;
  local?: boolean;
}

export interface InstanceRoleTestRequest {
  names: string[];
  roleARN: string;
  externalID: string;
}

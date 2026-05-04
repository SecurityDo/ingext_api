export interface Datalake {
  name?: string;
  managed?: boolean;
  storageDescription?: string;
  integrationID?: string;
  description?: string;
  createdAt?: string;
}

export interface DatalakeIndex {
  datalake?: string;
  datalakeIndex?: string;
  description?: string;
  schemaName?: string;
  storageDescription?: string;
}

export interface SchemaEntry {
  name: string;
  description: string;
  createdAt: string;
  updatedAt?: string;
  content: string;
}

export interface Field {
  name: string;
  type: string;
  convertedtype?: string;
  repetitiontype?: string;
  fields?: Field[];
  nullable?: boolean;
}

export interface Table {
  name: string;
  fields: Field[];
}

export interface DatalakeIndexListRequest {
  lake: string;
}

export interface DatalakeIndexListResponse {
  entries: DatalakeIndex[];
}

export interface DatalakeIndexAddRequest {
  entry: DatalakeIndex;
}

export interface DatalakeIndexDeleteRequest {
  lake: string;
  index: string;
}

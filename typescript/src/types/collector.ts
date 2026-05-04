export interface CollectorT {
  id: string;
  name: string;
  local: boolean;
  ebpf: boolean;
  description: string;
  createdOn: number;
  token: string;
  account?: string;
}

export interface CollectorForWeb extends CollectorT {
  lastPoll: number;
}

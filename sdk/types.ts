export interface Logger {
  info: (...args: any[]) => void;
  warn: (...args: any[]) => void;
  error: (...args: any[]) => void;
  debug: (...args: any[]) => void;
}

export type ClientParams = {
  providerUrl: string,
  contractAddr: string,
};

export type NetworkParams = {
  relayerIp: string,
  resolverPubKey: string,
};

export type JsonRequest = {
  Id: string;
  Method: string;
  Params: string[];
};

export type JsonResponse = {
  id: string;
  result: any;
};

export type PendingRequest = {
  resolve: any;
  reject: any;
  privKey: any;
}

export type ClientParams = {
  providerUrl: string,
  contractAddr: string,
  sessionId: string,
  signalingServer: string,
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

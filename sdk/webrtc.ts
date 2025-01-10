import axios from 'axios';
import { generateKeyPair, KeyType, encrypt, decrypt,  pemEncodedPublicKey, getKeyType } from "./crypto/util.ts";
import { ResolverRequestSchema, ResolverResponseSchema } from "./gen/resolver_pb.ts";
import { create, toBinary, toJson, toJsonString, fromJsonString} from "@bufbuild/protobuf";

export type ConnParams = {
  relayerAddress: string;
  channelName: string;
  stunServers: string[];
  encryptionType: KeyType;
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

type PendingRequest = {
  resolve: any;
  reject: any;
  keyPair: any;
}

let pc: RTCPeerConnection = null;
let connParams: ConnParams = null;
let pendingRequests = new Map<string, PendingRequest>();

const send = (url, msg) => {
    const headers = {
      'Content-Type': 'application/json'
    }
    return axios.post(connParams.relayerAddress + url, msg, {headers: headers})
}

const onicecandidate = ({candidate}) => {
  if (candidate !== null) {
    let sessionAndCandidate = {'session_id': 'firefox', 'candidate': candidate}
    log(`candidate: ${JSON.stringify(sessionAndCandidate)}`)
    send('/candidate', sessionAndCandidate)
  }
}

export const log = msg => {
  document.getElementById('logs').innerHTML += msg + '<br>'
}

let makingOffer = false
let sendChannel = null

const onnegotiationneeded = async () => {
  try {
    makingOffer = true;
    await pc.setLocalDescription();

    let sessionAndOffer = {'session_id': 'firefox', 'offer': pc.localDescription}
    let resp = await send('/sdp', sessionAndOffer)
    log(`resp: ${JSON.stringify(resp.data)}`)
    pc.setRemoteDescription(resp.data.answer)
  } catch (err) {
    console.error(err);
  } finally {
    makingOffer = false;
  }
};

const onmessage = async (ev) => {
  const data = ev.data;
  const protoResp = fromJsonString(ResolverRequestSchema, ev.data)  
  log(`chan msg: ${JSON.stringify(protoResp)}`)

  let payload = protoResp.payload
  let pendingReq = pendingRequests[protoResp.id];
  if (pendingReq == null) {
    return
  }
  let { resolve, reject, keyType, keyPair } = pendingReq;
  if (keyPair != null) {
    payload = await decrypt(keyType, keyPair, payload)
  } else {
    payload = new TextDecoder().decode(payload)
  }
  const resp: JsonResponse = JSON.parse(payload);
  log(`chan msg result: ${JSON.stringify(resp)}`)
  
  log(`resolve id: ${resp.id}`)
  resolve(resp);
  pendingRequests.delete(resp.id);
}

export async function connect(params: ConnParams) {
  pc = new RTCPeerConnection({ iceServers: [{urls: params.stunServers}]});
  connParams = params;
  pc.onnegotiationneeded = onnegotiationneeded
  pc.onicecandidate = onicecandidate
  sendChannel = pc.createDataChannel(params.channelName)
  log(`sendChannel: ${JSON.stringify(sendChannel)}`)
  sendChannel.onmessage = onmessage
  sendChannel.onopen = () => {
    log("channel open")
  }
  sendChannel.onclose = () => {
    log("channel closed")
  }
}

export async function execute(req: JsonRequest, pemPubKey: string): Promise<JsonResponse> {
  log("executing request")
  // Wrap it in ProtoBuf
  const reqStr = JSON.stringify(req);
  let keyPair = null;
  const keyType = getKeyType(pemPubKey);
  if (pemPubKey) {
    keyPair = await generateKeyPair(keyType);//connParams.encryptionType);
  }
  let reqStrBytes = new TextEncoder().encode(reqStr);
  log(`reqStrBytes: ${reqStrBytes}`);
  if (pemPubKey) {
    let reqStrEncr = await encrypt(pemPubKey, reqStrBytes)
    reqStrBytes = new Uint8Array(reqStrEncr)
  }
  let pkBytes = new TextEncoder().encode("PUBLICKEY");
  if (pemPubKey) {
    let pemEncodedKey = await pemEncodedPublicKey(keyType, keyPair)
    pkBytes = new TextEncoder().encode(pemEncodedKey);
  }
  let encrypted = pemPubKey != null;
  const protoReq = create(ResolverRequestSchema, {
    id: req.Id,
    payload: reqStrBytes,
    encrypted: encrypted,
    publicKey: pkBytes
  });
  log(`protoReq: ${protoReq}`);
  const reqJson = toJsonString(ResolverRequestSchema, protoReq);
  log(`reqJson: ${reqJson}`);
  sendChannel.send(reqJson);

  let resolve, reject;
  const promise = new Promise<JsonResponse>((res, rej) => {
    resolve = res;
    reject = rej;
  });

  log(`pending id: ${req.Id}`)
  pendingRequests[req.Id] = {resolve, reject, keyType, keyPair};
  return promise;
}

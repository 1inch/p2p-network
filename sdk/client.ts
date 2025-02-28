import { ClientParams, JsonRequest, JsonResponse, NetworkParams } from "./types";
import axios from 'axios';
import * as ecies from "eciesjs";
import { generateKeyPair, encrypt, decrypt } from "./crypto/util";
import { ResolverRequestSchema, ResolverResponseSchema } from "./gen/resolver_pb";
import { IncomingMessageSchema, OutgoingMessageSchema } from "./gen/relayer_pb";
import { createPublicClient, http } from 'viem'
import { registryAbi } from "./abi/NodeRegistry";
import { create, toJson, toJsonString, toBinary, fromBinary, fromJsonString} from "@bufbuild/protobuf";

type PendingRequest = {
  resolve: any;
  reject: any;
  privKey: any;
}

export class Client {
  pc: RTCPeerConnection;
  sendChannel: RTCDataChannel;
  makingOffer: boolean;
  networkParams: NetworkParams;
  connectionOpened: any;
  connectionClosed: any;
  pendingRequests: Map<string, PendingRequest>;

  constructor() {
    this.pendingRequests = new Map<string, PendingRequest>();
  }

  async init(params: ClientParams) {
    const networkParams = await this.fetchNetworkParams(params);
    this.networkParams = networkParams;
    let pc = new RTCPeerConnection({ iceServers: [{urls: defaultStunServers}]});

    this.pc = pc;

    pc.onnegotiationneeded = () => this.onnegotiationneeded();
    pc.onicecandidate = ({candidate}) => this.onicecandidate(candidate)

    let sendChannel = pc.createDataChannel(defaultChannelName);
    this.sendChannel = sendChannel;

    sendChannel.onmessage = (ev) => this.onmessage(ev);
    sendChannel.onopen = () => this.onopen();
    sendChannel.onclose = () => this.onclose();

    return new Promise<boolean>((res, rej) => {
      this.connectionOpened = res;
      this.connectionClosed = rej;
    });
  }

  onopen() {
      this.log("channel open")
      this.connectionOpened(true);
  }

  onclose() {
    this.log("channel closed")
      this.connectionClosed(true);
  }


  onicecandidate(candidate: RTCIceCandidate | null) {
    if (candidate !== null) {
      let sessionAndCandidate = {'session_id': 'firefox', 'candidate': candidate}
      this.log(`candidate: ${JSON.stringify(sessionAndCandidate)}`)
      this.send('/candidate', sessionAndCandidate)
    }
  }

  send(url: string, msg: { session_id: string; candidate?: RTCIceCandidate; offer?: RTCSessionDescription | null; }) {
    const headers = {
      'Content-Type': 'application/json'
    }
    console.log(`axios posting to url: ${url}`);
    let addr = "http://" + this.networkParams.relayerIp + url;
    console.log(`axios posting to addr: ${addr}`);
    return axios.post(addr, msg, {headers: headers})
  }

  log(msg: string) {
    document.getElementById('logs')!.innerHTML += msg + '<br>'
  }

  async encryptRequest(req: JsonRequest, resolverPubKey: string) {
    const reqStr = JSON.stringify(req);
    const reqStrBytes = new TextEncoder().encode(reqStr);
    this.log(`reqStrBytes: ${reqStrBytes}`);
    const reqStrEncr = await encrypt(resolverPubKey, reqStrBytes)
    const reqStrBytesEncrypted = new Uint8Array(reqStrEncr)
    return reqStrBytesEncrypted;
  }


  async execute(req: JsonRequest): Promise<JsonResponse> {
    this.log("executing request")
    // Wrap it in ProtoBuf

    const reqStrBytesEncrypted = await this.encryptRequest(req, this.networkParams.resolverPubKey);
    console.log("request encrypted");

    const privKey = generateKeyPair();

    const dappPubKeyBytes = privKey.publicKey.toBytes(true);
    const protoReq = create(ResolverRequestSchema, {
      id: req.Id,
      payload: reqStrBytesEncrypted,
      encrypted: true,
      publicKey: dappPubKeyBytes
    });
    this.log(`protoReq: ${protoReq}`);

    console.log(`create IncomingMessage with resolver key: ${this.networkParams.resolverPubKey}`);
    const incomingMsg = create(IncomingMessageSchema, {
      publicKeys: [ecies.PublicKey.fromHex(this.networkParams.resolverPubKey).toBytes(true)],
      request: protoReq,
    });
    console.log(`incomingMsg: ${JSON.stringify(incomingMsg)}`);
    // Try unwrap
    //
    const reqJson = toBinary(IncomingMessageSchema, incomingMsg);

    const extractedMsg = fromBinary(IncomingMessageSchema, reqJson);
    console.log(`requestId: ${extractedMsg.request?.id}`);
    console.log(`publicKeys: ${new TextDecoder().decode(extractedMsg.publicKeys[0])}`);


    this.log(`reqJson: ${reqJson}`);
    this.sendChannel.send(reqJson);

    let resolve, reject;
    const promise = new Promise<JsonResponse>((res, rej) => {
      resolve = res;
      reject = rej;
    });

    this.log(`pending id: ${req.Id}`)
    this.pendingRequests.set(req.Id, {resolve, reject, privKey});
    return promise;
  }

  async onmessage(ev: MessageEvent)  {
    const data = ev.data;
    console.log(`onmessage data: ${ev.data}`);
    console.log(`onmessage type: ${typeof(ev.data)}`);
    const bytes = new Uint8Array(data);
    const outgoingMsg = fromBinary(OutgoingMessageSchema, bytes);
    const protoResp = outgoingMsg.result;
    console.log(`chan msg: ${JSON.stringify(protoResp)}`)

    if (protoResp.case == "error" || protoResp.case == undefined) {
      console.log("error in response", protoResp.value);
      return;
    }
    let pendingReq = this.pendingRequests.get(protoResp.value.id);
    if (pendingReq == null) {
      return
    }
    let { resolve, reject, privKey } = pendingReq;
    const privKeyHex = privKey.toHex();
    const payload = await decrypt(privKeyHex, protoResp.value.payload)
    const resp: JsonResponse = JSON.parse(payload);
    console.log(`chan msg result: ${JSON.stringify(resp)}`)
    
    console.log(`resolve id: ${resp.id}`)
    resolve(resp);
    this.pendingRequests.delete(resp.id);
  }

  async onnegotiationneeded() {
    try {
      this.makingOffer = true;
      await this.pc.setLocalDescription();

      let sessionAndOffer = {'session_id': 'firefox', 'offer': this.pc.localDescription}
      let resp = await this.send('/sdp', sessionAndOffer)
      this.log(`resp: ${JSON.stringify(resp.data)}`)
      this.pc.setRemoteDescription(resp.data.answer)
    } catch (err) {
      console.error(err);
    } finally {
      this.makingOffer = false;
    }
  }

  async fetchNetworkParams(clientParams: ClientParams): Promise<NetworkParams> {
    const client = createPublicClient({ transport: http(clientParams.providerUrl) });
    const data = await client.readContract({
      address: clientParams.contractAddr,
      abi: registryAbi,
      functionName: 'getRelayer',
    });
    console.log(`registry res: ${JSON.stringify(data)}`);
    return {relayerIp: data[0], resolverPubKey: data[1][0]};
  }
};



const defaultStunServers = ['stun:stun.l.google.com:19302', 'stun:stun.services.mozilla.com'];
const defaultChannelName = 'default';

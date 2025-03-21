import { ClientParams, JsonRequest, JsonResponse, NetworkParams, Logger, PendingRequest } from "./types";
import axios from 'axios';
import * as ecies from "eciesjs";
import { generateKeyPair, encrypt, decrypt } from "./crypto/util";
import { Error as ResolverError, ResolverRequestSchema, ResolverResponse } from "./gen/resolver_pb";
import { IncomingMessageSchema, OutgoingMessageSchema } from "./gen/relayer_pb";
import { Address, createPublicClient, http } from 'viem'
import { registryAbi } from "./abi/NodeRegistry";
import { create, toJson, toJsonString, toBinary, fromBinary, fromJsonString} from "@bufbuild/protobuf";

export class Client {
  pc: RTCPeerConnection | null;
  sendChannel: RTCDataChannel | null;
  makingOffer: boolean = false;
  networkParams: NetworkParams | null;
  connectionOpened: any;
  connectionClosed: any;
  pendingRequests: Map<string, PendingRequest>;
  logger: Logger;

  constructor(logger: Logger) {
    this.pc = null;
    this.sendChannel = null;
    this.makingOffer = false;
    this.networkParams = null;
    this.pendingRequests = new Map<string, PendingRequest>();
    this.logger = logger;
  }

  async init(params: ClientParams) {
    this.logger.info("Initializing client with params:", params);
    const networkParams = await this.fetchNetworkParams(params);
    this.networkParams = networkParams;
    this.logger.debug("Network parameters details:", JSON.stringify(networkParams));

    const pc = new RTCPeerConnection({ iceServers: [{ urls: defaultStunServers }] });
    this.pc = pc;
    this.logger.info("RTCPeerConnection created");

    pc.onnegotiationneeded = () => {
      this.logger.info("Negotiation needed");
      this.onnegotiationneeded();
    };

    pc.onicecandidate = ({ candidate }) => {
      this.logger.debug("ICE candidate received");
      if (candidate) {
        this.logger.debug("Raw ICE candidate string:", candidate.candidate);
      }
      this.onicecandidate(candidate);
    };

    const sendChannel = pc.createDataChannel(defaultChannelName);
    this.sendChannel = sendChannel;
    this.logger.info("Data channel created with name:", defaultChannelName);

    sendChannel.onmessage = (ev) => {
      this.logger.info("Data channel message received");
      this.onmessage(ev);
    };
    sendChannel.onopen = () => {
      this.logger.info("Data channel opened");
      this.onopen();
    };
    sendChannel.onclose = () => {
      this.logger.info("Data channel closed");
      this.onclose();
    };

    return new Promise<boolean>((res, rej) => {
      this.connectionOpened = res;
      this.connectionClosed = rej;
    });
  }

  onopen() {
      this.logger.info("channel open")
      this.connectionOpened(true);
  }

  onclose() {
    this.logger.info("channel closed")
    this.connectionClosed(true);
  }


  onicecandidate(candidate: RTCIceCandidate | null) {
    if (candidate !== null) {
      const sessionAndCandidate = {
        session_id: "firefox",
        candidate: parsePionCandidate(candidate.candidate),
      };
      this.logger.debug(`Candidate processed: ${JSON.stringify(sessionAndCandidate)}`);
      this.send("/candidate", sessionAndCandidate);
    }
  }

  send(url: string, msg: { session_id: string; candidate?: any; offer?: RTCSessionDescription | null; }) {
    const headers = { "Content-Type": "application/json" };
    const addr = "http://" + (this.networkParams?.relayerIp || "") + url;
    this.logger.debug(`Send http POST request: ${addr}`);
    return axios.post(addr, msg, { headers });
  }

  async encryptRequest(req: JsonRequest, resolverPubKey: string) {
    const reqStr = JSON.stringify(req);
    this.logger.debug("Request string:", reqStr);
    const reqStrBytes = new TextEncoder().encode(reqStr);
    this.logger.debug("Request string bytes:", reqStrBytes);
    const reqStrEncr = await encrypt(resolverPubKey, reqStrBytes);
    const encrypted = new Uint8Array(reqStrEncr);
    this.logger.debug("Encrypted request bytes:", encrypted);
    return encrypted;
  }


  async execute(req: JsonRequest, shouldEncrypt: boolean = true): Promise<JsonResponse> {
    this.logger.info("Executing request");
    const resolverPubKey = this.networkParams?.resolverPubKey || "";
    let payloadBytes: Uint8Array;
    if (shouldEncrypt) {
      payloadBytes = await this.encryptRequest(req, resolverPubKey);
      this.logger.debug("Request encryption complete");
    } else {
      const reqStr = JSON.stringify(req);
      this.logger.debug("Request string (unencrypted):", reqStr);
      payloadBytes = new TextEncoder().encode(reqStr);
    }

    const privKey = generateKeyPair();
    this.logger.debug("Generated key pair:", { publicKey: privKey.publicKey.toBytes(true) });
    const dappPubKeyBytes = privKey.publicKey.toBytes(true);
    const protoReq = create(ResolverRequestSchema, {
      id: req.Id,
      payload: payloadBytes,
      encrypted: shouldEncrypt,
      publicKey: dappPubKeyBytes,
    });
    this.logger.debug("ProtoReq constructed:", JSON.stringify(protoReq));

    this.logger.info(`Creating IncomingMessage with resolver key: ${resolverPubKey}`);
    const incomingMsg = create(IncomingMessageSchema, {
      publicKeys: [ecies.PublicKey.fromHex(resolverPubKey).toBytes(true)],
      request: protoReq,
    });
    this.logger.debug("IncomingMsg created:", JSON.stringify(incomingMsg));

    const reqJson = toBinary(IncomingMessageSchema, incomingMsg);
    const extractedMsg = fromBinary(IncomingMessageSchema, reqJson);
    this.logger.debug(`Extracted RequestId: ${extractedMsg.request?.id}`);
    if (extractedMsg.publicKeys && extractedMsg.publicKeys.length > 0) {
      this.logger.debug(`Extracted PublicKey (decoded): ${new TextDecoder().decode(extractedMsg.publicKeys[0])}`);
    }
    this.logger.debug("Binary request (reqJson):", reqJson);

    this.sendChannel?.send(reqJson);

    let resolve, reject;
    const promise = new Promise<JsonResponse>((res, rej) => {
      resolve = res;
      reject = rej;
    });

    this.logger.info(`Pending request id: ${req.Id}`);
    this.pendingRequests.set(req.Id, { resolve, reject, privKey });
    return promise;
  }

  async onmessage(ev: MessageEvent)  {
    const data = ev.data;
    this.logger.info("onmessage data received");
    this.logger.debug("onmessage raw data:", data);

    const bytes = new Uint8Array(data);
    const outgoingMsg = fromBinary(OutgoingMessageSchema, bytes);
    const protoResp = outgoingMsg.result;
    this.logger.info("Channel message received");
    this.logger.debug("Channel message details:", JSON.stringify(protoResp));

    if (protoResp.case === "error" || protoResp.case === undefined) {
      this.logger.error("Error in response", JSON.stringify(protoResp.value));
      const errorObj =
        typeof protoResp.value === "object" && protoResp.value !== null
          ? protoResp.value
          : { message: "Unknown error in response" };
      const errorMsg = errorObj.message || "Unknown error in response";

      if ("id" in errorObj) {
        const errorId = (errorObj as { id: string }).id;
        const pendingReq = this.pendingRequests.get(errorId);
        if (pendingReq) {
          pendingReq.reject(new Error(errorMsg));
          this.pendingRequests.delete(errorId);
          return;
        }
      }
      return;
    }

    this.logger.debug("Response without error");
    if (!protoResp.value || typeof protoResp.value !== "object" || !("id" in protoResp.value)) {
      this.logger.warn("Invalid response: missing id", protoResp.value);
      return;
    }
    const successResp = protoResp.value as ResolverResponse;
    const pendingReq = this.pendingRequests.get(successResp.id);
    if (!pendingReq) {
      this.logger.warn(`No pending request found for response id: ${successResp.id}`);
      return;
    }
    const { resolve, reject, privKey } = pendingReq;
    const privKeyHex = privKey.toHex();
    
    // If ResolverResponse has some a error, reject pending request and give away this error
    if (successResp.result.case === "error" || successResp.result.case === undefined) {
      const error = successResp.result.value as ResolverError

      this.logger.error(`Received a response with an error on the 'Resolver', error message: ${error.message}, error code: ${error.code}`)
      pendingReq.reject(new Error(`Received a response with an error on the 'Resolver': ${error.message}`))
      this.pendingRequests.delete(successResp.id)
      return
    }

    const responseValue = successResp.result.value;
    if (successResp.encrypted || !this.tryParse(responseValue)) {
      try {   
        const payload = await decrypt(privKeyHex, responseValue);
        const resp: JsonResponse = JSON.parse(payload);
        this.logger.info("Channel message result processed (decrypted)");
        this.logger.debug("Processed response:", JSON.stringify(resp));
        this.logger.info(`Resolving request with id: ${resp.id}`);
        resolve(resp);
        this.pendingRequests.delete(resp.id);
      } catch (decryptionError) {
        this.logger.error("Error processing (decrypting) response:", decryptionError);
        reject(new Error("Failed to process response: " + decryptionError));
        this.pendingRequests.delete(successResp.id);
      }
    } else {
      try {
        const payload = new TextDecoder().decode(responseValue);
        const resp: JsonResponse = JSON.parse(payload);
        this.logger.info("Channel message result processed (unencrypted)");
        this.logger.debug("Processed response:", JSON.stringify(resp));
        this.logger.info(`Resolving request with id: ${resp.id}`);
        resolve(resp);
        this.pendingRequests.delete(resp.id);
      } catch (error) {
        this.logger.error("Error processing (unencrypted) response:", error);
        reject(new Error("Failed to process response (unencrypted): " + error));
        this.pendingRequests.delete(successResp.id);
      }
    }
  }

  tryParse(payload: any): boolean {
    try {
      JSON.parse(new TextDecoder().decode(payload));
      return true
    }
    catch {
      return false
    }
  }

  async onnegotiationneeded() {
    try {
      this.makingOffer = true;
      await this.pc?.setLocalDescription();
      const sessionAndOffer = { session_id: "firefox", offer: this.pc?.localDescription };
      const resp = await this.send("/sdp", sessionAndOffer);
      this.logger.info(`Response from SDP received`);
      this.logger.debug(`SDP response data: ${JSON.stringify(resp.data)}`);
      this.pc?.setRemoteDescription(resp.data.answer);
    } catch (err) {
      this.logger.error("Error in negotiation needed:", err);
    } finally {
      this.makingOffer = false;
    }
  }

  async fetchNetworkParams(clientParams: ClientParams): Promise<NetworkParams> {
    const client = createPublicClient({ transport: http(clientParams.providerUrl) });
    const data: any = await client.readContract({
      address: clientParams.contractAddr as Address,
      abi: registryAbi,
      functionName: "getRelayer",
    });
    return { relayerIp: data[0] as string, resolverPubKey: data[1][0] };
  }
};

function parsePionCandidate(candidateLine: string) {
  const parts = candidateLine.trim().split(/\s+/);
  if (!parts[0].startsWith("candidate:") || parts.length < 8) {
    throw new Error("Invalid ICE candidate string format: " + candidateLine);
  }
  const foundation = parts[0].split(":")[1];
  const component = parseInt(parts[1], 10);
  const protocol = parts[2];
  const priority = parseInt(parts[3], 10);
  const address = parts[4];
  const port = parseInt(parts[5], 10);
  if (parts[6] !== "typ") {
    throw new Error(`Expected 'typ' at position 6 but got: ${parts[6]}`);
  }
  const candidateType = parts[7];

  let relatedAddress = "";
  let relatedPort = 0;
  let tcpType = "";
  let i = 8;
  while (i < parts.length) {
    switch (parts[i]) {
      case "raddr":
        relatedAddress = parts[i + 1] || "";
        i += 2;
        break;
      case "rport":
        relatedPort = parseInt(parts[i + 1], 10) || 0;
        i += 2;
        break;
      case "tcptype":
        tcpType = parts[i + 1] || "";
        i += 2;
        break;
      default:
        i += 1;
    }
  }

  return {
    foundation,
    priority,
    address,
    protocol: parseProtocol(protocol),
    port,
    type: candidateType,
    component,
    relatedAddress,
    relatedPort,
    tcpType,
  };
}

function parseProtocol(protocolStr: string) {
  switch (protocolStr) {
    case "udp":
      return 1;
    case "tcp":
      return 2;
    default:
      return 0;
  }
}

const defaultStunServers = ['stun:stun.l.google.com:19302', 'stun:stun.services.mozilla.com'];
const defaultChannelName = 'default';

import * as ecies from "eciesjs";
import { Buffer } from 'buffer';

globalThis.Buffer = Buffer;

function base64StringToUint8Array(b64str: string) {
  var byteStr = atob(b64str);
  var bytes = new Uint8Array(byteStr.length);
  for (var i = 0; i < byteStr.length; i++) {
    bytes[i] = byteStr.charCodeAt(i);
  }
  return bytes;
}

export async function encrypt(pubKey: string, payload: any) {
  // // Compress
  // const resolverPubKeyCompressedBytes = new ecies.PublicKey(keyBytes).toBytes(true);
  // console.log("Encrypt with resolverPubKeyCompressedBytes");
  return ecies.encrypt(pubKey, payload)
}

export async function decrypt(privKey: any, payload: Uint8Array) {
  let decrypted = ecies.decrypt(privKey, payload);
  return decrypted.toString();
}

export function generateKeyPair() {
  const privKey = new ecies.PrivateKey();
  return privKey;
}

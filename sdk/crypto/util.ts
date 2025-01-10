import * as ecies from "eciesjs";
import { Buffer } from 'buffer';

globalThis.Buffer = Buffer;

export enum KeyType {
  Secp256k1 = "ECDSA",
  RSA4096 = "RSA",
};

const rsaAlgorithm = {
  name: "RSA-OAEP",
  modulusLength: 4096,
  publicExponent: new Uint8Array([1, 0, 1]),
  extractable: false,
  hash: {
    name: "SHA-256"
  }
};

// https://gist.github.com/mholt/813db71291de8e45371fe7c4749df99c
function pemEncode(label, data) {
	const base64encodedWrapped = data.replace(/(.{64})/g, "$1\n");
	return `-----BEGIN ${label}-----\n${base64encodedWrapped}\n-----END ${label}-----`;
}

export async function pemEncodedPublicKey(keyType: string, keyPair: any) {
  if (keyType == KeyType.Secp256k1) {
    const exported = keyPair.publicKey.toBytes(false);
    const str = exported.toBase64();
    return pemEncode("ECDSA PUBLIC KEY", str);
  } else if (keyType == KeyType.RSA4096) {
    const exported = await window.crypto.subtle.exportKey("spki", keyPair.publicKey);
    const str = new Uint8Array(exported).toBase64();
    return pemEncode("RSA PUBLIC KEY", str);
  } else {
    return null;
  }
}

function base64StringToUint8Array(b64str) {
  var byteStr = atob(b64str);
  var bytes = new Uint8Array(byteStr.length);
  for (var i = 0; i < byteStr.length; i++) {
    bytes[i] = byteStr.charCodeAt(i);
  }
  return bytes;
}

export function getKeyType(pemPubKey: string) {
  if (pemPubKey == null) {
    return null;
  }
  const typeSuffix = " PUBLIC KEY";
  let keyType = pemPubKey.substring(11, pemPubKey.indexOf(typeSuffix));
  return keyType;
}

function parsePem(pemPubKey: string) {
  const keyType = getKeyType(pemPubKey);
  let pemStart = "-----BEGIN ";
  pemStart = pemStart.concat(keyType).concat(" PUBLIC KEY-----");
  let pemEnd = "-----END ";
  pemEnd = pemEnd.concat(keyType).concat(" PUBLIC KEY-----");
  const keyStr = pemPubKey.replace(pemStart, "").replace(pemEnd, "").trim();
  return {
    keyType: keyType,
    keyBytes: base64StringToUint8Array(keyStr),
  }
}

export async function encrypt(pemPubKey: string, payload: any) {
  const {keyType, keyBytes} = parsePem(pemPubKey);

  if (keyType == KeyType.RSA4096) {
    let importedKey = await crypto.subtle.importKey("spki", keyBytes, rsaAlgorithm, true, ["encrypt"]);
    return crypto.subtle.encrypt(
      {
        name: "RSA-OAEP",
      },
      importedKey,
      payload,
    );
  }
  else if (keyType == KeyType.Secp256k1) {
    return ecies.encrypt(keyBytes, payload)
  }
}

export async function decrypt(keyType: any, keyPair: any, payload: string) {
  if (keyType == KeyType.Secp256k1) {
    let decrypted = ecies.decrypt(keyPair.toHex(), payload);
    return decrypted.toString();
  } else if (keyType == KeyType.RSA4096) {
    let res = await crypto.subtle.decrypt(
      rsaAlgorithm,
      keyPair.privateKey,
      payload,
    )
    return new TextDecoder().decode(res);
  } else {
    return null;
  }
}

export async function generateKeyPair(type: KeyType) {
  if (type == KeyType.Secp256k1) {
    const privKey = new ecies.PrivateKey();
    return privKey;
  } else if (type == KeyType.RSA4096) {
    let keypair = await window.crypto.subtle.generateKey(
      rsaAlgorithm,
      true,
      ["encrypt", "decrypt"],
    );
    return keypair;
  } else {
    return null;
  }
}

import axios from 'axios';
import {connect, execute, log, JsonRequest, JsonResponse, ConnParams} from './webrtc'

import OpenCrypto from 'opencrypto'

// Initialize new OpenCrypto instance
const crypt = new OpenCrypto()


function setup() {
  const connParams: ConnParams = {
    relayerAddress: 'http://localhost:8080',
    channelName: 'default',
    stunServers: ['stun:stun.l.google.com:19302', 'stun:stun.services.mozilla.com']
  }
  
  connect(connParams)
}

async function testExecute() {
  let req: JsonRequest = {
    Id: "TestID",
    Method: "GetWalletBalance",
    Params: ["0x38308C349fd2F9dad31Aa3bFe28015dA3EB67193", "latest"]
  }

  let resp: JsonResponse = await execute(req, null);
  log(`resp received: ${JSON.stringify(resp)}`);
}

async function testExecuteEncrypted() {
  let req: JsonRequest = {
    Id: "TestID",
    Method: "GetWalletBalance",
    Params: ["0x38308C349fd2F9dad31Aa3bFe28015dA3EB67193", "latest"]
  }
  let pemPubKey = document.getElementById("resolverPem").value

  let resp: JsonResponse = await execute(req, pemPubKey);
  log(`resp received: ${JSON.stringify(resp)}`);
}

window.onload = (ev) => {
  setup()
  document.getElementById('test_execute').onclick = testExecute
  document.getElementById('test_execute_encr').onclick = testExecuteEncrypted
}

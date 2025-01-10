import { JsonRequest, JsonResponse } from './types'
import { Client } from "./client";

async function testExecute() {

  let client = new Client();
  await client.init({
    providerUrl: "http://localhost:8545",
    contractAddr: "0x5FbDB2315678afecb367f032d93F642f64180aa3"
  });
  console.log("WebRTC initialized");
  const req: JsonRequest = {
    Id: "TestID",
    Method: "GetWalletBalance",
    Params: ["0x38308C349fd2F9dad31Aa3bFe28015dA3EB67193", "latest"]
  };

  let resp: JsonResponse = await client.execute(req);
  client.log(`resp received: ${JSON.stringify(resp)}`);
}

window.onload = async (ev) => {
  document.getElementById('test_execute').onclick = testExecute
}

import { JsonRequest, JsonResponse } from './types'
import { Client } from "./client";


// I added tags for input parameters for request. It is need for test. 
// In test you can check expected and actual result without hardcode
// I left default volumes in input tags.
async function callExecute() {
  let client = new Client();
  await client.init({
    providerUrl: "http://localhost:8545",
    contractAddr: "0x5FbDB2315678afecb367f032d93F642f64180aa3"
  });
  console.log("WebRTC initialized");
  const req: JsonRequest = {
    Id: getRequestIdFromInput(),
    Method: getMethodFromInput(),
    Params: getParamsFromInput()
  };

  let resp: JsonResponse = await client.execute(req);
  setResponseToDoc(resp)
}

window.onload = async (ev) => {
  document.getElementById('button-test-execute').onclick = callExecute
}

function getRequestIdFromInput(): string {
  return document.getElementById('input-request-id')?.value;
}

function getMethodFromInput(): string {
  return document.getElementById("input-method-name")?.value;
}

function getParamsFromInput(): string[] {
  let paramsStr = document.getElementById("input-params")?.value
  
  return paramsStr.split(",")
}

function setResponseToDoc(resp: JsonResponse) {
  let inputRequestIdResult = createNewInputElement("input-request-id-result", resp.id)
  let inputResult = createNewInputElement("input-result", resp.result)
  document.getElementById("td-for-request-id-result")?.appendChild(inputRequestIdResult)
  document.getElementById("td-for-result")?.appendChild(inputResult)
}

function createNewInputElement(id: string, value: string) {
  let newInputElement = document.createElement("input")
  newInputElement.type = "text"
  newInputElement.id = id
  newInputElement.size = 60
  newInputElement.value = value

  return newInputElement
}

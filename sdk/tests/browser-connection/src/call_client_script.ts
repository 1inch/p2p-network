import { JsonRequest, JsonResponse } from '../../../types'
import { Client } from "../../../client";


const defaultLogger = {
  debug: (...args: any[]) => console.log(`DEBUG, ${args.join(" ")}`),
  info: (...args: any[]) => console.log(`INFO ${args.join(" ")}`),
  warn: (...args: any[]) => console.log(`WARNING ${args.join(" ")}`),
  error: (...args: any[]) => console.log(`ERROR ${args.join(" ")}`)
};

defaultLogger.info("script loaded");

var client = new Client(defaultLogger);

function initDatachannel() {
  client.init({
    providerUrl: "http://localhost:8545",
    contractAddr: "0x5FbDB2315678afecb367f032d93F642f64180aa3"
  })
  .then((ok)=> {
    if (ok) {
      inputDatachannelState(ok, "")
    }
  })
}

// I added tags for input parameters for request. It is need for test. 
// In test you can check expected and actual result without hardcode
// I left default volumes in input tags.
async function callExecute() {
  deleteResponseInputs();
  
  const req: JsonRequest = {
    Id: getRequestIdFromInput(),
    Method: getMethodFromInput(),
    Params: getParamsFromInput()
  };

  // TODO fix this sleep
  // tthis pause is used to wait for the CI server to be unloaded
  await new Promise(f => setTimeout(f, 5000))

  client.execute(req)
  .then(resp => setResponseToDoc(resp))
  .catch(err => setErrorToDoc(err))
}

window.onload = async (ev) => {
  document.getElementById('button-test-execute')!.onclick = callExecute
  document.getElementById('button-init-datachannel')!.onclick = initDatachannel
}

function getRequestIdFromInput(): string {
  return getValueFromForm('input-request-id');
}

function getMethodFromInput(): string {
  return getValueFromForm("input-method-name");
}

function getParamsFromInput(): string[] {
  let paramsStr = getValueFromForm("input-params")
  return paramsStr.split(",")
}

function getValueFromForm(formId: string): string {
  let form = document.getElementById(formId)
  if (form) 
    return (form as HTMLFormElement).value
  throw new Error(`Form with Id: ${formId}, not found on test intex.html`)
}

function setResponseToDoc(resp: JsonResponse) {
  setRequestIdAndResultToInputs(resp.id, resp.result);
}

function setRequestIdAndResultToInputs(requestId: string, result: string) {
  let inputRequestIdResult = createNewInputElement("input-request-id-result", requestId);
  let inputResult = createNewInputElement("input-result", result);
  document.getElementById("td-for-request-id-result")?.appendChild(inputRequestIdResult);
  document.getElementById("td-for-result")?.appendChild(inputResult);
}

function setErrorToDoc(err: Error) {
  setRequestIdAndResultToInputs("nil", err.message);
}

function createNewInputElement(id: string, value: string): HTMLInputElement {
  let newInputElement = document.createElement("input")
  newInputElement.type = "text"
  newInputElement.id = id
  newInputElement.size = 60
  newInputElement.value = value

  return newInputElement
}

function createNewTextElement(id: string, value: string): HTMLLabelElement {
  let newLabelElement = document.createElement("label")
  newLabelElement.id = id
  newLabelElement.textContent = value

  return newLabelElement
}

function deleteResponseInputs() {
  let inputForResult = document.getElementById("input-result");

  if (inputForResult) {
    console.log("remove input for result");
    inputForResult.remove();
  }

  let inputForRequestId = document.getElementById("input-request-id-result");

  if (inputForRequestId) {
    console.log("remove input for request id");
    inputForRequestId.remove();
  }
}

function inputDatachannelState(status: boolean, reason: string) {
  let labelForTextConnected = createNewTextElement("label-text-connected", "Connected: ")
  let labelForState = createNewTextElement("label-status", status.toString())

  let blockForState = document.getElementById("div-state")

  blockForState?.appendChild(labelForTextConnected)
  blockForState?.appendChild(labelForState)
}

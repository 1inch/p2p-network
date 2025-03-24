import { JsonRequest, JsonResponse } from '../../../types'
import { Client } from "../../../client";


const defaultLogger = {
  debug: (...args: any[]) => console.log(`DEBUG, ${args.join(" ")}`),
  info: (...args: any[]) => console.log(`INFO ${args.join(" ")}`),
  warn: (...args: any[]) => console.log(`WARNING ${args.join(" ")}`),
  error: (...args: any[]) => console.log(`ERROR ${args.join(" ")}`)
};

defaultLogger.info("script loaded");

// I added tags for input parameters for request. It is need for test. 
// In test you can check expected and actual result without hardcode
// I left default volumes in input tags.
async function callExecute() {
  
  let client = new Client(defaultLogger);
  await client.init({
    providerUrl: "http://localhost:8545",
    contractAddr: "0x5FbDB2315678afecb367f032d93F642f64180aa3"
  });
  defaultLogger.info("WebRTC initialized");
  const req: JsonRequest = {
    Id: getRequestIdFromInput(),
    Method: getMethodFromInput(),
    Params: getParamsFromInput()
  };

  client.execute(req)
  .then(resp => setResponseToDoc(resp))
  .catch(err => setErrorToDoc(err))
}

window.onload = async (ev) => {
  document.getElementById('button-test-execute')!.onclick = callExecute
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
  let inputRequestIdResult = createNewInputElement("input-request-id-result", resp.id)
  let inputResult = createNewInputElement("input-result", resp.result)
  document.getElementById("td-for-request-id-result")?.appendChild(inputRequestIdResult)
  document.getElementById("td-for-result")?.appendChild(inputResult)
}

function setErrorToDoc(err: Error) {
  let inputResult = createNewInputElement("input-result", err.message)
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

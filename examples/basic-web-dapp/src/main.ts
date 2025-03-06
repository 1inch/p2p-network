import { JsonRequest, JsonResponse } from '../../../sdk/types';
import { Client } from '../../../sdk/client';

// Simple logger function to append messages to the UI logs area
function logMessage(msg: string): void {
  const logsDiv = document.getElementById("logs");
  if (logsDiv) {
    logsDiv.innerHTML += msg + "<br />";
    logsDiv.scrollTop = logsDiv.scrollHeight;
  } else {
    console.log(msg);
  }
}

async function executeRequest() {
  // Create a new WebRTC client instance with signaling server settings
  const client = new Client();

  // Initialize the client with blockchain parameters (replace with actual values)
  await client.init({
    providerUrl: "http://localhost:8545",
    contractAddr: "0x5FbDB2315678afecb367f032d93F642f64180aa3"
  });
  logMessage("WebRTC initialized");

  // Build the request from the UI inputs
  const req: JsonRequest = {
    Id: "TestID",
    Method: "GetWalletBalance",
    Params: ["0x38308C349fd2F9dad31Aa3bFe28015dA3EB67193", "latest"]
  };

  try {
    logMessage("Sending request...");
    const resp: JsonResponse = await client.execute(req);
    logMessage("Response received: " + JSON.stringify(resp));
  } catch (error) {
    logMessage("Error executing request: " + error);
  }
}

// Attach event handler when the window loads
window.onload = () => {
  const btn = document.getElementById("sendBtn");
  if (btn) {
    btn.onclick = executeRequest;
  }
};

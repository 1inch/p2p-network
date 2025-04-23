import { JsonRequest, JsonResponse } from '../../../types';
import { Client } from '../../../client';

// Provide node endpoint and discovery contract address
const providerUrl = "http://blockchain:8545";
const contractAddr = "0x5FbDB2315678afecb367f032d93F642f64180aa3";

function appendLog(level: string, message: string): void {
  const logsDiv = document.getElementById("logs");
  if (logsDiv) {
    const p = document.createElement("p");
    switch (level) {
      case "DEBUG":
        p.className = "text-muted";
        break;
      case "INFO":
        p.className = "text-secondary";
        break;
      case "WARN":
        p.className = "text-warning";
        break;
      case "ERROR":
        p.className = "text-danger";
        break;
      default:
        p.className = "text-dark";
    }
    p.innerHTML = message;
    logsDiv.appendChild(p);
    logsDiv.scrollTop = logsDiv.scrollHeight;
  }
}

function showError(message: string): void {
  const alertContainer = document.getElementById("alertContainer");
  if (alertContainer) {
    const alertDiv = document.createElement("div");
    alertDiv.className = "alert alert-danger alert-dismissible fade show";
    alertDiv.role = "alert";
    alertDiv.innerHTML =
      message +
      '<button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>';
    alertContainer.appendChild(alertDiv);
  }
}

const baseLogger = {
  debug: (...args: any[]) => {
    const msg = args.map(arg => (typeof arg === "object" ? JSON.stringify(arg) : arg)).join(" ");
    appendLog("DEBUG", msg);
  },
  info: (...args: any[]) => {
    const msg = args.map(arg => (typeof arg === "object" ? JSON.stringify(arg) : arg)).join(" ");
    appendLog("INFO", msg);
  },
  warn: (...args: any[]) => {
    const msg = args.map(arg => (typeof arg === "object" ? JSON.stringify(arg) : arg)).join(" ");
    appendLog("WARN", msg);
  },
  error: (...args: any[]) => {
    const msg = args.map(arg => (typeof arg === "object" ? JSON.stringify(arg) : arg)).join(" ");
    appendLog("ERROR", msg);
    showError(msg);
  }
};

const mainLogger = {
  debug: (...args: any[]) => baseLogger.debug(`[MAIN] DEBUG ${args.join(" ")}`),
  info: (...args: any[]) => baseLogger.info(`[MAIN] INFO ${args.join(" ")}`),
  warn: (...args: any[]) => baseLogger.warn(`[MAIN] WARNING ${args.join(" ")}`),
  error: (...args: any[]) => baseLogger.error(`[MAIN] ERROR ${args.join(" ")}`)
};

// provide logger to the SDK (in production enforce log level)
const clientLogger = {
  debug: (...args: any[]) => baseLogger.debug(`[CLIENT] DEBUG ${args.join(" ")}`),
  info: (...args: any[]) => baseLogger.info(`[CLIENT] INFO ${args.join(" ")}`),
  warn: (...args: any[]) => baseLogger.warn(`[CLIENT] WARNING ${args.join(" ")}`),
  error: (...args: any[]) => baseLogger.error(`[CLIENT] ERROR ${args.join(" ")}`)
};

function updateConnectionState(state: string): void {
  const connectionStatusEl = document.getElementById("connectionStatus");
  const connectionStateEl = document.getElementById("connectionState");
  if (connectionStatusEl) connectionStatusEl.textContent = state;
  if (connectionStateEl) {
    if (state === "Connected") {
      connectionStateEl.className = "alert alert-success";
    } else if (state === "Disconnected") {
      connectionStateEl.className = "alert alert-danger";
    } else if (state === "Connecting") {
      connectionStateEl.className = "alert alert-secondary";
    } else {
      connectionStateEl.className = "alert alert-secondary";
    }
  }
}

async function initializeClient(): Promise<Client> {
  const client = new Client(clientLogger);
  try {
    await client.init({
      providerUrl: providerUrl,
      contractAddr: contractAddr
    });
    updateConnectionState("Connected");
    mainLogger.info("WebRTC initialized");
  } catch (err) {
    updateConnectionState("Disconnected");
    mainLogger.error("Failed to initialize WebRTC: " + err);
    showError("Failed to initialize WebRTC: " + err);
    throw err;
  }
  return client;
}

function shouldEncryptRequest(): boolean {
  const encryptCheckbox = document.getElementById("encryptCheckbox") as HTMLInputElement;
  return encryptCheckbox ? encryptCheckbox.checked : true;
}

async function getBalance(client: Client): Promise<void> {
  const chainId = (document.getElementById("chainIdInput") as HTMLInputElement).value;
  const address = (document.getElementById("addressInput") as HTMLInputElement).value;
  const balanceField = document.getElementById("balanceList") as HTMLInputElement;
  const req: JsonRequest = {
    Id: "TestID-GetBalance",
    Method: "GetWalletBalance",
    Params: [chainId, address]
  };
  try {
    mainLogger.info("Sending GetBalance request...");
    const resp: JsonResponse = await client.execute(req, shouldEncryptRequest());
    mainLogger.info("GetBalance response received:", JSON.stringify(resp));
    if (resp && (resp as any).result) {
      renderBalances((resp as any).result);
    } else {
      balanceField.value = "No balance returned";
    }
  } catch (error) {
    mainLogger.error("Error executing GetBalance: " + error);
  }
}

async function sendFunds(client: Client): Promise<void> {
  const req: JsonRequest = {
    Id: "TestID-SendFunds",
    Method: "SendFunds",
    Params: ["recipient-address", "100"]
  };
  try {
    mainLogger.info("Sending SendFunds request...");
    const resp: JsonResponse = await client.execute(req, shouldEncryptRequest());
    mainLogger.info("SendFunds response received:", JSON.stringify(resp));
    const resultField = document.getElementById("sendFundsResult") as HTMLInputElement;
    if (resp && (resp as any).result) {
      resultField.value = (resp as any).result.toString();
    } else {
      resultField.value = "No result returned";
    }
  } catch (error) {
    mainLogger.error("Error executing SendFunds: " + error);
    showError("Error executing SendFunds: " + error);
  }
}

window.onload = async () => {
  let client: Client;
  try {
    client = await initializeClient();
  } catch (e) {
    return;
  }
  
  const getBalanceBtn = document.getElementById("getBalanceBtn");
  if (getBalanceBtn) {
    getBalanceBtn.onclick = () => {
      getBalance(client);
    };
  }
  
  const sendFundsBtn = document.getElementById("sendFundsBtn");
  if (sendFundsBtn) {
    sendFundsBtn.onclick = () => {
      sendFunds(client)
    };
  }
};

const COMMON_TOKENS: Record<string, string> = {
  '0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee': 'ETH or NATIVE',   // native coin
  '0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48': 'USDC',
  '0xdac17f958d2ee523a2206206994597c13d831ec7': 'USDT',
} as const;

const toLower = (addr: string) => addr.toLowerCase();

function renderBalances(balances: Record<string, string>): void {
  let container = document.getElementById('balanceList') as HTMLDivElement | null;

  if (!container) {
    const fieldParent = document.getElementById('balanceField')!.parentElement!;
    container = document.createElement('div');
    container.id = 'balanceList';
    container.className = 'accordion-body mt-2';
    fieldParent.appendChild(container);
  }
  container.innerHTML = '';

  const sorted = Object.entries(balances).sort(([aToken, aAmt], [bToken, bAmt]) => {
    const aBalance = BigInt(aAmt);
    const bBalance = BigInt(bAmt);
    if (aBalance > bBalance) return -1;
    if (aBalance < bBalance) return 1;
    return 0;
  });

  sorted.forEach(([token, amount]) => {
    const symbol = COMMON_TOKENS[toLower(token)] ?? '';
    const row = document.createElement('div');
    row.className =
      'd-flex justify-content-between align-items-center border-bottom py-1 small';
    if (symbol) {
      row.innerHTML = `
        <span class="border rounded px-2 me-2">${symbol}</span>
        <span class="text-break flex-grow-1">${token}</span>
        <span class="fw-semibold">${amount}</span>
      `;
    } else {
      row.innerHTML = `
        <span class="text-break flex-grow-1">${token}</span>
        <span class="fw-semibold">${amount}</span>
      `;
    }

    container!.appendChild(row);
  });
}

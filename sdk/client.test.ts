import puppeteer, { Page } from 'puppeteer';

const webPageForTest = "http://localhost:3000/index.html"
const inputIdForRequestIdResult = "input-request-id-result"
const inputIdForResult = "input-result"


test('send GetWalletBalance request to relayer', async () => {
  let requestId = 'test-request-id';
  let methodName = 'GetWalletBalance';
  let params = ['0x38308C349fd2F9dad31Aa3bFe28015dA3EB67193', 'latest'];

  const browser = await puppeteer.launch();
  const page = await browser.newPage();

  await page.goto(webPageForTest);
  console.log("open page index.html");

  // Input values for request in inputs
  console.log("start input values to input for request")
  await inputValueInInputOnPage(page, 'input-request-id', requestId); 
  await inputValueInInputOnPage(page, 'input-method-name', methodName); 
  await inputValueInInputOnPage(page, 'input-params', mapParamsForInput(params));
  console.log("values successful inputed")

  // Find button test-exection and click
  console.log("click to button 'Test execute'")
  await page.click("#button-test-execute");

  await waitWhenRelayerGiveResponse(page)

  // Get values from response
  let actualRequestId = await getValueFromInputOnPage(page, inputIdForRequestIdResult);
  let actualResult = await getValueFromInputOnPage(page, inputIdForResult);


  expect(actualRequestId).toEqual(requestId);
  expect(actualResult).not.toBeUndefined();

  browser.close();
});

async function inputValueInInputOnPage(page: Page,  inputId: string, newValue: string) {
  await page.$eval(`#${inputId}`, (element, value) => element.value = value, newValue);
}

async function getValueFromInputOnPage(page: Page, inputId: string): Promise<string> {
  return await page.$eval(`#${inputId}`, element => element.value);
}

function mapParamsForInput(params: string[]): string {
  let separator = ",";
  return params.join(separator);
}

async function waitWhenRelayerGiveResponse(page: Page) {
  console.log("start wait response")
  await page.waitForSelector(`#${inputIdForRequestIdResult}`)
  await page.waitForSelector(`#${inputIdForResult}`)
  console.log("response received")
}

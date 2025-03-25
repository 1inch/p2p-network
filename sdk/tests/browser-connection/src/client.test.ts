import puppeteer, { Browser, Page } from 'puppeteer';

const webPageForTest = "http://localhost:3000/index.html"
const inputIdForRequestIdResult = "input-request-id-result"
const inputIdForResult = "input-result"


class Request {
  requestId: string
  methodName: string
  params: string[]

  constructor(requestId: string, methodName: string, params: string[]) {
    this.requestId = requestId
    this.methodName = methodName
    this.params = params
  } 
}

class Response {
  requestId: string
  result: any

  constructor(requestId: string, result: any) {
    this.requestId = requestId
    this.result = result
  }
}

type CheckingFunc = (request: Request, actualResponse: Response) => void

class TestCase {
  name: string
  request: Request
  checkingFunc: CheckingFunc

  constructor(name: string, request: Request, checkingFunc: CheckingFunc) {
    this.name = name
    this.request = request
    this.checkingFunc = checkingFunc
  }
}

describe("SDK integration tests",  ()=> {
  let browser: Browser
  // initialize test cases
  const testCases = [
    new TestCase(
      'send GetWalletBalance request to relayer',
      new Request(
        'positive-test-request-id', 
        'GetWalletBalance', 
        ['0x38308C349fd2F9dad31Aa3bFe28015dA3EB67193', 'latest']
      ),
      function(request: Request, actualResponse: Response) {
        expect(actualResponse.requestId).toEqual(request.requestId)
        expect(actualResponse.result).not.toBeUndefined()
      }
    ),
    new TestCase(
      'send request with unknown method name for api handler',
      new Request(
        'unknown-method-request-id',
        'GetBlockNumber',
        ['latest']
      ),
      function(request: Request, actualResponse: Response) {
        console.log(actualResponse.result)
        expect(actualResponse.result).toContain("unrecognized method")
      }
    ),
    new TestCase(
      'send request with differnet params',
      new Request(
        'different-params-request-id',
        'GetWalletBalance',
        ['latest']
      ),
      function(request: Request, actualResponse: Response) {
        console.log(actualResponse.result)
        expect(actualResponse.result).toContain("wrong number of params")
      }
    ),
    new TestCase(
      'send request with incorrect address',
      new Request(
        'incorrect-address-request-id',
        'GetWalletBalance',
        ['incorrect-address', 'latest']
      ),
      function(request: Request, actualResponse: Response) {
        console.log(actualResponse.result)
        expect(actualResponse.result).toContain("invalid format for address")
      }
    )
  ]

  // start browser before all tests
  beforeAll(async ()=> {
    browser = await launchBrowser()
  })

  // close browser after tests
  afterAll(async ()=> {
    browser.close();
  })

  // start test cases
  testCases.forEach((testCase) => {
    test(testCase.name, async ()=> {
      let page = await browser.newPage()
  
      await page.goto(webPageForTest);
      console.log("open page index.html")

      await inputRequestOnPage(page, testCase.request)

      await clickTestExecutionButton(page)
      await waitWhenRelayerGiveResponse(page)

      let actualResponse = await getResponseFromPage(page)
      console.log(`received response with requestId: ${actualResponse.requestId}`)
      testCase.checkingFunc(testCase.request, actualResponse)

      page.close()
    })
  })
})

// Find button test-exection and click
async function clickTestExecutionButton(page: Page) {
  console.log("click to button 'Test execute'")
  await page.click("#button-test-execute");
}

// Input values for request in inputs
async function inputRequestOnPage(page: Page, request: Request) {
  console.log("start input values to input for request")
  await inputValueInInputOnPage(page, 'input-request-id', request.requestId); 
  await inputValueInInputOnPage(page, 'input-method-name', request.methodName); 
  await inputValueInInputOnPage(page, 'input-params', mapParamsForInput(request.params));
  console.log("values successful inputed")
}

// Get values from response
async function getResponseFromPage(page: Page): Promise<Response> {
  let actualRequestId = await getValueFromInputOnPage(page, inputIdForRequestIdResult);
  let actualResult = await getValueFromInputOnPage(page, inputIdForResult);
  
  return new Response(actualRequestId, actualResult)
}

async function inputValueInInputOnPage(page: Page,  inputId: string, newValue: string) {
  await page.$eval(`#${inputId}`, (element, value) => (element as HTMLInputElement).value = value, newValue);
}

async function getValueFromInputOnPage(page: Page, inputId: string): Promise<string> {
  return await page.$eval(`#${inputId}`, element =>  (element as HTMLInputElement).value);
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

async function launchBrowser(): Promise<Browser>{
  return await puppeteer.launch({
    args: [
      '--no-sandbox',
      '--disable-setuid-sandbox'
    ]
  });
}

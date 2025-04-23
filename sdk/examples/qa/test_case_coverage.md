✅ Test case for local launch of Web dApp (Ethereum RPC)

* [QA] Test case coverage 1 — Checking the start of Anvil
  Check that Anvil launches without any bugs
STR
1.Add a command
docker run -p 8545:8545 --platform linux/amd64 ghcr.io/foundry-rs/foundry:latest "anvil --host 0.0.0.0"
2.Wait text like  Listening on 0.0.0.0:{8545}
Anvil starts successfully, without critical fixes
  Warning: This is a nightly build of Foundry. It is recommended to use the latest stable version. Visit https://book.getfoundry.sh/announcements for more information.
  To mute this warning set `FOUNDRY_DISABLE_NIGHTLY_WARNING` in your environment.
  
* [QA] Test case coverage 2 — Checking RPC availability
  Make sure RPC is running and responding
STR
Send a request for eth_chainId to http://localhost:8545
2.Check the response  
Expected Returns 0x7a69 (тобто Chain ID 31337) 

* [QA] Test case coverage 3 — Checking available accounts
  Verify that 10 accounts with balances have been created
STR
1.Сheck eth_accounts validation
2.Check that 10 addresses were eth
Expected 10 addresses, each with 10,000 ETH

  
* [QA] Test case coverage 4 — Verifying private keys
  Verify that private keys match addresses
  STR
1.Get private key from logs
2.Import into Metamask/web3
3.Check public address
Expected The address matches the key pair from the logs


* [QA] Test case coverage 5 — Smart contract deployment
Verify that the smart contract is successfully deployed
  STR
Send eth_sendRawTransaction with contract code
2.Wait for contract created in logs
Expected The address of the new contract is returned, visible in the logs, gas used is indicated


* [QA] Test case coverage 6 — Executing a contract call
Verify that eth_call is working correctly
STR
1.Execute eth_call to existing contract
2.Get the expected result
Expected  Contract returns value without errors

  
* [QA] Test case coverage 7 — Block generation
Make sure new blocks are created
  STR
1.Send multiple transactions
2.Keep track of Block Number in logs 
Expected New blocks are incremented (1 → 2 → 3...)

  
* [QA] Test case coverage 8 — Checking eth_getTransactionReceipt
Check that you can get a receipt for the transaction
  STR
1.Send transaction
2.Wait eth_getTransactionReceipt
Expected Returns receipt with gasUsed, blockNumber

* [QA] Test case coverage 9 — Checking for invalid request
  Check that you can get a receipt for the transaction
  STR
  1.Send transaction Invalid request check (eth_fakeMethod)
* Expected should return an error.

* [QA] Test case coverage 10 — Test for chainId constraints.
  Verify that the ChainId field accepts only valid numeric values ​​within the allowed range (for example, only positive integers).
  STR
1. Open the "Get Balance" form.
2. In the ChainId field, enter a test value from the list.
3. Click the "Get Balance" button.
4. Check the error message or query result.
  Expected 
If the value is valid (positive integer within supported networks): the request is sent, the balance is returned.

If the value is invalid (negative, not a number, fractional, very large, or empty): the user sees an error message, the request is not sent.


[QA] Test case coverage 11 — Checking for an invalid RPC method (eth_fakeMethod)
  Make sure that when calling a non-existent RPC method (for example, eth_fakeMethod), the system returns the appropriate error, and it is correctly displayed in the UI or logs.
  Prerequisites:
  Web dApp is loaded and connected (Connection State: Connected).
Correct chainId and address are set (it doesn't matter because the method is incorrect).
STR
  1.Open the console (or modify the code/send the request manually if the UI doesn't allow you to enter methods directly).
  2.Initiate a JSON-RPC request with the eth_fakeMethod method:
{
  "jsonrpc": "2.0",
  "method": "eth_fakeMethod",
  "params": [],
  "id": 1
  }
(can be done via fetch or via DevTools → Console)
  Expected 
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
  "code": -32601,
  "message": "Method not found"
  }
  An error message should appear in the UI (if any) — something like "Method not found" or similar.
There should be a record in the logs (#logs) about a failed request with an error.


TODO: Сhecklist for launching a web application (dApp)
Setup Verification
✅ Сhecklist 1: Anvil Starts and Exposes RPC Port
Precondition: Docker image successfully pulled and container started.

Accounts Validation
✅ Сhecklist 2: Verify Availability of Default Accounts

✅ Сhecklist 3: Verify Private Keys are Accessible
Cross-check private keys from the logs with known account addresses.
Expected Result: Each listed address matches a corresponding private key.

Contract Deployment
✅ Сhecklist 4: Contract Deployment Transaction
Transaction Hash: 0x7a1de742012573a97ec78ac0d90d23fa1099529487cc7d8890e3404f3c64fbe6
Contract Address: 0x5FbDB2315678afecb367f032d93F642f64180aa3
Gas Used: 559385
Expected Result: Smart contract is deployed in Block #1 and is callable via eth_call.

Transaction Validation
✅ Сhecklist 5: Validate Second Transaction
Transaction Hash: 0x590a066495914a9e8bf71cceae104321fd62c1d8aecf011d490b0cb2c6660871
Gas Used: 135574
Block: #2
Expected Result: Transaction is mined and receipt is retrievable via eth_getTransactionReceipt.

✅ Сhecklist 6: Validate Third Transaction
Transaction Hash: 0xf46ce388703c80c2a9cf00d865244f60fa4dc90772930fdaddc2d14044cf0c4f
Gas Used: 44550
Block: #3
Expected Result: Transaction successfully included in block; smart contract call returns expected result.

Chain Functionality
✅ Сhecklist 7: Validate Supported JSON-RPC Methods
All methods respond with valid and expected data formats.

Genesis Configuration
✅ Сhecklist 8: Validate Genesis Timestamp and Block
Expected Genesis Block Number: 0
Expected Genesis Timestamp: "Sat, 5 Apr 2025 14:47:49 +0000" (Unix: 1743864469)
 Block metadata should match provided timestamp and parameters.

///////////////
✅ Test cases web application (dApp)

  * [QA] Test case coverage 12 — Checking for connect
  STR
  1.Open the web application (localhost: ****)
  2.Check the connection status in the "Connection state" section.
Expected.
The message "Connection state: Connected" should be displayed.

  * [QA] Test case coverage 13 — Checking the Logs Section
    STR
    1.Open the web application (localhost: ****)
    2.Go to the Logs section.
    3.Expand the Logs tab.
    4.Check for informational and debug messages.
    Expected.
    Messages about client initialization, connection creation, ICE candidate processing, and other messages should be displayed.

* [QA] Test case coverage 14 — Checking whether the "Encrypt Request" checkbox is present.
  STR
  1.Open the web application (localhost: ****)
  2.Check/uncheck the checkbox.
  3.Determine whether the checkbox state is preserved across page reloads.
  Expected.
  Encryption should be enabled by default.
  When the checkbox state changes, the encryption state changes.



* [QA] Test case coverage 16 — Checking the "Get Balance" Section
  STR
  1.Open the web application (localhost: ****)
  2.Go to the "Get Balance" section.
  3.Check for the presence of the "ChainId" and "Address" input fields.
  4.Enter a value in the "ChainId" field (for example, 1).
  Enter a test address in the "Address" field (for example, 0x1234567890abcdef1234567890abcdef12345678).
  5.Click the "Get Balance" button.
 6.Check if the result appears.
  Expected.
  The balance should be displayed according to the entered parameters.


* [QA] Test case coverage 17 — Checking the "Send Funds" button
  STR
  1.Open the web application (localhost: ****)
  2.Go to the "Send Funds" section.
  3.Check for the presence of the "Send Funds" button.
  4.Click on the button.
  5.Check if the result is displayed in the "Result" field.
  Expected.
  A message should be displayed about the result of the operation (success or error).


* [QA] Test case coverage 18 — Verifying the correctness of the "Send Funds" operation result
  STR
  1.Open the web application (localhost: ****)
  2.Use the test data set to verify the transactions.
  3.Enter the correct address and amount to send funds.
  4.Click the "Send Funds" button.
  5.Check the result in the "Result" field.
  Expected.
  A message about a successful transaction or a sending error.


  [QA] Test case coverage 19 — Checking the validity of the address in the "Address" field
  STR
  1.Open the web application (localhost: ****)
  2.Enter an invalid address in the "Address" field (for example, a shorter or incorrectly formatted address).
  3.Click "Get Balance".
  Expected.
  An error message about the invalid address should be displayed.


* [QA] Test case coverage 20 — Verifying the correctness of the result in the "Balance" field
  STR
  1.Open the web application (localhost: ****)
  2.Enter the correct address in the "Address" field.
  3.Click the "Get Balance" button.
  4.Verify the result is displayed.
  Expected.
  The correct balance for the specified address should be displayed.


* [QA] Test case coverage 21 — Check Accordion Functionality (Expand)
  STR
  1.Open the web application (localhost: ****)
  2.Check if the accordion functionality works for the "Logs" and "Methods" sections.
  3.Click on the "Logs", "Get Balance", "Send Funds" headings and check if the corresponding section opens/closes.
  Expected.
  The sections should open and close when clicking on the corresponding headings.




AQA TEST
// The test cases cover:
// - Validate correct handling of parameter order, ensuring the chainId comes before the address.
// - Check that unrecognized method names are properly rejected.
// - Verify that valid parameters result in the expected response.
// - Test handling of incorrect parameter counts, both too few and too many.
// - Ensure proper handling of empty or nil parameter lists.
// - Verify that different valid addresses and chain IDs are accepted.
// - Test the handler's response to various block number formats.
// - Check handling of edge cases like zero addresses and very long addresses.
// - Validate proper error responses for empty chainId or address parameters.
// - Test rejection of non-numeric chainId values.
// - Verify acceptance of addresses without the '0x' prefix.
// - Check handling of very large block numbers.
// - Test rejection of non-numeric characters in the block number.
// - Verify handling of whitespace in parameters.
// - Test case sensitivity in address handling.

=== RUN   TestDefaultApiHandler_GetWalletBalance
=== RUN   TestDefaultApiHandler_GetWalletBalance/Incorrect_params_order
=== RUN   TestDefaultApiHandler_GetWalletBalance/Unrecognized_method_name
=== RUN   TestDefaultApiHandler_GetWalletBalance/Valid_params
=== RUN   TestDefaultApiHandler_GetWalletBalance/Wrong_number_of_params_-_too_few
=== RUN   TestDefaultApiHandler_GetWalletBalance/Wrong_number_of_params_-_too_many
=== RUN   TestDefaultApiHandler_GetWalletBalance/Empty_params
=== RUN   TestDefaultApiHandler_GetWalletBalance/Nil_params
=== RUN   TestDefaultApiHandler_GetWalletBalance/Valid_params_with_different_address
=== RUN   TestDefaultApiHandler_GetWalletBalance/Valid_params_with_different_block
=== RUN   TestDefaultApiHandler_GetWalletBalance/Valid_params_with_numeric_block
=== RUN   TestDefaultApiHandler_GetWalletBalance/Valid_params_with_zero_address
=== RUN   TestDefaultApiHandler_GetWalletBalance/Count_params_less_than_needed
=== RUN   TestDefaultApiHandler_GetWalletBalance/Address_param_is_empty
=== RUN   TestDefaultApiHandler_GetWalletBalance/ChainId_param_is_empty
=== RUN   TestDefaultApiHandler_GetWalletBalance/ChainId_is_not_numeric
=== RUN   TestDefaultApiHandler_GetWalletBalance/Very_long_address
=== RUN   TestDefaultApiHandler_GetWalletBalance/Address_without_0x_prefix
=== RUN   TestDefaultApiHandler_GetWalletBalance/Very_large_block_number
=== RUN   TestDefaultApiHandler_GetWalletBalance/Non-numeric_characters_in_block_number
=== RUN   TestDefaultApiHandler_GetWalletBalance/Whitespace_in_parameters
=== RUN   TestDefaultApiHandler_GetWalletBalance/Case_sensitivity_test
--- PASS: TestDefaultApiHandler_GetWalletBalance (0.00s)
--- PASS: TestDefaultApiHandler_GetWalletBalance/Incorrect_params_order (0.00s)
--- PASS: TestDefaultApiHandler_GetWalletBalance/Unrecognized_method_name (0.00s)
--- PASS: TestDefaultApiHandler_GetWalletBalance/Valid_params (0.00s)
--- PASS: TestDefaultApiHandler_GetWalletBalance/Wrong_number_of_params_-_too_few (0.00s)
--- PASS: TestDefaultApiHandler_GetWalletBalance/Wrong_number_of_params_-_too_many (0.00s)
--- PASS: TestDefaultApiHandler_GetWalletBalance/Empty_params (0.00s)
--- PASS: TestDefaultApiHandler_GetWalletBalance/Nil_params (0.00s)
--- PASS: TestDefaultApiHandler_GetWalletBalance/Valid_params_with_different_address (0.00s)
--- PASS: TestDefaultApiHandler_GetWalletBalance/Valid_params_with_different_block (0.00s)
--- PASS: TestDefaultApiHandler_GetWalletBalance/Valid_params_with_numeric_block (0.00s)
--- PASS: TestDefaultApiHandler_GetWalletBalance/Valid_params_with_zero_address (0.00s)
--- PASS: TestDefaultApiHandler_GetWalletBalance/Count_params_less_than_needed (0.00s)
--- PASS: TestDefaultApiHandler_GetWalletBalance/Address_param_is_empty (0.00s)
--- PASS: TestDefaultApiHandler_GetWalletBalance/ChainId_param_is_empty (0.00s)
--- PASS: TestDefaultApiHandler_GetWalletBalance/ChainId_is_not_numeric (0.00s)
--- PASS: TestDefaultApiHandler_GetWalletBalance/Very_long_address (0.00s)
--- PASS: TestDefaultApiHandler_GetWalletBalance/Address_without_0x_prefix (0.00s)
--- PASS: TestDefaultApiHandler_GetWalletBalance/Very_large_block_number (0.00s)
--- PASS: TestDefaultApiHandler_GetWalletBalance/Non-numeric_characters_in_block_number (0.00s)
--- PASS: TestDefaultApiHandler_GetWalletBalance/Whitespace_in_parameters (0.00s)
--- PASS: TestDefaultApiHandler_GetWalletBalance/Case_sensitivity_test (0.00s)
PASS
ok      command-line-arguments  0.143s




# WEBRTC tester

## Prerequisites

Local Node.js installation with npx and http-server package.

## Running environment
1. Return to root of project directory
2. run docker compose file
```
docker compose up --build
```

## Running the tool
```
npx http-server .
```

## Sample message
```
{
  "id": "1",
  "method": "GetWalletBalance",
  "params": ["1", "0x4838B106FCe9647Bdf1E7877BF73cE8B0BAD5f97"]
}
```

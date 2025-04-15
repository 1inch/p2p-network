FROM node:18-alpine

WORKDIR /app

RUN mkdir -p /app/sdk/examples/basic-web-dapp

COPY sdk /app/sdk

COPY sdk/examples/basic-web-dapp/package*.json /app/sdk/examples/basic-web-dapp/
COPY sdk/examples/basic-web-dapp/ /app/sdk/examples/basic-web-dapp/

WORKDIR /app/sdk/examples/basic-web-dapp

RUN npm install
RUN npm run build
RUN npm install -g serve

EXPOSE 9999

CMD ["serve", "-s", "dist", "-l", "9999"]

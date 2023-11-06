# The 1inch P2P Decentralized Network White Paper

## Abstract

In the emergent era of decentralized technologies, the centralization of web3 applications and backend services presents a paradox that undermines the foundational principles of the blockchain ethos. The 1inch Network is envisaged as a solution to this problem—a peer-to-peer, decentralized network infrastructure that acts as a substrate for the development of diverse web3 products. At its core, the network addresses the quintessential vulnerabilities associated with centralized web services, such as susceptibility to censorship and single points of failure, by providing a robust framework that is resistant to Denial of Service (DoS) attacks through an innovative payment channel system.

Each participant in the network, whether they are a user, a relayer, or a resolver, plays a critical role in sustaining the ecosystem. Users, interacting through lightweight web3 applications, initiate requests within the network. Relayers, operating in a capacity similar to traditional proxies, facilitate the transmission of these requests. Resolvers, the backbone of the network, perform the computations and return the requested data, thereby completing the cycle of interaction. These resolvers are selected based on the 1inch community's trust, conferred through a staking mechanism of the native 1INCH token that engenders a democratized validator-like model.

The amalgamation of these roles is harmonized through a payment system wherein transactions are signed and encapsulated within requests, safeguarding against unwanted network spam and incentivizing honest, efficient responses. This economic model is designed to foster a self-regulating network where quality of service is matched with appropriate compensation, ensuring scalability and sustainability.

The 1inch Network not only revolutionizes the way services are hosted and accessed in the web3 space but also lays the groundwork for an array of decentralized applications—ranging from file storage solutions to streaming services, and from decentralized finance (DeFi) protocols to blockchain query handling—thus heralding a new epoch of internet where decentralization is not just a promise, but an operational reality.

## 1. Introduction

Centralized hosting solutions are at odds with the decentralized ethos of web3. The 1inch Network proposes a P2P network to offer decentralized service alternatives and solve the problems of centralized control and vulnerability to service disruptions.

## 2. Background

Web3's reliance on centralized services undermines its decentralized promise. The 1inch Network's P2P approach aims to resolve this, offering decentralized alternatives to traditional web hosting and backend services.

## 3. The 1inch Network Overview

The network is comprised of three actor types: users, relayers, and resolvers. Users operate through a web3 dApp thin client, relayers function as proxies, and resolvers handle computation and data responses, with each role incentivized via a payment system.

## 4. Technical Architecture

Relayers are registered through an Ethereum smart contract, making their IP address and connection details publicly available for network connections. The RPC approach eliminates the need for DNS, as requests are sent directly to

# Docker compose run Guide

## Step for run docker compose
- make sure that config in dir **assets/** for relayer and resolver are up to date
- call command: ```docker compose up -d```

## Description of docker compose file
Docker compose file consists of 4 containers:
- **blockchain** - represend blockchain network using **anvil** from **fountry-rc**
- **deploy_contract** - container in which the deployment NodeRegistry contract using **forge** from **fountry-rc**
- **relayer**
- **resolver**
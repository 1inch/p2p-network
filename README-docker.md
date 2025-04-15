# Docker compose run Guide

## Step for run docker compose
- make sure that config in dir **assets/** for relayer and resolver are up to date
- call command: ```docker compose up -d```

## Description of docker compose file
Docker compose file consists of 6 containers:
- **blockchain** - represents blockchain network using **anvil** from **foundry-rs**
- **deploy_contract** - container in which the deployment NodeRegistry contract using **forge** from **foundry-rs**
- **relayer** - handles message relaying between nodes
- **resolver** - resolves node addresses and handles registration
- **prometheus** - metrics collection and monitoring
- **grafana** - visualization dashboard for metrics

## Accessing the services
- **Blockchain**: http://localhost:8545
- **Relayer**: http://localhost:8080
- **Resolver**: http://localhost:8001, http://localhost:8081
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000

## Monitoring
The system includes Prometheus for metrics collection and Grafana for visualization:
- Prometheus collects metrics from the relayer and resolver services
- Grafana provides dashboards to visualize the collected metrics
- Default Grafana admin password is set to "admin"

### Grafana Configuration
Grafana is pre-configured with datasources and dashboards through the `assets/provisioning` directory:

1. **Data Sources**: 
   - Prometheus is automatically configured as a data source
   - The connection URL is set to http://prometheus:9090

2. **Dashboards**:
   - Two pre-configured dashboards are automatically loaded:
     - Relayer Dashboard: Shows metrics related to the relayer service
     - Resolver Dashboard: Shows metrics related to the resolver service

3. **Accessing Grafana**:
   - Navigate to http://localhost:3000
   - Log in with the default credentials:
     - Username: admin
     - Password: admin
   - When prompted to change the password, you can either:
     - Set a new password
     - Skip by clicking "Skip" (not recommended for production)

4. **Additional Configuration (Optional)**:
   - Set up alerts:
     - Go to Alerting > Notification channels
     - Add notification channels (email, Slack, etc.)
     - Configure alert rules for important metrics
   - Import additional dashboards:
     - Go to Dashboards > Import
     - Upload JSON files or use dashboard IDs from grafana.com

### Provisioning Directory Structure
The `assets/provisioning` directory contains configuration files for Grafana:

#### Data Sources (`assets/provisioning/datasources/`)
- **datasource.yaml**: Configures the Prometheus data source
  - Sets the name to "Prometheus"
  - Configures the URL to http://prometheus:9090
  - Sets it as the default data source

#### Dashboards (`assets/provisioning/dashboards/`)
- **dashboard.yaml**: Configures the dashboard provider
  - Sets up file-based dashboard provisioning
  - Points to the directory containing dashboard JSON files

- **relayer_dashboard.json**: Dashboard definition for the relayer service
  - Contains panels for monitoring relayer metrics
  - Includes graphs for message throughput, latency, and error rates
  - Provides visualizations for network connectivity and performance

- **resolver_dashboard.json**: Dashboard definition for the resolver service
  - Contains panels for monitoring resolver metrics
  - Includes graphs for registration requests, resolution queries, and response times
  - Provides visualizations for node registry status and performance

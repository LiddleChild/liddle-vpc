# liddle vpc

My personal VPC. The whole VPC is hosted on Oracle Cloud Infrastructure.
Services are exposed through Tailscale Services. Services are hosted with docker and are managed through Portainer.

# Setup

Firstly, apply `docker-compose.yaml` to start up core services.

## Tailscale

Advertise services behind tailscale network. Hosted in container [https://tailscale.com/docs/features/containers/docker/how-to/connect-docker-container]

To serve new service

```bash
# Start serving service
tailscale serve --service=svc:<service-name> <url>

# Stop serving
tailscale serve --service=svc:<service-name> off
```

## Portainer

Manage containers in web ui

# Services
- Grist


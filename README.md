# liddle vpc

My personal VPC. The whole VPC is hosted on Oracle Cloud Infrastructure.
Services are exposed through Tailscale Services. Services are hosted with docker and are managed through Portainer.

## Tailscale

Installation guide from tailscale [https://pkgs.tailscale.com/stable/#oracle-9]

```bash
sudo tailscale up --auth-key=<tskey-client-key> --advertise-tags=tag:<tag>
```

Serving new service

```bash
# Start serving service
tailscale serve --service=svc:<service-name> --https=443 <port>

# Stop serving
tailscale serve --service=svc:<service-name> off
```

## Services

- Portainer
- Grist


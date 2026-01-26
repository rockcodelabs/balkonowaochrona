# Balkonowa Ochrona

Professional website for balcony protection services - pigeon net installation and balcony cleaning services.

**Live:** https://balkonowaochrona.pl

## Project Structure

```
balkonowaochrona/
├── main.go             # Go server with Resend email integration
├── go.mod              # Go module definition
├── index.html          # Main HTML page
├── style.css           # Stylesheet
├── favicon.svg         # Site favicon
├── Dockerfile          # Optimized Docker configuration (7.7MB image)
├── config/
│   └── deploy.yml      # Kamal deployment configuration
├── .kamal/
│   └── secrets         # Bitwarden secrets integration
├── .github/
│   └── workflows/
│       └── deploy.yml  # GitHub Actions CI/CD
├── start.sh            # Local start script with Bitwarden
└── README.md           # This file
```

## Features

- Clean, modern responsive design
- Contact form with email via Resend API
- Mobile-friendly layout
- Ultra-minimal Docker image (**7.7MB** using Go + scratch)
- Deployed via Kamal to Raspberry Pi 5
- Cloudflare Tunnel for secure access

## DNS Configuration

**Cloudflare Nameservers (for JDM.pl):**
- `hope.ns.cloudflare.com`
- `keaton.ns.cloudflare.com`

**DNS Records:**
| Type  | Name | Target |
|-------|------|--------|
| CNAME | @    | `c43349ad-0aba-4694-b89f-ac4e5acebcfe.cfargotunnel.com` |
| CNAME | www  | `c43349ad-0aba-4694-b89f-ac4e5acebcfe.cfargotunnel.com` |

## Prerequisites

- Docker
- Kamal (`gem install kamal`)
- Bitwarden CLI (`brew install bitwarden-cli`)
- Ruby 3.2.2 (for Kamal)

## Local Development

### Run with Docker:

```bash
# Unlock Bitwarden
export BW_SESSION=$(bw unlock --raw)

# Start with Bitwarden secrets
./start.sh --build
```

Then open http://localhost:4001

### Run with Go directly:

```bash
go run main.go
```

## Deployment

### Automatic (GitHub Actions)

Push to `main` branch → automatically deploys to Raspberry Pi 5

### Manual

```bash
# Unlock Bitwarden
export BW_SESSION=$(bw unlock --raw)
bw sync

# Deploy
kamal deploy
```

## Secrets (Bitwarden)

Required items in Bitwarden Password Manager:

| Item Name | Field | Description |
|-----------|-------|-------------|
| `registry-credentials` | username | Docker Hub username |
| `registry-credentials` | password | Docker Hub password/token |
| `kalkowski-resend` | notes | Resend API key |

## GitHub Actions Secrets

| Secret | Description |
|--------|-------------|
| `KAMAL_REGISTRY_USERNAME` | Docker Hub username |
| `KAMAL_REGISTRY_PASSWORD` | Docker Hub password |
| `RESEND_API_KEY` | Resend API key |
| `SSH_PRIVATE_KEY` | SSH key for Pi access |
| `PI_HOST` | Pi hostname/IP |

## Infrastructure

- **Server:** Raspberry Pi 5 (`pi5main.local`)
- **Proxy:** Kamal Proxy (port 80)
- **Tunnel:** Cloudflare Tunnel (`kw-staging`)
- **Image Size:** 7.7MB (Go binary + static files)

## Cloudflare Tunnel Config

Location on Pi: `/etc/cloudflared/config.yml`

```yaml
tunnel: c43349ad-0aba-4694-b89f-ac4e5acebcfe
credentials-file: /root/.cloudflared/c43349ad-0aba-4694-b89f-ac4e5acebcfe.json

ingress:
  - hostname: panel.taterniczek.pl
    service: http://localhost:80
  - hostname: balkonowaochrona.pl
    service: http://localhost:80
  - hostname: www.balkonowaochrona.pl
    service: http://localhost:80
  - service: http_status:404
```

Restart tunnel: `sudo systemctl restart cloudflared`

## Useful Commands

```bash
# SSH to Pi
ssh rege@pi5main.local

# Check containers
docker ps | grep balkonowa

# View logs
kamal app logs

# Restart app
kamal app boot

# Check tunnel status
sudo systemctl status cloudflared

# Kamal proxy list
docker exec kamal-proxy kamal-proxy list
```

## Contact

- **Phone:** 503 508 987
- **Email:** kalkowski123@gmail.com

## License

© 2024 Balkonowa Ochrona. All rights reserved.
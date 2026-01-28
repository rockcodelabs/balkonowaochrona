# Balkonowa Ochrona

Professional website for balcony protection services - pigeon net installation and balcony cleaning services.

**Live:** https://balkonowaochrona.pl | https://siatkanakota.pl

## Architecture

```
Internet → Cloudflare → Cloudflare Tunnel → Pi5 (localhost:80)
                                               ↓
                                          kamal-proxy
                                               ↓
                                    Docker container (port 4001)
```

## Project Structure

```
balkonowaochrona/
├── main.go              # Go server with Resend email integration
├── index.html           # Main HTML page with gallery
├── style.css            # Stylesheet
├── photos/              # Gallery images
├── favicon.svg          # Site favicon
├── Dockerfile           # Multi-stage Docker build (~8MB image)
├── go.mod               # Go module definition
├── config/
│   └── deploy.yml       # Kamal config (reference only)
├── .github/
│   └── workflows/
│       └── deploy.yml   # GitHub Actions CI/CD
├── CLAUDE.md            # AI assistant guide
└── README.md            # This file
```

## Features

- Clean, modern responsive design
- Photo gallery with lightbox zoom
- Contact form with email via Resend API
- Mobile-friendly layout
- Ultra-minimal Docker image (~8MB using Go + scratch)
- Automated deployment via GitHub Actions
- Cloudflare Tunnel for secure access (no exposed ports)

## Infrastructure

| Component | Description |
|-----------|-------------|
| **Server** | Raspberry Pi 5 (`pi5main.local`) |
| **Proxy** | kamal-proxy (routes by hostname) |
| **Tunnel** | Cloudflare Tunnel (secure, no open ports) |
| **Runner** | Self-hosted GitHub Actions (Docker on Pi) |
| **Registry** | Docker Hub (`regedarek/balkonowaochrona`) |

## Local Development

### Quick static preview:

```bash
python3 -m http.server 8080
# Open http://localhost:8080
```

### Full Docker test:

```bash
docker build -t balkonowaochrona-test .
docker run --rm -p 4001:4001 -e RESEND_API_KEY="your_key" balkonowaochrona-test
# Open http://localhost:4001
```

## Deployment

### Automatic (recommended)

Push to `main` branch → GitHub Actions automatically deploys to Pi

### What happens on deploy:

1. Runner builds Docker image
2. Pushes to Docker Hub
3. Stops old container, starts new one on `kamal` network
4. Registers with kamal-proxy for domain routing

## Secrets

Managed via Bitwarden Secrets Manager, fetched at deploy time:

| Secret | Description |
|--------|-------------|
| `KAMAL_REGISTRY_USERNAME` | Docker Hub username |
| `KAMAL_REGISTRY_PASSWORD` | Docker Hub password |
| `RESEND_API_KEY` | Resend API key for emails |

## Useful Commands

```bash
# SSH to Pi
ssh rege@pi5main.local

# Check container status
docker ps | grep balkonowa

# View logs
docker logs balkonowaochrona

# Restart container
docker restart balkonowaochrona

# Check proxy routing
docker exec kamal-proxy kamal-proxy list

# Check tunnel status
sudo systemctl status cloudflared
```

## DNS Configuration

**Cloudflare Nameservers (configured at registrar):**
- `hope.ns.cloudflare.com`
- `keaton.ns.cloudflare.com`

**DNS Records (in Cloudflare):**

| Type  | Name | Target |
|-------|------|--------|
| CNAME | @    | `c43349ad-0aba-4694-b89f-ac4e5acebcfe.cfargotunnel.com` |
| CNAME | www  | `c43349ad-0aba-4694-b89f-ac4e5acebcfe.cfargotunnel.com` |

## Contact

- **Phone:** 503 508 987
- **Email:** kalkowski123@gmail.com

## License

© 2024 Balkonowa Ochrona. All rights reserved.
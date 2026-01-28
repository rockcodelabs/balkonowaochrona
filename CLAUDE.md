# CLAUDE.md - Project Guide for AI Assistants

## Project Overview

**Balkonowa Ochrona** is a simple static website for balcony protection services (pigeon net installation and cleaning). It's a Go server serving static HTML with a contact form that sends emails via Resend API.

## Architecture

```
Internet → Cloudflare → Cloudflare Tunnel → Pi5 (localhost:80)
                                               ↓
                                          kamal-proxy
                                               ↓
                                    Docker container (port 4001)
```

### Components

| Component | Description |
|-----------|-------------|
| **Cloudflare** | DNS, CDN, DDoS protection |
| **Cloudflare Tunnel** | Secure tunnel from Cloudflare to Pi (no exposed ports) |
| **kamal-proxy** | Reverse proxy on Pi, routes by hostname to containers |
| **Go server** | Minimal server serving static files + `/api/contact` endpoint |
| **Resend** | Email API for contact form |

### Infrastructure

- **Server:** Raspberry Pi 5 (`pi5main.local`)
- **Runner:** Self-hosted GitHub Actions runner (Docker container on Pi)
- **Network:** Containers use `kamal` Docker network
- **Domains:** balkonowaochrona.pl, siatkanakota.pl (+ www variants)

## Deployment

### How it works

1. Push to `main` branch triggers GitHub Actions
2. Runner (Docker on Pi) builds and pushes image to Docker Hub
3. Runner stops old container, starts new one on `kamal` network
4. Runner registers container with kamal-proxy

### Key files

| File | Purpose |
|------|---------|
| `.github/workflows/deploy.yml` | CI/CD pipeline |
| `Dockerfile` | Multi-stage build (Go → scratch, ~8MB image) |
| `config/deploy.yml` | Kamal config (kept for reference, not actively used) |
| `main.go` | Go server with static files + email endpoint |

### Secrets (Bitwarden Secrets Manager)

Fetched at deploy time via `bws` CLI:
- `KAMAL_REGISTRY_USERNAME` - Docker Hub username
- `KAMAL_REGISTRY_PASSWORD` - Docker Hub password
- `RESEND_API_KEY` - Resend API key for emails

### Manual deployment commands

```bash
# SSH to Pi
ssh rege@pi5main.local

# Check running containers
docker ps | grep balkonowa

# View logs
docker logs balkonowaochrona

# Restart container
docker restart balkonowaochrona

# Check kamal-proxy routing
docker exec kamal-proxy kamal-proxy list

# Update kamal-proxy to point to new container
docker exec kamal-proxy kamal-proxy deploy balkonowaochrona-web \
  --target <container_id>:4001 \
  --host balkonowaochrona.pl \
  --host www.balkonowaochrona.pl \
  --host siatkanakota.pl \
  --host www.siatkanakota.pl

# Check tunnel status
sudo systemctl status cloudflared
```

## Project Structure

```
balkonowaochrona/
├── main.go              # Go server (static files + /api/contact)
├── index.html           # Main page with gallery and contact form
├── style.css            # Styles
├── photos/              # Gallery images
├── favicon.svg          # Favicon
├── Dockerfile           # Multi-stage build
├── go.mod               # Go module
├── config/
│   └── deploy.yml       # Kamal config (reference only)
└── .github/
    └── workflows/
        └── deploy.yml   # GitHub Actions CI/CD
```

## Development

### Local testing

```bash
# With Python (static only, no API)
python3 -m http.server 8080

# With Docker (full functionality)
docker build -t balkonowaochrona-test .
docker run --rm -p 4001:4001 -e RESEND_API_KEY="your_key" balkonowaochrona-test
```

### Contact form

- **Endpoint:** `POST /api/contact`
- **Payload:** `{ "name", "email", "phone", "message" }`
- **Email sent to:** Configured via `TO_EMAIL` env var
- **Email sent from:** `onboarding@resend.dev` (Resend default)

## Important Notes

1. **No Ruby/Kamal CLI needed** - Deployment uses pure Docker
2. **kamal-proxy required** - Container must be on `kamal` network and registered with proxy
3. **Cloudflare Tunnel** - All traffic goes through tunnel, no ports exposed on Pi
4. **Self-hosted runner** - Runs as Docker container on Pi, has access to host Docker socket
# Balkonowa Ochrona - Static Website

Professional website for balcony protection services - pigeon net installation and balcony cleaning services.

## Project Structure

```
balkonowaochrona/
├── index.html          # Main HTML page
├── style.css           # Stylesheet
├── favicon.svg         # Site favicon
├── server.js           # Node.js server with Resend integration
├── package.json        # Node.js dependencies
├── package-lock.json   # Dependency lock file
├── Dockerfile          # Optimized Docker configuration
├── docker-compose.yml  # Docker Compose configuration
├── start.sh            # Start script with Bitwarden integration
└── README.md           # This file
```

## Features

- Clean, modern responsive design
- Contact form with email, name, and message fields
- Email sending via Resend API
- Mobile-friendly layout
- Professional styling with gradient backgrounds
- Optimized Docker image (~189MB using distroless)

## Prerequisites

- Docker
- Bitwarden CLI (`bw`) for secrets management
- Resend API key

## Setup

### 1. Get Resend API Key

1. Create an account at [resend.com](https://resend.com)
2. Go to API Keys and create a new API key
3. (Optional) Add and verify your domain for custom "from" addresses

### 2. Store API Key in Bitwarden

1. Install Bitwarden CLI:
   ```bash
   brew install bitwarden-cli
   ```

2. Log in to Bitwarden:
   ```bash
   bw login
   ```

3. Create a new login item in Bitwarden:
   - **Name:** `resend-api-key`
   - **Password:** Your Resend API key (e.g., `re_xxxxxxxx`)

   Or via CLI:
   ```bash
   bw get template item | jq '.name="resend-api-key" | .login.password="re_your_api_key_here"' | bw encode | bw create item
   ```

### 3. Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `RESEND_API_KEY` | Yes | Your Resend API key (fetched from Bitwarden) |
| `TO_EMAIL` | No | Recipient email (default: kalkowski123@gmail.com) |
| `FROM_EMAIL` | No | Sender email (default: onboarding@resend.dev) |
| `PORT` | No | Server port (default: 4001) |

## Running with Bitwarden (Recommended)

Use the provided start script that automatically fetches secrets from Bitwarden:

```bash
# First run (builds the image)
./start.sh --build

# Subsequent runs
./start.sh
```

The script will:
1. Check Bitwarden CLI is installed
2. Log in / unlock your vault if needed
3. Fetch the Resend API key from Bitwarden
4. Start the Docker container with the secret

## Running with Docker Manually

### Build the optimized image:

```bash
docker build -t balkonowa-ochrona .
```

### Run with environment variables:

```bash
docker run -d \
  --name balkonowa \
  -p 4001:4001 \
  -e RESEND_API_KEY=your_api_key_here \
  -e TO_EMAIL=kalkowski123@gmail.com \
  balkonowa-ochrona
```

Then open http://localhost:4001 in your browser.

### Stop and remove container:

```bash
docker stop balkonowa && docker rm balkonowa
```

## Running with Docker Compose

### With Bitwarden (recommended):

```bash
export RESEND_API_KEY=$(bw get password resend-api-key)
docker-compose up -d
```

### View logs:

```bash
docker-compose logs -f
```

### Stop:

```bash
docker-compose down
```

## Docker Image Optimization

The image has been optimized for minimal size:

| Tag | Size | Description |
|-----|------|-------------|
| `distroless` | **189MB** | Production-ready, minimal attack surface |
| `slim` | 193MB | Alpine-based with esbuild bundle |
| `optimized` | 216MB | Alpine with production deps only |
| `latest` | 266MB | Original unoptimized |

Optimizations applied:
- **Multi-stage build** - Build dependencies don't end up in final image
- **esbuild bundling** - Single JS file, no node_modules needed
- **Distroless base** - Minimal OS, only Node.js runtime
- **Non-root user** - Security best practice
- **.dockerignore** - Excludes unnecessary files from build context

## Running Locally (without Docker)

### 1. Install dependencies:

```bash
npm install
```

### 2. Set environment variables and run:

```bash
export RESEND_API_KEY=$(bw get password resend-api-key)
npm start
```

## Deployment Options

### Docker on VPS

Use the included Dockerfile to deploy on any VPS with Docker support:

```bash
# On your VPS
docker pull your-registry/balkonowa-ochrona
docker run -d -p 80:4001 -e RESEND_API_KEY=xxx balkonowa-ochrona
```

### Railway / Render / Fly.io

1. Connect your GitHub repository
2. Set environment variables in the dashboard
3. Deploy

### Kubernetes

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: balkonowa-secrets
type: Opaque
stringData:
  RESEND_API_KEY: "your-api-key"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: balkonowa
spec:
  replicas: 1
  selector:
    matchLabels:
      app: balkonowa
  template:
    metadata:
      labels:
        app: balkonowa
    spec:
      containers:
      - name: web
        image: balkonowa-ochrona:distroless
        ports:
        - containerPort: 4001
        envFrom:
        - secretRef:
            name: balkonowa-secrets
```

## Form Configuration

The contact form sends emails via Resend API with the following fields:
- **Name** (required)
- **Email** (required) - set as reply-to address
- **Phone** (optional)
- **Message** (required)

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)
- Mobile browsers (iOS Safari, Chrome for Android)

## Troubleshooting

### Bitwarden CLI issues

```bash
# Check status
bw status

# Force sync
bw sync

# Unlock vault
export BW_SESSION=$(bw unlock --raw)
```

### Container issues

```bash
# View logs
docker logs balkonowa

# Check if running
docker ps | grep balkonowa

# Restart
docker restart balkonowa
```

### Email not sending

1. Check RESEND_API_KEY is set correctly
2. Verify API key is valid at resend.com
3. Check container logs for error messages

## License

© 2024 Balkonowa Ochrona. All rights reserved.
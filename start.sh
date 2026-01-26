#!/bin/bash
set -e

# Balkonowa Ochrona - Start script with Bitwarden secrets
# Usage: ./start.sh [--build]

CONTAINER_NAME="balkonowa"
IMAGE_NAME="balkonowa-ochrona"
PORT=4001

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🏠 Balkonowa Ochrona - Starting...${NC}"

# Check if Bitwarden CLI is installed
if ! command -v bw &> /dev/null; then
    echo -e "${RED}❌ Bitwarden CLI (bw) is not installed.${NC}"
    echo "Install it with: brew install bitwarden-cli"
    echo "Or visit: https://bitwarden.com/help/cli/"
    exit 1
fi

# Check Bitwarden login status
BW_STATUS=$(bw status | grep -o '"status":"[^"]*"' | cut -d'"' -f4)

if [ "$BW_STATUS" = "unauthenticated" ]; then
    echo -e "${YELLOW}🔐 Please log in to Bitwarden:${NC}"
    bw login
    BW_STATUS="locked"
fi

if [ "$BW_STATUS" = "locked" ]; then
    echo -e "${YELLOW}🔓 Unlocking Bitwarden vault...${NC}"
    export BW_SESSION=$(bw unlock --raw)
    if [ -z "$BW_SESSION" ]; then
        echo -e "${RED}❌ Failed to unlock Bitwarden vault${NC}"
        exit 1
    fi
fi

# Sync vault to get latest secrets
echo -e "${GREEN}🔄 Syncing Bitwarden vault...${NC}"
bw sync

# Fetch Resend API key from Bitwarden
# The item should be named "resend-api-key" or you can use the item ID
echo -e "${GREEN}🔑 Fetching Resend API key from Bitwarden...${NC}"

RESEND_API_KEY=$(bw get password "resend-api-key" 2>/dev/null) || {
    echo -e "${RED}❌ Could not find 'resend-api-key' in Bitwarden.${NC}"
    echo ""
    echo "Please create a login item in Bitwarden with:"
    echo "  - Name: resend-api-key"
    echo "  - Password: your Resend API key"
    echo ""
    echo "Or use a custom field. You can also use the item ID instead of the name."
    exit 1
}

if [ -z "$RESEND_API_KEY" ]; then
    echo -e "${RED}❌ Resend API key is empty${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Successfully retrieved Resend API key${NC}"

# Stop existing container if running
echo -e "${GREEN}🛑 Stopping existing container...${NC}"
docker rm -f $CONTAINER_NAME 2>/dev/null || true

# Build image if requested or if it doesn't exist
if [ "$1" = "--build" ] || [ -z "$(docker images -q $IMAGE_NAME 2>/dev/null)" ]; then
    echo -e "${GREEN}🔨 Building Docker image...${NC}"
    docker build -t $IMAGE_NAME .
fi

# Run container with secrets
echo -e "${GREEN}🚀 Starting container...${NC}"
docker run -d \
    --name $CONTAINER_NAME \
    --restart unless-stopped \
    -p $PORT:$PORT \
    -e RESEND_API_KEY="$RESEND_API_KEY" \
    -e TO_EMAIL="kalkowski123@gmail.com" \
    -e FROM_EMAIL="onboarding@resend.dev" \
    $IMAGE_NAME

# Wait for container to be healthy
echo -e "${GREEN}⏳ Waiting for container to be ready...${NC}"
sleep 2

# Check if container is running
if docker ps | grep -q $CONTAINER_NAME; then
    echo -e "${GREEN}✅ Container is running!${NC}"
    echo -e "${GREEN}🌐 Website available at: http://localhost:$PORT${NC}"
    echo ""
    docker logs $CONTAINER_NAME
else
    echo -e "${RED}❌ Container failed to start${NC}"
    docker logs $CONTAINER_NAME
    exit 1
fi
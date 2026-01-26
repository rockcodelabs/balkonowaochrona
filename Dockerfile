# Stage 1: Bundle application with esbuild
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install all dependencies (including dev for esbuild)
RUN npm ci

# Copy source files
COPY server.js ./

# Bundle into single file with all dependencies
RUN npx esbuild server.js --bundle --platform=node --target=node20 --outfile=server.bundle.js --minify

# Stage 2: Distroless production image (ultra-minimal)
FROM gcr.io/distroless/nodejs20-debian12:nonroot

WORKDIR /app

# Copy only the bundled server (no node_modules needed!)
COPY --from=builder /app/server.bundle.js ./server.js

# Copy static files
COPY index.html ./
COPY style.css ./
COPY favicon.svg ./

# Expose port
EXPOSE 4001

# Start the server
CMD ["server.js"]
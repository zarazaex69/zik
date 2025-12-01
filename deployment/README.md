# ZIK Deployment Guide

Production deployment configuration for ZIK using Docker Compose and Cloudflare Tunnel.

## Prerequisites

- Docker and Docker Compose installed
- Cloudflare account with tunnel configured
- Domain: `zarazaex.xyz`

## Directory Structure

```
~/srv/
├── containers/
│   └── compose.yaml          # Main compose with all services
├── zik/                      # This repository
│   ├── apps/
│   │   ├── server/
│   │   │   └── Dockerfile
│   │   └── web/
│   │       └── Dockerfile
│   └── deployment/
│       ├── compose.yaml      # ZIK-specific compose
│       ├── .env.example
│       └── README.md
```

## Setup

### Option 1: Add to Existing Compose

Add ZIK services to your existing `~/srv/containers/compose.yaml`:

```yaml
services:
  # ... existing services (zarazaex, zcell, cloudflared)

  zik-server:
    container_name: zik-server
    build:
      context: ../zik/apps/server
      dockerfile: Dockerfile
    restart: unless-stopped
    ports:
      - "127.0.0.1:8804:8804"
    environment:
      - PORT=8804
      - DEBUG=false
    networks:
      - server-net
    stop_grace_period: 5s

  zik-web:
    container_name: zik-web
    build:
      context: ../zik/apps/web
      dockerfile: Dockerfile
    restart: unless-stopped
    ports:
      - "127.0.0.1:8805:8805"
    environment:
      - PORT=8805
      - API_URL=https://zik-api.zarazaex.xyz
    networks:
      - server-net
    depends_on:
      - zik-server
    stop_grace_period: 5s
```

### Option 2: Separate Compose File

Use the provided `deployment/compose.yaml` separately:

```bash
cd ~/srv/zik/deployment
cp .env.example .env
docker compose up -d
```

## Cloudflare Configuration

Your existing Cloudflare Tunnel already handles routing. Add these routes in Cloudflare dashboard:

1. **API Server**:
   - Hostname: `zik-api.zarazaex.xyz`
   - Service: `http://localhost:8804`
   - Path: `*`

2. **Web UI**:
   - Hostname: `zik.zarazaex.xyz`
   - Service: `http://localhost:8805`
   - Path: `*`

## Deployment Commands

### Build and Start

```bash
# From ~/srv/containers/
docker compose up -d zik-server zik-web

# Or rebuild
docker compose up -d --build zik-server zik-web
```

### View Logs

```bash
docker compose logs -f zik-server
docker compose logs -f zik-web
```

### Stop Services

```bash
docker compose stop zik-server zik-web
```

### Update Deployment

```bash
cd ~/srv/zik
git pull
cd ~/srv/containers
docker compose up -d --build zik-server zik-web
```

## Health Checks

Both services include health checks:

```bash
# Check server health
curl http://localhost:8804/health

# Check web UI
curl http://localhost:8805/

# Docker health status
docker compose ps
```

## Environment Variables

Create `.env` file in `~/srv/containers/` or `~/srv/zik/deployment/`:

```bash
# Optional Z.AI token
ZAI_TOKEN=your_token_here

# Server config
PORT=8804
DEBUG=false
MODEL=GLM-4-6-API-V1
```

## Monitoring

### Container Status

```bash
docker compose ps
```

### Resource Usage

```bash
docker stats zik-server zik-web
```

### Logs

```bash
# Real-time logs
docker compose logs -f zik-server zik-web

# Last 100 lines
docker compose logs --tail=100 zik-server
```

## Troubleshooting

### Server Not Starting

```bash
# Check logs
docker compose logs zik-server

# Check if port is available
netstat -tuln | grep 8804

# Restart service
docker compose restart zik-server
```

### Web UI Can't Connect to API

1. Verify API is running: `curl http://localhost:8804/health`
2. Check environment variable: `API_URL=https://zik-api.zarazaex.xyz`
3. Verify Cloudflare route for `zik-api.zarazaex.xyz`

### Cloudflare Tunnel Issues

```bash
# Check cloudflared logs
docker compose logs cloudflared

# Verify routes in Cloudflare dashboard
# Ensure zik-api.zarazaex.xyz -> localhost:8804
# Ensure zik.zarazaex.xyz -> localhost:8805
```

## Security Notes

1. Services bind to `127.0.0.1` only (not exposed directly)
2. All external access goes through Cloudflare Tunnel
3. Rate limiting: 30 requests/minute per IP
4. No authentication required (public API)

## Backup

No persistent data to backup. Configuration is in git repository.

## Updates

```bash
# Pull latest changes
cd ~/srv/zik
git pull

# Rebuild and restart
cd ~/srv/containers
docker compose up -d --build zik-server zik-web
```

## Production Checklist

- [ ] Clone repository to `~/srv/zik`
- [ ] Add services to `~/srv/containers/compose.yaml`
- [ ] Configure Cloudflare routes
- [ ] Build and start containers
- [ ] Verify health checks
- [ ] Test API: `curl https://zik-api.zarazaex.xyz/health`
- [ ] Test Web UI: `curl https://zik.zarazaex.xyz/`
- [ ] Monitor logs for errors

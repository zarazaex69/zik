# Build stage
FROM docker.io/oven/bun:1-alpine AS builder

WORKDIR /build

# Copy package files
COPY package.json bun.lock* ./
RUN bun install --frozen-lockfile

# Copy source code
COPY . .

# Runtime stage
FROM docker.io/oven/bun:1-alpine

WORKDIR /app

# Copy everything from builder (source + dependencies)
COPY --from=builder /build ./

# Expose port
EXPOSE 8805

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8805/ || exit 1

# Run as non-root user
USER bun

# Set production environment
ENV NODE_ENV=production

CMD ["bun", "src/index.ts"]

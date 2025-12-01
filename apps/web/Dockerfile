# Build stage
FROM oven/bun:1-alpine AS builder

WORKDIR /build

# Copy package files
COPY package.json bun.lock* ./
RUN bun install --frozen-lockfile

# Copy source code
COPY . .

# Build application
RUN bun run build

# Runtime stage
FROM oven/bun:1-alpine

WORKDIR /app

# Copy built application
COPY --from=builder /build/dist ./dist
COPY --from=builder /build/src ./src
COPY --from=builder /build/package.json ./

# Install production dependencies only
RUN bun install --production --frozen-lockfile

# Expose port
EXPOSE 8805

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8805/ || exit 1

# Run as non-root user
USER bun

CMD ["bun", "start"]

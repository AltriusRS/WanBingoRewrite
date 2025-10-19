# Multi-stage build for production
FROM oven/bun:1.3 AS builder

WORKDIR /app

# Copy package files
COPY package.json bun.lock ./

# Install dependencies
RUN bun install --frozen-lockfile

# Copy source code
COPY . .
RUN rm -rf .env .env.example

# Set build args
ARG NEXT_PUBLIC_API_ROOT
ENV NEXT_PUBLIC_API_ROOT=$NEXT_PUBLIC_API_ROOT

ARG NEXT_PUBLIC_POSTHOG_KEY
ENV NEXT_PUBLIC_POSTHOG_KEY=$NEXT_PUBLIC_POSTHOG_KEY

ARG NEXT_PUBLIC_POSTHOG_HOST
ENV NEXT_PUBLIC_POSTHOG_HOST=$NEXT_PUBLIC_POSTHOG_HOST

# Build the application
RUN bun run build

# Production stage
FROM oven/bun:1.3 AS runner

WORKDIR /app

# Copy built application
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/package.json ./
COPY --from=builder /app/next.config.mjs ./next.config.mjs

# Install only production dependencies
RUN bun install --omit=dev

# Expose port
EXPOSE 3000

# Start the application
CMD ["bun", "run", "start"]
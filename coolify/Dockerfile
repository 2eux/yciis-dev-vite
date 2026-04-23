# Edusys Pro - Coolify Deployment
# =============================================
# This Dockerfile is optimized for deployment on Coolify
# 
# Environment Variables (set in Coolify):
# - SUPABASE_URL: Your Supabase project URL
# - SUPABASE_ANON_KEY: Your Supabase anon key  
# - SUPABASE_SERVICE_KEY: Your Supabase service role key
# - JWT_SECRET: A strong random string for JWT signing
# =============================================

FROM oven/bun:1-alpine AS base
WORKDIR /app

# =============================================
# Dependencies
# =============================================
FROM base AS deps
RUN bun install

# =============================================
# Builder
# =============================================
FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

# =============================================
# Production
# =============================================
FROM base AS runner
WORKDIR /app

# Set production environment
ENV NODE_ENV=production

# Copy built files
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000/api/health || exit 1

# Start command
CMD ["bun", "x", "hyper", "dist/index.js"]
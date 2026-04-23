FROM node:18-alpine

WORKDIR /app

# Copy package management files first for caching
COPY package.json package-lock.json* ./

# Install all dependencies (including devDependencies required for build)
RUN NODE_ENV=development npm install

# Copy all application files
COPY . .

# Build the production application bundle
RUN npm run build

# Expose the internal port specified in docker network
EXPOSE 3000

# Start Vite's preview server bound to all network interfaces.
# This serves the built /dist folder and supports API proxies.
CMD ["npm", "run", "preview", "--", "--host", "0.0.0.0", "--port", "3000"]

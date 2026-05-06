FROM node:18-alpine AS builder
WORKDIR /app

ARG VITE_CONVEX_URL
ENV VITE_CONVEX_URL=${VITE_CONVEX_URL}

COPY web/package.json ./
RUN npm install

COPY web/ ./
RUN npm run build

FROM nginx:alpine

RUN apk add --no-cache curl

RUN rm /etc/nginx/conf.d/default.conf
COPY web/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=5s --retries=3 \
    CMD curl -f http://localhost:80/ || exit 1

CMD ["nginx", "-g", "daemon off;"]
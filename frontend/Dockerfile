FROM node:20-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM node:20-alpine AS runner

WORKDIR /app
ENV NODE_ENV=production

RUN npm install --global serve@14.2.4

COPY --from=builder /app/dist ./dist

EXPOSE 4173

CMD ["sh", "-c", "printf 'window.__APP_CONFIG__ = { API_BASE_URL: \"%s\" };\\n' \"${VITE_API_BASE_URL:-}\" > /app/dist/runtime-config.js && serve -s dist -l ${PORT:-4173}"]

FROM node:20-alpine3.22 AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .

ARG VUE_APP_API_URL
ENV VUE_APP_API_URL=${VUE_APP_API_URL}

RUN npm run build

# --- Финальный этап (Final) ---
FROM nginx:stable-alpine3.21

RUN rm /etc/nginx/conf.d/default.conf

COPY nginx.conf /etc/nginx/templates/default.conf.template

COPY --from=builder /app/dist/ /usr/share/nginx/html/

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]

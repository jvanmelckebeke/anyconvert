FROM node:21-alpine3.17 as build

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY public ./public
COPY src ./src
COPY index.html ./
COPY vite.config.js ./


RUN npm run build

FROM nginx:1.21-alpine as runtime
EXPOSE 80

ENV API_HOST=api
ENV API_PORT=8000

COPY nginx/resolver.conf /etc/nginx/conf.d/resolver.conf
COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf
COPY nginx/variables.conf.template /etc/nginx/templates/10-variables.conf.template

COPY --from=build /app/dist /usr/share/nginx/html
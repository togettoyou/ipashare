FROM node:lts-alpine as build-stage
WORKDIR /app
COPY package*.json ./
RUN yarn install --registry=https://registry.npmmirror.com
COPY . .
RUN yarn run build:prod

FROM nginx:stable-alpine as production-stage
COPY --from=build-stage /app/dist /usr/share/nginx/html
RUN rm /etc/nginx/conf.d/default.conf
COPY nginx.template /etc/nginx/conf.d
EXPOSE 80
WORKDIR /etc/nginx/conf.d
ENTRYPOINT envsubst '${BACKEND_URL}' < nginx.template > nginx.conf && cat nginx.conf && nginx -g 'daemon off;'

# Build stage
FROM node:lts-alpine AS build-stage

WORKDIR /app

COPY package*.json ./
RUN npm ci 

COPY . .
RUN npm run build

# Production stage
FROM nginx:alpine

RUN rm -rf /usr/share/nginx/html*

COPY --from=build-stage /app/dist /usr/share/nginx/html

COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
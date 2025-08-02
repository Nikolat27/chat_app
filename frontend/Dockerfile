# Build stage
FROM node:22-alpine AS build-stage

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci --only=production

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage
FROM nginx:alpine AS production-stage

# Copy built application from build stage
COPY --from=build-stage /app/dist /usr/share/nginx/html

# Copy custom nginx configuration template
COPY nginx.conf /etc/nginx/conf.d/default.conf.template

# Expose port 80
EXPOSE 80

# Create startup script
RUN echo '#!/bin/sh\n\
envsubst < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf\n\
nginx -g "daemon off;"' > /start.sh && chmod +x /start.sh

# Start nginx with environment variable substitution
CMD ["/start.sh"] 
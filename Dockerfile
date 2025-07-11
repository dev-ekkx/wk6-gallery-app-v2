# # Step 1: Build Angular frontend
# FROM node:20-alpine AS frontend
# WORKDIR /app/frontend
# COPY frontend/ .
# RUN npm install && npm run build --prod

# # Step 2: Build Go backend
# FROM golang:1.24-alpine AS backend
# WORKDIR /app/backend
# COPY backend/ .
# RUN go build -o app .

# # Step 3: Final Image with Nginx + Go + Frontend
# FROM nginx:alpine

# # Copy built Angular files to Nginx's web root
# COPY --from=frontend /app/frontend/dist/frontend/browser /usr/share/nginx/html/

# # Replace Nginx default config with custom config
# COPY frontend/nginx.conf /etc/nginx/nginx.conf

# # Copy the compiled Go binary
# COPY --from=backend /app/backend/app /app/app
# COPY .env /app/.env

# # Copy startup script
# COPY start.sh /start.sh
# RUN chmod +x /start.sh

# # Expose HTTP port
# EXPOSE 8081

# # Start Go backend and Nginx
# CMD ["/start.sh"]



# Step 1: Build Angular frontend
FROM node:20-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/ .
RUN npm install && npm run build --prod

# Step 2: Build Go backend
FROM golang:1.24-alpine AS backend
WORKDIR /app/backend
COPY backend/ .
RUN go build -o app .

# Step 3: Final Image with Nginx + Go + Frontend
FROM nginx:alpine

# Install curl for debugging (optional)
RUN apk add --no-cache curl

# Set working directory
WORKDIR /app

# Copy built Angular files to Nginx's web root
COPY --from=frontend /app/frontend/dist/frontend/browser /usr/share/nginx/html/

# Replace Nginx default config with custom config
COPY frontend/nginx.conf /etc/nginx/nginx.conf

# Copy the compiled Go binary
COPY --from=backend /app/backend/app /app/app

# Optional: copy .env file if your Go code uses godotenv
COPY .env /app/.env

# Copy startup script
COPY start.sh /start.sh
RUN chmod +x /start.sh

# Expose the port served by Nginx
EXPOSE 8081

# Start Go backend and Nginx
CMD ["/start.sh"]

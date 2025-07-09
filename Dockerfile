# Build Angular frontend
FROM node:20-alpine as frontend-builder
WORKDIR /frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ .
RUN npm run build -- --configuration production

# Build Go backend
FROM golang:1.24-alpine as backend-builder
WORKDIR /backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN go build -o server .

# Final stage: merge both builds
FROM nginx:alpine

# Copy Angular app to Nginx
COPY frontend/nginx.conf /etc/nginx/nginx.conf
COPY --from=frontend-builder /frontend/dist/ /usr/share/nginx/html/

# Copy backend binary
COPY --from=backend-builder /backend/server /app/server

# Install tini to handle PID 1 correctly
RUN apk add --no-cache tini

# Custom entrypoint: run backend & nginx together
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 80
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/entrypoint.sh"]

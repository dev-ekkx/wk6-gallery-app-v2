#!/bin/sh

# Start Go API server in background
/app/app &

# Start Nginx in foreground
nginx -g "daemon off;"

#!/bin/sh

# Start Go backend in background
/app/server &

# Start Nginx (frontend)
nginx -g "daemon off;"

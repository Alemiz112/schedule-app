#!/bin/sh
set -e

echo "Applying runtime config..."
envsubst < /app/frontend/dist/config.js > /tmp/config.js
mv /tmp/config.js /app/frontend/dist/config.js

echo "Starting server..."
exec ./server -release=true

#!/bin/bash

# For Execute
# sed -i 's/\r$//' setup.sh && bash setup.sh

echo "[setup.sh] Start to execute..."

echo "--> Preparing directories..."
mkdir -p logs
touch logs/golang.log logs/golang-error.log 2>/dev/null || true

echo "--> Setting permissions..."
sudo chown -R www:www . 2>/dev/null || true
sudo find . \( -path ./node_modules -o -path ./vendor \) -prune -o -type d -exec chmod 755 {} \;
sudo find . \( -path ./node_modules -o -path ./vendor \) -prune -o -type f -exec chmod 644 {} \;

echo "--> Setting writable permissions..."
sudo chmod 755 runner-app

echo "--> Supervisor setup..."
sudo cp ./deploy/supervisor.conf /etc/supervisor/conf.d/simple_pos_api_novaardiansyah_id.conf

sudo supervisorctl reread
sudo supervisorctl update
sudo supervisorctl restart simple_pos_api_novaardiansyah_id

echo "--> Securing credentials files..."
sudo chmod 600 .env .env.local .env.production .well-known .git artisan Makefile setup.sh 2>/dev/null

echo "[setup.sh] Script has been executed successfully..."

#!/bin/bash

GCR_SERVER=gcr.io/messaging-1174/portal-server
GCR_DB=gcr.io/messaging-1174/portal/db

echo "Stopping existing Docker containers..."

docker rm -f portal_db
docker rm -f portal_api
docker rm -f portal_gcm

echo "Fetching latest images..."

gcloud docker pull $GCR_DB:latest
gcloud docker pull $GCR_SERVER:latest

echo "Starting the portal database..."

docker run -d --restart=always --name portal_db \
           -v /var/lib/postgresql/data:/var/lib/postgresql/data \
           -p 5432:5432 $GCR_DB

setup_db() {
    echo "Setting up portal database..."
    docker exec -u postgres portal_db ./setup.sh

    echo "Performing database operation: $1..."
    docker run --rm --link portal_db:postgres $GCR_SERVER ./dbtool $1

    echo "Setting up database permissions..."
    docker exec -u postgres portal_db ./permissions.sh
}

if docker logs portal_db 2>/dev/null | grep -q 'will be initialized'; then
    until docker logs portal_db 2>/dev/null | grep -q 'init process complete'; do
        sleep 1
    done
    setup_db create
else
    setup_db migrate
fi

echo "Starting Portal GCM server..."
docker run -d --restart=always --name portal_gcm \
           --link portal_db:postgres \
           $GCR_SERVER ./portalgcm

echo "Starting Portal API on port 8080..."
docker run -d --restart=always --name portal_api \
           --link portal_db:postgres -p 8080:8080 \
           $GCR_SERVER ./portalapi

echo "Successfully deployed containers"

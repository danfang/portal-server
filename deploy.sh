#!/bin/bash

DB_CONTAINER=gcr.io/messaging-1174/portal/db
API_CONTAINER=gcr.io/messaging-1174/portal/api
GCM_CONTAINER=gcr.io/messaging-1174/portal/gcm

echo "Stopping existing Docker containers..."

docker rm -f portal_db
docker rm -f portal_api
docker rm -f portal_gcm

echo "Fetching latest images..."

gcloud docker pull $DB_CONTAINER:latest
gcloud docker pull $API_CONTAINER:latest
gcloud docker pull $GCM_CONTAINER:latest

echo "Starting the portal database..."

docker run -v /var/lib/postgresql/data:/var/lib/postgresql/data -d --name portal_db $DB_CONTAINER

setup_db() {
    echo "Setting up portal database..."
    docker exec -u postgres portal_db ./setup.sh

    echo "Performing database operation: $1..."
    docker run --rm --name portal_api --link portal_db:postgres $API_CONTAINER ./dbtool $1

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
docker run -d --name portal_gcm --link portal_db:postgres $GCM_CONTAINER

echo "Running Portal API on port 8080 (Ctrl-p Ctrl-q to daemonize)"
docker run -p 8080:8080 -ti --name portal_api --link portal_db:postgres $API_CONTAINER

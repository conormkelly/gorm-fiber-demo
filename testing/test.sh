#!/bin/bash

docker compose --profile tests up --detach
echo "Sleeping for 10 seconds..." && sleep 10 # TODO: find a better way of making this more consistent
docker compose logs newman
docker compose down

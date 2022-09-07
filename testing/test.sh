#!/bin/bash

docker-compose up -d && sleep 0.5
docker-compose logs newman
docker-compose down

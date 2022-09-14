#!/bin/bash

docker-compose up -d && sleep 5 # increase sleep timeout if not seeing postman tests in logs
docker-compose logs newman
docker-compose down

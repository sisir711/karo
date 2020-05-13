#!/usr/bin/env bash

echo "Building karo-api ..."
cd $1 && echo $PWD
docker-compose down #-v db_data
docker-compose up -d 
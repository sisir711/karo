#!/usr/bin/env bash

echo "Building karo-cli ... $1"
cd $1 && go build

#Logistics
mv "$1" ../build/release
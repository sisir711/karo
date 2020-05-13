#!/usr/bin/env bash

echo "Building karo-workforce ... $1"
cd $1 && go build

#Logistics
mv "workforce-go" ../build/release
#!/bin/bash

while true; do
    curl -X POST 'http://ubuntu-1/api/v2/root/0/pan/255'
    sleep 10
    curl -X POST 'http://ubuntu-1/api/v2/root/0/pan/0'
    sleep 10
done
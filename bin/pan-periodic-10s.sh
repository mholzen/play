#!/bin/bash

resource=1
period=10

while true; do
    echo 255 | curl --data @- -X POST "http://ubuntu-1/api/v2/root/$resource/pan"
    sleep $period
    echo 0 | curl --data @- -X POST "http://ubuntu-1/api/v2/root/$resource/pan"
    sleep $period
done
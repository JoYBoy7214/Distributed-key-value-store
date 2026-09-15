#!/bin/bash

trap "kill 0" SIGINT

cd /home/gowtham/development/GO/Distributed_key_value_store/proxy

go build -o proxy-bin || exit



./proxy-bin &
PROXY_PID=$!
echo "Proxy started with PID: $PROXY_PID"
sleep 1

cd ..
cd node

go build -o node-bin || exit



for ((i=8081;i<8084;i++)); do
    ./node-bin -port $i &
    echo "Node on port $i is running with PID: $!"
    
done 

wait 



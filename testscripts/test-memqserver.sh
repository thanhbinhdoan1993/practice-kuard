#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

HOST=127.0.0.1:8080
URLPREFIX=/memq/server

curl -X PUT ${HOST}${URLPREFIX}/queues/work
curl -X POST ${HOST}${URLPREFIX}/queues/work/enqueue -d "message 1"
curl -X POST ${HOST}${URLPREFIX}/queues/work/enqueue -d "message 2"
curl -X POST ${HOST}${URLPREFIX}/queues/work/enqueue -d "message 3"
curl ${HOST}${URLPREFIX}/stats
curl -X POST ${HOST}${URLPREFIX}/queues/work/dequeue
curl -X POST ${HOST}${URLPREFIX}/queues/work/dequeue
curl -X POST ${HOST}${URLPREFIX}/queues/work/dequeue
curl -X POST ${HOST}${URLPREFIX}/queues/work/enqueue -d "message 1"
curl -X POST ${HOST}${URLPREFIX}/queues/work/enqueue -d "message 2"
curl -X POST ${HOST}${URLPREFIX}/queues/work/enqueue -d "message 3"
curl -X POST ${HOST}${URLPREFIX}/queues/work/drain
curl ${HOST}${URLPREFIX}/stats
curl -X DELETE ${HOST}${URLPREFIX}/queues/work

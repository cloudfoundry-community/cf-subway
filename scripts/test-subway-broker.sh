#!/bin/bash

set -e
set -x

sudo apt-get update
sudo apt-get install curl uuid -y

curl -f -X GET https://warreng:natedogg@haash-broker-1.cfapps.io/v2/catalog
curl -f -X GET https://username:password@subway-broker.cfapps.io/v2/catalog

id=$(uuid)
curl -f -X PUT https://username:password@subway-broker.cfapps.io/v2/service_instances/$id -d '{"plan_id": "1", "service_id":"1"}'
curl -f -X PUT https://username:password@subway-broker.cfapps.io/v2/service_instances/$id/service_bindings/$id -d '{"plan_id": "1", "service_id":"1", "app_guid": "x"}'
# curl -fv -X DELETE https://username:password@subway-broker.cfapps.io/v2/service_instances/$id/service_bindings/$id -d '{"plan_id": "1", "service_id":"1"}'
# curl -fv -X DELETE https://username:password@subway-broker.cfapps.io/v2/service_instances/$id -d '{"plan_id": "1", "service_id":"1"}'

id=$(uuid)
curl -f -X PUT https://username:password@subway-broker.cfapps.io/v2/service_instances/$id -d '{"plan_id": "1", "service_id":"1"}'
curl -f -X PUT https://username:password@subway-broker.cfapps.io/v2/service_instances/$id/service_bindings/$id -d '{"plan_id": "1", "service_id":"1", "app_guid": "x"}'

id=$(uuid)
curl -f -X PUT https://username:password@subway-broker.cfapps.io/v2/service_instances/$id -d '{"plan_id": "1", "service_id":"1"}'
curl -f -X PUT https://username:password@subway-broker.cfapps.io/v2/service_instances/$id/service_bindings/$id -d '{"plan_id": "1", "service_id":"1", "app_guid": "x"}'

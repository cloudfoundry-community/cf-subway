#!/bin/bash

set -e
set -x

SUBWAY_USERNAME=${SUBWAY_USERNAME:-username}
SUBWAY_PASSWORD=${SUBWAY_PASSWORD:-password}

sudo apt-get update
sudo apt-get install curl uuid -y

curl -f -X GET https://warreng:natedogg@haash-broker-1.cfapps.io/v2/catalog
curl -f -X GET https://$SUBWAY_USERNAME:$SUBWAY_PASSWORD@subway-broker.cfapps.io/v2/catalog

id=$(uuid)
curl -f -X PUT https://$SUBWAY_USERNAME:$SUBWAY_PASSWORD@subway-broker.cfapps.io/v2/service_instances/$id -d '{"plan_id": "1", "service_id":"1"}'
curl -f -X PUT https://$SUBWAY_USERNAME:$SUBWAY_PASSWORD@subway-broker.cfapps.io/v2/service_instances/$id/service_bindings/$id -d '{"plan_id": "1", "service_id":"1", "app_guid": "x"}'
# curl -fv -X DELETE https://$SUBWAY_USERNAME:$SUBWAY_PASSWORD@subway-broker.cfapps.io/v2/service_instances/$id/service_bindings/$id -d '{"plan_id": "1", "service_id":"1"}'
# curl -fv -X DELETE https://$SUBWAY_USERNAME:$SUBWAY_PASSWORD@subway-broker.cfapps.io/v2/service_instances/$id -d '{"plan_id": "1", "service_id":"1"}'

id=$(uuid)
curl -f -X PUT https://$SUBWAY_USERNAME:$SUBWAY_PASSWORD@subway-broker.cfapps.io/v2/service_instances/$id -d '{"plan_id": "1", "service_id":"1"}'
curl -f -X PUT https://$SUBWAY_USERNAME:$SUBWAY_PASSWORD@subway-broker.cfapps.io/v2/service_instances/$id/service_bindings/$id -d '{"plan_id": "1", "service_id":"1", "app_guid": "x"}'

id=$(uuid)
curl -f -X PUT https://$SUBWAY_USERNAME:$SUBWAY_PASSWORD@subway-broker.cfapps.io/v2/service_instances/$id -d '{"plan_id": "1", "service_id":"1"}'
curl -f -X PUT https://$SUBWAY_USERNAME:$SUBWAY_PASSWORD@subway-broker.cfapps.io/v2/service_instances/$id/service_bindings/$id -d '{"plan_id": "1", "service_id":"1", "app_guid": "x"}'

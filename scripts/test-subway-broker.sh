#!/bin/bash

set -e
set -x

SUBWAY_USERNAME=${SUBWAY_USERNAME:-username}
SUBWAY_PASSWORD=${SUBWAY_PASSWORD:-password}
SUBWAY_HOST=${SUBWAY_HOST}

sudo apt-get update
sudo apt-get install curl uuid -y

curl -f -X GET https://$BACKEND/v2/catalog
curl -f -X GET https://$SUBWAY_USERNAME:$SUBWAY_PASSWORD@$SUBWAY_HOST/v2/catalog

id=$(uuid)
curl -f -X PUT https://$SUBWAY_USERNAME:$SUBWAY_PASSWORD@$SUBWAY_HOST/v2/service_instances/$id -d "{\"plan_id\": \"$PLAN_ID\", \"service_id\":\"$SERVICE_ID\"}"
curl -f -X PUT https://$SUBWAY_USERNAME:$SUBWAY_PASSWORD@$SUBWAY_HOST/v2/service_instances/$id/service_bindings/$id -d "{\"plan_id\": \"$PLAN_ID\", \"service_id\":\"$SERVICE_ID\", \"app_guid\": \"x\"}"
# curl -fv -X DELETE https://$SUBWAY_USERNAME:$SUBWAY_PASSWORD@$SUBWAY_HOST/v2/service_instances/$id/service_bindings/$id -d "{\"plan_id\": \"$PLAN_ID\", \"service_id\":\"$SERVICE_ID\"}"
# curl -fv -X DELETE https://$SUBWAY_USERNAME:$SUBWAY_PASSWORD@$SUBWAY_HOST/v2/service_instances/$id -d "{\"plan_id\": \"$PLAN_ID\", \"service_id\":\"$SERVICE_ID\"}"

id=$(uuid)
curl -f -X PUT https://$SUBWAY_USERNAME:$SUBWAY_PASSWORD@$SUBWAY_HOST/v2/service_instances/$id -d "{\"plan_id\": \"$PLAN_ID\", \"service_id\":\"$SERVICE_ID\"}"
curl -f -X PUT https://$SUBWAY_USERNAME:$SUBWAY_PASSWORD@$SUBWAY_HOST/v2/service_instances/$id/service_bindings/$id -d "{\"plan_id\": \"$PLAN_ID\", \"service_id\":\"$SERVICE_ID\", \"app_guid\": \"x\"}"

id=$(uuid)
curl -f -X PUT https://$SUBWAY_USERNAME:$SUBWAY_PASSWORD@$SUBWAY_HOST/v2/service_instances/$id -d "{\"plan_id\": \"$PLAN_ID\", \"service_id\":\"$SERVICE_ID\"}"

echo create binding to subway should be same as to backend
subway=$(curl -f -X PUT https://$SUBWAY_USERNAME:$SUBWAY_PASSWORD@$SUBWAY_HOST/v2/service_instances/$id/service_bindings/$id -d "{\"plan_id\": \"$PLAN_ID\", \"service_id\":\"$SERVICE_ID\", \"app_guid\": \"x\"}")
echo $subway
backend=$(curl -f -X PUT https://$BACKEND/v2/service_instances/$id/service_bindings/$id -d "{\"plan_id\": \"$PLAN_ID\", \"service_id\":\"$SERVICE_ID\", \"app_guid\": \"x\"}")
echo $backend
[[ $subway == $backend ]]

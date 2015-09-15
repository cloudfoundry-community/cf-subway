#!/bin/bash

if [[ "${BOSH_TARGET}X" == "X" ]]; then
  echo "Required \$BOSH_TARGET, \$BOSH_USERNAME, \$BOSH_PASSWORD"
  exit 1
fi

BOSH_DEPLOYMENT_NAME=${BOSH_DEPLOYMENT_NAME:-subway-postgresql-docker}

cat > $HOME/.bosh_config << YAML
---
auth:
  ${BOSH_TARGET}:
    username: ${BOSH_USERNAME}
    password: ${BOSH_PASSWORD}
YAML

bosh target $BOSH_TARGET

if [[ "$(bosh releases | grep ' docker ')X" == "X" ]]; then
  bosh upload release https://bosh.io/d/github.com/cf-platform-eng/docker-boshrelease
fi
if [[ "$(bosh releases | grep ' postgresql-docker ')X" == "X" ]]; then
  bosh upload release https://bosh.io/d/github.com/cloudfoundry-community/postgresql-docker-boshrelease
fi

bosh -n delete deployment ${BOSH_DEPLOYMENT_NAME}

set -e

cd /tmp
git clone https://github.com/cloudfoundry-community/postgresql-docker-boshrelease.git postgresql-docker
cd postgresql-docker

mkdir -p tmp
cat > tmp/scaling.yml << YAML
---
name: ${BOSH_DEPLOYMENT_NAME}
update:
  canaries: 0

jobs:
  - name: postgresql_docker_z1
    instances: 3
YAML


./templates/make_manifest warden broker embedded tmp/scaling.yml
bosh -n deploy

backend_ips=$(bosh vms ${BOSH_DEPLOYMENT_NAME} | grep running | awk  '{print $8}')
cd -

counter=1
for ip in $backend_ips; do
  echo "export BACKEND_BROKER_${counter}=http://containers:containers@${ip}" >> backends.env
  echo "http://containers:containers@${ip}" >> backends
  counter=$((counter + 1))
done

Subway
======

A subway to scale out Cloud Foundry Service Brokers.

```
________   ______________________>__
[]_[]||[| |]||[]_[]_[]|||[]_[]_[]||[|
===o-o==/_\==o-o======o-o======o-o==/______
:::::::::::::::::::::::::::::::::::::::::::
```

Subway is a multiplexing service broker that allows you to scale out single node brokers, such as cf-containers-broker/docker-boshrelease and cf-redis-broker/cf-redis-release.

-	CI pipeline https://ci.starkandwayne.com/pipelines/subway

Deployment
----------

The subway broker can be deployed as a Cloud Foundry application. It needs to be given connection information to the 1+ backend brokers via environment variables.

For example, consider below that you are deploying Subway to tunnel to 1+ existing single-node Postgresql backends (say deployed via [postgresql-docker-boshrelease](https://github.com/cloudfoundry-community/postgresql-docker-boshrelease)\):

```
git clone https://github.com/cloudfoundry-community/cf-subway
cd cf-subway
cf push subway-postgresql-docker --no-start
broker_url=$(cf app subway-p-redis | grep urls | awk '{print $2}')
```

Now set one environment variable for each backend Postgresql node. The variable must start with `BACKEND_BROKER`.

```
cf set-env subway-postgresql-docker BACKEND_BROKER_1 http://containers:containers@10.10.10.10
cf set-env subway-postgresql-docker BACKEND_BROKER_2 http://containers:containers@10.10.10.11
cf set-env subway-postgresql-docker BACKEND_BROKER_3 http://containers:containers@10.10.10.12
```

Now start the Subway broker app:

```
cf start subway-postgresql-docker
```

Finally, register the broker with Cloud Foundry (requires you to login as an admin at the moment).

```
cf create-service-broker subway-postgresql-docker username password ${broker_url}
```

Finally:

-	Run `cf service-access` to discover the available services & plans
-	Run `cf enable-service-access` to enable access to everyone, or specific plans/org combinations

### Usage

You should now be able to see the service/plans in the marketplace:

```
cf marketplace
```

And be able to create/bind/unbind/deprovision service instances:

```
cf create-service postgresql shared test-pg
cf delete-service postgresql shared test-pg
```

### Update backends

To update the location or number of backend brokers you simply:

-	update/add the variables - `cf set-env subway-postgresql-docker BACKEND_BROKER_XYZ` variables
-	restart the subway app - `cf restart subway-postgresql-docker`

Development
-----------

To update dependencies:

```

godep save ./... git add Godeps git commit -m "update deps"

```

To run tests locally:

```

ginkgo \*

```

To run tests within Concourse:

```

fly e

```

CI
--

The CI pipeline is publicly visible at https://ci.starkandwayne.com/pipelines/subway

To update CI pipeline in Concourse:

```
fly -t snw c -c pipeline.yml subway
```

Thanks
------

-	ASCII art for Subway tram is from http://www.retrojunkie.com/asciiart/vehicles/trains.htm

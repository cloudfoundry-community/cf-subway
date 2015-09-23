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
-	Blog post intro https://blog.starkandwayne.com/2015/09/21/how-to-scale-out-any-cloud-foundry-service/

How does it work?
-----------------

When an end user runs `cf create-service` the Provision request goes to the Subway broker; which randomly assigns it to backend brokers until one of them accepts it (backend brokers might not have capacity for the requested service plan).

When the end user runs `cf bind-service` the Bind request goes to the Subway broker; which forwards the Bind request on to the correct backend broker; and returns the resulting binding credentials to Cloud Controller.

Applications bound to services then communicate directly with the backend service - neither Subway nor the Backend Broker is involved in direct application-to-service communication.

Similarly to Bind, when the end user runs `cf unbind-service` or `cf delete-service` the request goes to the Subway broker, which forwards the request to the correct broker.

As an aside, the implementation of "forwards the request to the correct broker" involves first sending the requests to random incorrect brokers until the correct broker is discovered. Subway is stateless and does not remember how it assigned service instances to backend brokers.

![diagram](https://www.gliffy.com/go/publish/image/8949413/L.png)

Deployment
----------

The subway broker can be deployed as a Cloud Foundry application. It needs to be given connection information to the 1+ backend brokers via environment variables.

For example, consider below that you are deploying Subway to tunnel to 1+ existing single-node Postgresql backends (say deployed via [postgresql-docker-boshrelease](https://github.com/cloudfoundry-community/postgresql-docker-boshrelease)\):

```
git clone https://github.com/cloudfoundry-community/cf-subway
cd cf-subway
appname=subway-postgresql-docker
cf push $appname --no-start
cf set-env $appname SUBWAY_USERNAME secretusername
cf set-env $appname SUBWAY_PASSWORD secretpassword
broker_url=$(cf app $appname | grep urls | awk '{print $2}')
```

The credentials `secretusername` and `secretpassword` are used later when registering the broker (`cf create-service-broker`). By default they are `username` and `password` respectively; if you forget to explicit set them.

Now set one environment variable for each backend Postgresql node. The variable must start with `BACKEND_BROKER`.

```
cf set-env $appname BACKEND_BROKER_1 http://containers:containers@10.10.10.10
cf set-env $appname BACKEND_BROKER_2 http://containers:containers@10.10.10.11
cf set-env $appname BACKEND_BROKER_3 http://containers:containers@10.10.10.12
```

Now start the Subway broker app:

```
cf start $appname
```

Finally, register the broker with Cloud Foundry (requires you to login as an admin at the moment).

```
cf create-service-broker postgresql-docker secretusername secretpassword http://${broker_url}
```

If you are replacing an existing broker with Subway then you will run:

```
cf update-service-broker postgresql-docker secretusername secretpassword http://${broker_url}
```

Finally:

-	Run `cf service-access` to discover the available services & plans
-	Run `cf enable-service-access` to enable access to everyone, or specific plans/org combinations

If your backend brokers update their catalog, then Subway will automatically pick up the new catalog when you next `cf update-service-broker`.

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
godep save ./...
git add Godeps
git commit -m "update deps"
```

To run tests locally:

```
ginkgo -r
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

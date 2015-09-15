Subway
======

A subway to scale out Cloud Foundry Service Brokers.

Subway is a multiplexing service broker that allows you to scale out single node brokers, such as cf-containers-broker/docker-boshrelease and cf-redis-broker/cf-redis-release.

-	CI pipeline https://ci.starkandwayne.com/pipelines/subway

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
ginkgo *
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

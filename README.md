Subway for Cloud Foundry Service Brokers
========================================

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

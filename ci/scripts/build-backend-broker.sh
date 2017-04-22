#!/bin/bash

cd haash-broker

export TERM=${TERM:-dumb}
./gradlew assemble

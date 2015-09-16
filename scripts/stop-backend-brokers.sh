#!/bin/bash

set -e

cf login -a ${api} -u ${username} -p ${password} -o ${organization} -s ${space}

cf a | grep started | awk '{print $1}' | xargs -L1 cf stop

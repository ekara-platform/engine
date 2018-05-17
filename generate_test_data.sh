#!/bin/bash

(
    cd `dirname $0`/testdata
    mkdir -p components/lagoon-platform/core
    (
        cd components/lagoon-platform/core
        git init
        echo "v1.0.0" > VERSION
        git add .
        git commit -m "Content for v1.0.0"
        git tag v1.0.0
        echo "v1.0.1" > VERSION
        git commit -a -m "Content for v1.0.1"
        git tag v1.0.1
        echo "v2.0.0" > VERSION
        git commit -a -m "Content for v2.0.0"
        git tag v2.0.0
        echo "v2.0.1" > VERSION
        git commit -a -m "Content for v2.0.1"
    )
)
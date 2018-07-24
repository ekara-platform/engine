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
        git tag stable
        echo "v2.0.1" > VERSION
        git commit -a -m "Content for v2.0.1"
    )

    mkdir -p components/lagoon-platform/aws-provider
    (
        cd components/lagoon-platform/aws-provider
        git init
        echo "v1.0.0" > VERSION
        git add .
        git commit -m "Content for v1.0.0"
        git tag v1.0.0
        echo "v1.1.0" > VERSION
        git commit -a -m "Content for v1.1.0"
        git tag v1.1.0
        echo "v1.2.0" > VERSION
        git commit -a -m "Content for v1.2.0"
        git tag v1.2.0
        echo "v1.2.1" > VERSION
        git commit -a -m "Content for v1.2.1"
    )

    mkdir -p components/lagoon-platform/swarm-orchestrator
    (
        cd components/lagoon-platform/swarm-orchestrator
        git init
        echo "v1.0.0" > VERSION
        git add .
        git commit -m "Content for v1.0.0"
        git tag v1.0.0
        echo "v1.1.0" > VERSION
        git commit -a -m "Content for v1.1.0"
        git tag v1.1.0
        echo "v1.2.0" > VERSION
        git commit -a -m "Content for v1.2.0"
        git tag v1.2.0
        echo "v1.2.1" > VERSION
        git commit -a -m "Content for v1.2.1"
    )

    mkdir -p sample
    (
        cd sample
        git init
        git add .
        git commit -m "Content for v1.0.0"
        git checkout -b test
        git checkout master
        git tag v1.0.0
    )
)
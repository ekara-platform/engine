#!/bin/bash

mkdir -p `dirname $0`/testdata
(
    cd `dirname $0`/testdata

    mkdir -p components/ekara-platform/distribution
    (
        cd components/ekara-platform/distribution
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

    mkdir -p components/ekara-platform/aws-provider
    (
        cd components/ekara-platform/aws-provider
        git init
        echo "1.0.0" > VERSION
        git add .
        git commit -m "Content for 1.0.0"
        git tag v1.0.0
        echo "1.1.0" > VERSION
        git commit -a -m "Content for 1.1.0"
        git tag v1.1.0
        echo "1.2.0" > VERSION
        git commit -a -m "Content for 1.2.0"
        git tag v1.2.0
        echo "1.2.1" > VERSION
        git commit -a -m "Content for 1.2.1"
    )

    mkdir -p components/ekara-platform/swarm-orchestrator
    (
        cd components/ekara-platform/swarm-orchestrator
        git init
        echo "1.0.0" > VERSION
        git add .
        git commit -m "Content for 1.0.0"
        git tag v1.0.0
        echo "1.1.0" > VERSION
        git commit -a -m "Content for 1.1.0"
        git tag v1.1.0
        echo "1.2.0" > VERSION
        git commit -a -m "Content for 1.2.0"
        git tag v1.2.0
        echo "1.2.1" > VERSION
        git commit -a -m "Content for 1.2.1"
    )

    mkdir -p components/ekara-platform/monitoring-stack
    (
        cd components/ekara-platform/monitoring-stack
        git init
        echo "1.0.0" > VERSION
        git add .
        git commit -m "Content for 1.0.0"
        git tag v1.0.0
        echo "1.1.0" > VERSION
        git commit -a -m "Content for 1.1.0"
        git tag v1.1.0
        echo "1.2.0" > VERSION
        git commit -a -m "Content for 1.2.0"
        git tag v1.2.0
        echo "1.2.1" > VERSION
        git commit -a -m "Content for 1.2.1"
    )

    mkdir -p sample
    (
        cd sample
	    rm -rf .git
        git init
        git add .
        git commit -m "Content for 1.0.0"
        git checkout -b test
        git checkout master
        git tag v1.0.0
    )
)
mkdir -p `dirname $0`/component/testdata
(
    cd `dirname $0`/component/testdata

    mkdir -p components/ekara-platform/c1
    (
        cd components/ekara-platform/c1
        git init
        mkdir modules
        echo "DUMMY" > modules/dummy
        git add .
        git commit -m "Content for 1.0.0"
        git tag v1.0.0
    )
    mkdir -p components/ekara-platform/c2
    (
        cd components/ekara-platform/c2
        git init
        mkdir inventory
        echo "DUMMY" > inventory/dummy
        git add .
        git commit -m "Content for 1.0.0"
        git tag v1.0.0
    )
    mkdir -p components/ekara-platform/c3
    (
        cd components/ekara-platform/c3
        git init
        mkdir modules
        echo "DUMMY" > modules/dummy
        mkdir inventory
        echo "DUMMY" > inventory/dummy
        git add .
        git commit -m "Content for 1.0.0"
        git tag v1.0.0
    )
    mkdir -p components/ekara-platform/c4
    (
        cd components/ekara-platform/c4
        git init
        echo "DUMMY" > dummy
        git add .
        git commit -m "Content for 1.0.0"
        git tag v1.0.0
    )
)

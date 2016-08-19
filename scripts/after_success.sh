#!/usr/bin/env bash

bash <(curl -s https://codecov.io/bash)
if [ "TRAVIS_GO_VERSION" == "1.7" ]; then
	docker login -u $DOCKER_USER -p $DOCKER_PASSWORD
	if [ ! -z "$TRAVIS_TAG" ] && [ "$TRAVIS_PULL_REQUEST" == "false" ]; then
		export VCS_REF=$TRAVIS_TAG
		make publish
	fi
	if [ $TRAVIS_BRANCH == 'master' ] && [ "$TRAVIS_PULL_REQUEST" == "false" ]; then
		export VERSION=latest
		export VCS_REF=master
		make publish
	fi
fi
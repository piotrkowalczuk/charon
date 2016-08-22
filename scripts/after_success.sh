#!/usr/bin/env bash

bash <(curl -s https://codecov.io/bash)

echo $TRAVIS_TAG
echo $TRAVIS_PULL_REQUEST
echo $TRAVIS_BRANCH

if [ "TRAVIS_GO_VERSION" == "1.7" ]; then
	docker login -u $DOCKER_USER -p $DOCKER_PASSWORD
	if [ ! -z "$TRAVIS_TAG" ]; then
		export VCS_REF=$TRAVIS_TAG
		make publish
	fi
	if [ $TRAVIS_BRANCH == 'master' ] && [ "$TRAVIS_PULL_REQUEST" == "false" ]; then
		export VERSION=latest
		export VCS_REF=master
		make publish
	fi
fi
# Charon [![Build Status](https://travis-ci.org/piotrkowalczuk/charon.svg?branch=master)](https://travis-ci.org/piotrkowalczuk/charon)

[![GoDoc](https://godoc.org/github.com/piotrkowalczuk/charon?status.svg)](http://godoc.org/github.com/piotrkowalczuk/charon)
[![Docker Pulls](https://img.shields.io/docker/pulls/piotrkowalczuk/charon.svg?maxAge=604800)](https://hub.docker.com/r/piotrkowalczuk/charon/)
[![codecov](https://codecov.io/gh/piotrkowalczuk/charon/branch/master/graph/badge.svg)](https://codecov.io/gh/piotrkowalczuk/charon)
[![Code Climate](https://codeclimate.com/github/piotrkowalczuk/charon/badges/gpa.svg)](https://codeclimate.com/github/piotrkowalczuk/charon)
[![Go Report Card](https://goreportcard.com/badge/github.com/piotrkowalczuk/charon)](https://goreportcard.com/report/github.com/piotrkowalczuk/charon)
[![pypi](https://img.shields.io/pypi/v/charon-client.svg)](https://pypi.python.org/pypi/charon-client)

<img src="/data/logo/charon.png?raw=true" width="300">

## Quick Start

### Installation

```bash
$ go install github.com/piotrkowalczuk/charon/cmd/charond
$ go install github.com/piotrkowalczuk/charon/cmd/charonctl
```

### Superuser

```bash
$ charonctl register -address=localhost:8080 -auth.disabled -register.superuser=true -register.username="j.snow@gmail.com" -register.password=123 -register.firstname=John -register.lastname=Snow
```
## Example

### Python

```python
from  charonrpc import auth_pb2, auth_pb2_grpc
import grpc

charonChannel = grpc.insecure_channel('ADDRESS')
auth = auth_pb2_grpc.AuthStub(charonChannel)
try:
	res = auth.Login(auth_pb2.LoginRequest(
		username="USERNAME",
		password="PASSWORD",
	))

	print "access token: %s" % res.value
except grpc.RpcError as e:
	if e.code() == grpc.StatusCode.UNAUTHENTICATED:
		print "login failure, username and/or password do not match"
	else:
	    print "grpc error: %s" % e
except Exception as e:
	print "unexpected error: %s" % e
```

## Contribution

@TODO

### Documentation

@TODO

### TODO
- [x] Auth
    - [x] login
    - [x] logout
    - [x] is authenticated
    - [x] subject
    - [x] is granted
    - [x] belongs to
- [x] Permission
	- [x] get
    - [x] list
    - [x] register
- [x] Group
    - [x] get
    - [x] list
    - [x] modify
    - [x] delete
    - [x] create
    - [x] set permissions
    - [x] list permissions
- [x] User
    - [x] get
    - [x] list
    - [x] modify
    - [x] delete
    - [x] create
    - [x] set permissions
    - [x] set groups
    - [x] list permissions
    - [x] list groups

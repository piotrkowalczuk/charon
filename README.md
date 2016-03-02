Charon [![GoDoc](https://godoc.org/github.com/piotrkowalczuk/charon?status.svg)](http://godoc.org/github.com/piotrkowalczuk/charon)&nbsp;[![Build Status](https://travis-ci.org/piotrkowalczuk/charon.svg?branch=master)](https://travis-ci.org/piotrkowalczuk/charon)&nbsp;[![codecov.io](https://codecov.io/github/piotrkowalczuk/charon/coverage.svg?branch=master)](https://codecov.io/github/piotrkowalczuk/charon?branch=master)
=============

<img src="/data/logo/charon.png?raw=true" width="300">

## Quick Start

### Installation

```bash
$ go install github.com/piotrkowalczuk/charon/charond
$ go install github.com/piotrkowalczuk/charon/charonctl
```

### Superuser

```bash
$ charonctl register -noauth -r.username="j.snow@gmail.com" -r.password=123 -r.firstname=John -r.lastname=Snow
```


## Contribution

### Documentation
Documentation is available on [charon.readme.io](http://charon.readme.io).

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
- [ ] Group
    - [x] get
    - [x] list
    - [x] modify
    - [x] delete
    - [x] create
    - [ ] set permissions
    - [ ] list permissions
- [ ] User
    - [x] get
    - [x] list
    - [x] modify
    - [x] delete
    - [x] create
    - [ ] set permissions
    - [x] set groups
    - [x] list permissions
    - [x] list groups

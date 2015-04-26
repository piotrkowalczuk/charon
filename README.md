Auth Service
=============

@TODO description

Installation
------------
1. Set you GOPATH properly (http://golang.org/doc/code.html#GOPATH)
2. `go get github.com/go-soa/auth`
3. `go get` if some dependencies are missing
4. Create `conf/{env}.xml` based on `conf/{env}.xml.dist`
5. Set `$AUTH_SERVICE_ENV` global variable to `test`, `development` or `production`

Commands
--------

#### Build
```bash
go build
```

#### Service
```bash
./auth initdb - execute data/sql/schema_{adapter}.sql against configured database.
./auth run - starts server.
./auth help [command] - display help message about available commands
```

Dependencies
------------
- PostgreSQL
- MySQL *(not supported yet)*

TODO
----
- [ ] Commands
	- [x] Initialize database
	- [x] Start server
- [ ] Views
	- [x] Registration
- [ ] REST API
	- [ ] Registration
- [ ] RPC API
	- [ ] Registration

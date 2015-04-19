Auth Service
=============

@TODO description

Installation
------------
1. Set you GOPATH properly (http://golang.org/doc/code.html#GOPATH)
2. `go get github.com/piotrkowalczuk/auth-service`
3. `go get` if some dependencies are missing
4. Create `conf/{env}.xml` based on `conf/{env}.xml.dist`
5. Set `$AUTH_SERVICE_ENV` global variable to `test`, `development` or `production`

Commands
--------

#### Build

    go build

#### Service

    ./auth-service initdb - execute data/sql/schema_{env}.sql against configured database.
    ./auth-service help [command] - display help message about available commands
    
Dependencies
------------
- PostgreSQL
- MySQL *(not supported yet)*
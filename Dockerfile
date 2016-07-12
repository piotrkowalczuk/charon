FROM golang
MAINTAINER  Piotr Kowalczuk <p.kowalczuk.priv@gmail.com>

ADD . /go/src/github.com/piotrkowalczuk/charon

WORKDIR /go/src/github.com/piotrkowalczuk/charon

RUN make get
RUN go install github.com/piotrkowalczuk/charon/cmd/charond
RUN go install github.com/piotrkowalczuk/charon/cmd/charonctl
RUN rm -rf /go/src

EXPOSE 8080

CMD ["/go/bin/charond", "-host=0.0.0.0", "-namespace=charon", "-mnemo.address=mnemosyne:8080", "-p.address=postgres://postgres:postgres@postgres/postgres?sslmode=disable"]
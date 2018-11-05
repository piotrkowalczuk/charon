#@IgnoreInspection BashAddShebang
: ${PROTOC:="/usr/local/bin/protoc"}
SERVICE="charon"
PROTO_INCLUDE="-I=./${SERVICE}rpc -I=/usr/include -I=./vendor/github.com/piotrkowalczuk"

case $1 in
    lint)
        ${PROTOC} ${PROTO_INCLUDE} --lint_out=. ./${SERVICE}rpc/*.proto
        ;;
    python)
        python -m grpc_tools.protoc ${PROTO_INCLUDE} --python_out=. --grpc_python_out=. ./${SERVICE}rpc/*.proto
        ;;
    java)
        rm -rf ./tmp/java
        mkdir -p ./tmp/java
        ${PROTOC} ${PROTO_INCLUDE} --java_out=./tmp/java ./${SERVICE}rpc/*.proto
        ;;
    golang | go)
        ${PROTOC} ${PROTO_INCLUDE} --go_out=plugins=grpc:${GOPATH}/src ./${SERVICE}rpc/*.proto
        goimports -w ./${SERVICE}rpc
        ;;
	*)
	    echo "code generation failure due to unknown language: ${1}"
        exit 1
        ;;
esac

: ${PROTOC:="/usr/local/bin/protoc"}
PROTO_INCLUDE="-I=/usr/include -I=."

case $1 in
    lint)
        ${PROTOC} ${PROTO_INCLUDE} --lint_out=. ./mnemosynerpc/*.proto
        ;;
    python)
        python -m grpc_tools.protoc ${PROTO_INCLUDE} --python_out=. --grpc_python_out=. ./mnemosynerpc/*.proto
        ;;
    java)
        rm -rf ./tmp/java
        mkdir -p ./tmp/java
        ${PROTOC} ${PROTO_INCLUDE} --java_out=./tmp/java ./mnemosynerpc/*.proto
        ;;
    golang | go)
        ${PROTOC} ${PROTO_INCLUDE} --go_out=plugins=grpc:${GOPATH}/src ./mnemosynerpc/*.proto
        goimports -w ./mnemosynerpc
        ;;
	*)
	    echo "code generation failure due to unknown language: ${1}"
        exit 1
        ;;
esac
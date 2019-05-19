#@IgnoreInspection BashAddShebang
SERVICE="charon"
SERVICE_PATH=${GOPATH}/src/github.com/piotrkowalczuk/${SERVICE}
PROTO_INCLUDE="-I=./tmp/pb -I=./vendor/github.com/piotrkowalczuk -I=${GOPATH}/src"

: ${PROTOC:="${SERVICE_PATH}/tmp/bin/protoc"}

protobufs=(
    "rpc/${SERVICE}d/v1"
)
for protobuf in "${protobufs[@]}"
do
    case $1 in
        lint)
            ${PROTOC} ${PROTO_INCLUDE} --lint_out=. ${SERVICE_PATH}/pb/${protobuf}/*.proto
            ;;
        python)
            python -m grpc_tools.protoc ${PROTO_INCLUDE} --python_out=publish/python --grpc_python_out=publish/python ${SERVICE_PATH}/pb/${protobuf}/*.proto
            cp publish/python/github.com/piotrkowalczuk/charon/pb/${protobuf}/* publish/python/github/com/piotrkowalczuk/charon/pb/${protobuf}/
            rm -rf publish/python/github.com
            ;;
        java)
            rm -rf ./publish/java
            mkdir -p ./publish/java
            ${PROTOC} ${PROTO_INCLUDE} --java_out=publish/java ${SERVICE_PATH}/pb/${protobuf}/*.proto
            ;;
        golang | go)
            ${PROTOC} ${PROTO_INCLUDE} --go_out=plugins=grpc:${GOPATH}/src ${SERVICE_PATH}/pb/${protobuf}/*.proto
            mockery -case=underscore -dir=./pb/${protobuf} -all -outpkg=$(basename $(dirname "./pb/${protobuf}mock"))mock -output=./pb/${protobuf}mock
            goimports -w ${SERVICE_PATH}/pb
            ;;
        *)
            echo "code generation failure due to unknown language: ${1}"
            exit 1
            ;;
    esac
done
curl -L https://github.com/google/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip > protoc.zip

rm -rf ./tmp/protoc ./tmp/pb/google
mkdir -p ./tmp/protoc ./tmp/bin ./tmp/pb/google

unzip protoc.zip -d ./tmp/protoc

mv -f ./tmp/protoc/bin/protoc ./tmp/bin/protoc
mv -f ./tmp/protoc/include/google/ ./tmp/pb/google

rm -rf ./tmp/protoc

protoc --version
curl -L https://github.com/google/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip > protoc.zip

rm -rf ./tmp/protoc
mkdir -p ./tmp/protoc
unzip protoc.zip -d ./tmp/protoc

mv -f ./tmp/protoc/bin/protoc /tmp/bin/protoc

rm -rf /tmp/include/google
mv -f ./tmp/protoc/include/google/ /tmp/include/google

protoc --version
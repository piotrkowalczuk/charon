curl -L https://github.com/google/protobuf/releases/download/v3.6.1/protoc-3.6.1-linux-x86_64.zip > protoc.zip

rm -rf ./tmp/protoc
mkdir -p ./tmp/protoc
unzip protoc.zip -d ./tmp/protoc

sudo mv -f ./tmp/protoc/bin/protoc /usr/local/bin/protoc

sudo rm -rf /usr/include/google
sudo mv -f ./tmp/protoc/include/google/ /usr/include/google

protoc --version
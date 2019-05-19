rm -rf tmp/tools/$1
mkdir -p tmp/tools/$1
cd tmp/tools/$1
GO111MODULE=on go mod init
GO111MODULE=on go get $1@$2
GO111MODULE=on go install $1
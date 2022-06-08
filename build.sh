#/bin/bash

version=$(go version | awk '{print $3}' | awk '{split($0,b,".");print b[2]}')
if [ $version -lt 16 ]; then
    echo "[failed] go version should be greater than or equal to 1.16"
    exit
fi

cd $(dirname $0)

files=$(ls std/*.gsproto 2> /dev/null | wc -l)
if [ "$files" == "0" ] ;then
touch std/.gsproto  # make go complier happy for embedding at least one static file
fi

if [ ! -d "bin" ]; then
    mkdir bin
fi

go build -o bin/util util/main.go
bin/util std

go build -o bin/gsc main.go
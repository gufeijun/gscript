#/bin/bash

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
#!/bin/bash

version=$(go version | awk '{print $3}' | awk '{split($0,b,".");print b[2]}')
if [ $version -lt 16 ]; then
    echo "[failed] go version should be greater than or equal to 1.16"
    exit
fi

cd $(dirname $0)

prepare() {
    files=$(ls std/*.gsproto 2> /dev/null | wc -l)
    if [ "$files" = "0" ]; then
        touch std/.gsproto  # make go complier happy for embedding at least one static file
    fi
    
    if [ ! -d "bin" ]; then
        mkdir bin
    fi
}

build_util() {
    prepare
    go build -o bin/util util/main.go
    bin/util std
}

if [ $# -eq 0 ]; then
    build_util
    go build -o bin/gsc main.go
    exit 0
fi

# $1=arch $2=os
build() {
    go env -w GOARCH=$1
    go env -w GOOS=$2
    util_bin=util
    gsc_bin=gsc
    if [ $2 = "windows" ]; then
        util_bin=${util_bin}.exe
        gsc_bin=${gsc_bin}.exe
    fi

    go build -o bin/${gsc_bin} main.go
    cd bin
    if [ $2 = "windows" ]; then
        zip gsc_$2_$1.zip ${gsc_bin}
    else
        tar cvf gsc_$2_$1.tar.gz ${gsc_bin}
    fi
    cd ..

}

if [ $1 = "all" ]; then
    build_util
    old_os=$(go env GOOS)
    old_arch=$(go env GOARCH)

    build "amd64" "windows"
    build "arm" "windows"
    build "386" "windows"
    build "amd64" "linux"
    build "arm" "linux"
    build "386" "linux"
    build "amd64" "darwin"

    go env -w GOOS=$old_os
    go env -w GOARCH=$old_arch
elif [ $1 = "clean" ]; then
    rm bin std/*.gsproto std/.gsproto -rf
else
    echo "unkown command:" $1
    exit
fi

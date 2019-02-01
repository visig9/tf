PACKAGE=github.com/visig9/tf/cmd
NAME=tf
VERSION=`git describe`
TMPDIR=/tmp/tf-$VERSION

build_dist() {
    local goos=$1
    local goarch=$2
    local target=$3

    if [[ -z $target ]]; then
        local target=dist/$VERSION/${NAME}.${goos}-${goarch}.$VERSION
    fi

    if [[ $goos = windows ]]; then
        target=$target.exe
    fi

    GOOS=${goos} GOARCH=${goarch} go get ./...\
    && GOOS=${goos} GOARCH=${goarch} go build\
        -ldflags "-X main.version=${VERSION}"\
        -o $target\
        -v $PACKAGE\
    && echo "build $target successed"
}

build_all_dist() {
    build_dist linux 386
    build_dist linux amd64
    build_dist linux arm
}

install() {
    local install_dir=`go env GOBIN`
    if [[ -z $install_dir ]]; then
        install_dir=`go env GOPATH`/bin
    fi

    if [[ ! -d $install_dir ]]; then
        mkdir -p $install_dir
    fi

    go build\
        -ldflags "-X main.version=${VERSION}"\
        -o $NAME\
        -v $PACKAGE\
    && mv $NAME $install_dir/$NAME\
    && echo "install '$PACKAGE' as '$NAME' successed"
}

help_exit() {
    echo "Usage:"
    echo "  $0 build|install"
    echo ""
    echo "Example:"
    echo "  $0 install  # build and install locally"
    echo "  $0 build    # build everythings"
    echo "  $0 build [GOOS] [GOARCH] [TARGET_FILENAME]"
    echo ""
    echo "cancel."

    exit 1
}

if [[ $1 = build ]]; then
    if [[ -z $2 ]]; then
        build_all_dist
    elif [[ -n $2 && -n $3 ]]; then
        build_dist $2 $3 $4
    else
        help_exit
    fi
elif [[ $1 = install ]]; then
    install
else
    help_exit
fi

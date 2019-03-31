#!/bin/bash
if [ ! -f tupass ]; then
    echo "Unable to find tupass binary to put into package!"
    echo "Building tupass binary using 'make build-local'..."
    make build-local
fi

USER=$(whoami)
VERSION=$(git describe --tags)
DIR="tupass-$VERSION"

if [ -d $DIR ]; then
    echo "Cleaning up old $DIR"
    rm -rf $DIR
fi

cp -r debian $DIR
mkdir -p $DIR/usr/bin
cp tupass $DIR/usr/bin/
sed -i -e "s/Version: XYZ/Version: $VERSION/g" $DIR/DEBIAN/control

if ! [ -x "$(command -v sudo)" ]; then
    if [ $USER = "root" ]; then
        echo "Bulding package $DIR using root user"
        chown -R root:root $DIR
        dpkg --build $DIR
        rm -rf $DIR
    else
	    echo 'Unable to build .deb package: sudo was not found and you are not root!'
    fi
else 
    echo "Building package $DIR using sudo"
    sudo chown -R root:root $DIR
    dpkg --build $DIR
    sudo rm -rf $DIR
fi

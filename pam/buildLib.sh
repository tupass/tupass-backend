#!/bin/bash
#Build libtupass.so (change TUPASS_PWLIST to a password list in folder passwords if wanted)
#Requires of dependencies as done in Makefile dep-go
DIR="$(cd "$(dirname "$0")" && pwd)"
USER=$(whoami)
cd ../api
rice embed-go
cd $DIR
go build -o libtupass.so -buildmode=c-shared libtupass.go
#gcc -o libtupass-test libtupass-test.c libtupass.so

# Move neccessary files to their respective places for execution
if ! [ -x "$(command -v sudo)" ]; then
  if [ $USER = "root" ]; then
    cp libtupass.so /usr/lib/
    cp libtupass.h /usr/include/
  else
    echo "Error: sudo was not found and you are not root, please move files into appropriate locations"
    echo "libtupass.so -> /usr/lib/libtupass.so and libtupass.h -> /usr/include/libtupass.h"
  fi
else
  sudo cp libtupass.so /usr/lib/
  sudo cp libtupass.h /usr/include/
fi

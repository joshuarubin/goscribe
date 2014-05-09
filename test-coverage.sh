#!/bin/bash

set -e

echo "mode: set" > acc.out
for dir in $(find . \( -path ./Godeps -o -path ./.git \) -prune -o -type d -print); do
  if ! ls $dir/*.go &> /dev/null; then
    continue
  fi

  if [ "$WERCKER" == "true" ]; then
    BUILD_TAGS="-tags wercker"
  fi

  cmd="godep go test $BUILD_TAGS -v -coverprofile=profile.out $dir"
  echo $cmd
  eval $cmd

  if [ -f profile.out ]; then
    cat profile.out | grep -v "mode: set" >> acc.out
  fi
done

rm -f ./profile.out

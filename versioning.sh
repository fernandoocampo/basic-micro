#!/bin/bash

### Increments the part of the string
## $1: version itself
## $2: number of part: 0 – major, 1 – minor, 2 – patch

increment_version() {
  local delimiter=.
  local version=$(echo "$1" | tr -d v)
  local array=($(echo "$version" | tr $delimiter '\n'))
  array[$2]=$((array[$2]+1))
  if [ $2 -lt 2 ]; then array[2]=0; fi
  if [ $2 -lt 1 ]; then array[1]=0; fi
  echo $(local IFS=$delimiter ; echo "v${array[*]}")
}

if [ "$CURRENT_ACTION" = "feat" ]; then
    increment_version $CURRENT_VERSION 1
    exit;
fi

if [ "$CURRENT_ACTION" = "fix" ]; then
    increment_version $CURRENT_VERSION 2
    exit;
fi
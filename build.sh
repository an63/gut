#!/bin/bash

set -e

CHOICE="$1"
TARGET="$2"

# define target
if [[ -z "$TARGET" ]]; then
    TARGET="./y*"
else
    TARGET=("$TARGET")
fi

# prepare environment
source ./prepare.sh

printf "====== Begin at %s, OS: %s, Mode: %s - %s ======\n" "$(date '+%Y-%m-%d %H:%M:%S %z')" "$OS_NAME" "$CHOICE" "$TARGET"

COUNT=0
for FOLDER in $TARGET
do
    PACKAGE="${FOLDER##*/}"
    if [[ ! -d "$PACKAGE" ]]; then
        continue
    fi

    if [[ "ci" == "$CHOICE" && 0 -eq $COUNT ]]; then
        printf "\n###### Go Environment ######\n"
        go env
    fi

    printf "\n###### Working on package '%s' ######\n" "$PACKAGE"
    case "$CHOICE" in
    all)
        make fmt PACKAGE="$PACKAGE"
        make build PACKAGE="$PACKAGE"
        make test PACKAGE="$PACKAGE"
        make bench PACKAGE="$PACKAGE"
        make cover PACKAGE="$PACKAGE"
        make doc PACKAGE="$PACKAGE"
        ;;
    ci)
        make build PACKAGE="$PACKAGE"
        make test PACKAGE="$PACKAGE"
        make bench PACKAGE="$PACKAGE"
        make cover PACKAGE="$PACKAGE"
        ;;
    dev)
        make fmt PACKAGE="$PACKAGE"
        make testdev PACKAGE="$PACKAGE"
        make benchdev PACKAGE="$PACKAGE"
        ;;
    *)
        printf "Unknown build option: [%s]\n" "$CHOICE"
        exit 1
        ;;
    esac

    COUNT=$((COUNT+1))
done

printf "\n====== End at %s, Packages: %d ======\n" "$(date '+%Y-%m-%d %H:%M:%S %z')" "$COUNT"

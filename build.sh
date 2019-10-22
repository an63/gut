#!/bin/bash

CHOICE="$1"

set -eu

printf "====== Begin at %s, Mode: %s ======\n" "$(date '+%Y-%m-%d %H:%M:%S %z')" "$CHOICE"

COUNT=0
for FOLDER in ./y*
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
        make cover PACKAGE="$PACKAGE"
        make bench PACKAGE="$PACKAGE"
        ;;
    *)
        printf "Unknown build option: [%s]\n" "$CHOICE"
        exit 1
        ;;
    esac

    COUNT=$((COUNT+1))
done

printf "\n====== End at %s, Packages: %d ======\n" "$(date '+%Y-%m-%d %H:%M:%S %z')" "$COUNT"

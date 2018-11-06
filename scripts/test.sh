#!/bin/bash -e
require() { command -v "$1" &>/dev/null || { echo >&2 "Error: This script requires the command $1."; exit 1; }; }

require kobopatch
require sed

versions=src/versions/*/
if [[ ! -z "$1" ]]; then
    versions=src/versions/$1/
fi

for f in $versions; do
    version="$(basename "$f")"
    printf "Testing patches for %s\n" "$version"
    (
        echo "version: 4.11.11980";
        echo "in: /dev/null";
        echo "out: /dev/null";
        echo "log: /dev/null";
        echo "patchFormat: kobopatch";
        echo "patches:";
        for pf in src/versions/$version/*/*.yaml; do
            binary="$(basename "$(dirname "$pf")" | sed "s/.yaml//g")"
            creator="$(basename "$pf" | sed "s/.yaml//g")"
            echo "  src/versions/$version/$binary.yaml/$creator.yaml: usr/local/Kobo/$binary"
        done
    ) | kobopatch -tf "testdata/$version.tar.xz" -  || {
        echo "Test failed for $version! Stopping.";
        rm -rf "$temp" "$testtemp";
        exit 1;
    }
done
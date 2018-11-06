#!/bin/bash -e
require() { command -v "$1" &>/dev/null || { echo >&2 "Error: This script requires the command $1."; exit 1; }; }
travis_fold_start() { [[ $TRAVIS == true ]] && echo -en "travis_fold:start:$1\\r"; }
travis_fold_end() { [[ $TRAVIS == true ]] && echo -en "travis_fold:end:$1\\r"; }

require kobopatch
require sed

versions=src/versions/*/
if [[ ! -z "$1" ]]; then
    versions=src/versions/$1/
fi

for f in $versions; do
    version="$(basename "$f")"
    travis_fold_start "test.$version"
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
    travis_fold_end "test.$version"
done
#!/bin/bash -e
require() { command -v "$1" &>/dev/null || { echo >&2 "Error: This script requires the command $1."; exit 1; }; }

require xz
require zip
require tar
require grep
require sed

testtemp="$PWD/testtemp"
rm -rf "$testtemp"
mkdir "$testtemp"

for f in build/kobopatch_*.zip; do
    version="$(echo "$f" | grep -Eo '[0-9.]+[0-9]+')"
    printf "Testing built patches for version %s\n" "$version"
    
    temp="$testtemp/$version"
    mkdir "$temp"

    echo "--> Unzipping patch zip"
    unzip -qq "$f" -d "$temp"

    kobopatch -tf "testdata/$version.tar.xz" "$temp/kobopatch_$version/kobopatch.yaml" || { echo "Test failed for $version! Stopping."; rm -rf "$temp" "$testtemp"; exit 1; }

    rm -rf "$temp"
    printf "\n"
done

echo "Cleaning up"
rm -rf "$testtemp"
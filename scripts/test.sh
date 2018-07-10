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
    unzip -j "$f" "*/src/*.yaml" -d "$temp"

    echo "--> Untarring test data"
    [[ -e "testdata/$version.tar.xz" ]] || { echo "Warning: no test data for $version, skipping"; continue; }
    tar -xJvf "testdata/$version.tar.xz" -C "$temp"

    for pf in $temp/*.yaml; do
        pfn="$(basename "$pf")"
        bfn="$(echo "$pfn" | sed 's/.yaml//g')"
        echo "--> Testing patch file $pfn on $bfn"
        [[ -e "$temp/$bfn" ]] || { echo "Test failed for $version - $pfn! Missing binary $bfn in testdata. Stopping."; rm -rf "$temp" "$testtemp"; exit 1; }
        go run tools/test-patches/main.go "$temp/$pfn" "$temp/$bfn" || { echo "Test failed for $version - $pfn! Stopping."; rm -rf "$temp" "$testtemp"; exit 1; }
    done

    rm -rf "$temp"
    printf "\n"
done

echo "Cleaning up"
rm -rf "$testtemp"
#!/bin/bash -e
prepend() { sed "s/^/$1/"; }
require() { command -v "$1" &>/dev/null || { echo >&2 "Error: This script requires the command $1."; exit 1; }; }
travis_fold_start() { [[ $TRAVIS == true ]] && echo -en "travis_fold:start:$1\\r" || true; }
travis_fold_end() { [[ $TRAVIS == true ]] && echo -en "travis_fold:end:$1\\r" || true; }

require unix2dos
require zip
require wget

cd "$(dirname "$0")"
cd ..

rm -rf build temp
mkdir -p build temp

kobopatch="v0.11.0"

echo "Downloading tools"
dl="$PWD/dl/$kobopatch"
mkdir -p "$dl"
for tool in kobopatch-darwin-64bit cssextract-darwin-64bit kobopatch-linux-32bit cssextract-linux-32bit kobopatch-linux-64bit cssextract-linux-64bit koboptch-windows.exe cssextract-windows.exe; do
    [[ -f "$dl/$tool" ]] || wget --show-progress --progress=bar:force:noscroll -cqP "$dl/" "https://github.com/geek1011/kobopatch/releases/download/$kobopatch/$tool"
done
ls "$dl"

for f in src/versions/*/; do
    version="$(basename "$f")"
    travis_fold_start "build.$version"
    printf "Creating patch zip for %s\n" "$version"

    temp="$PWD/temp/$version"
    build="$PWD/build"
    mkdir -p $temp $build
    
    echo "--> Copying template"
    cp -vr src/template/* "$temp" |& prepend "    "
    chmod +x "$temp/kobopatch.sh"

    echo "--> Replacing template variables"
    for f in $temp/readme.txt $temp/kobopatch.yaml; do
        sed -i "s/{{version}}/$version/g" "${f}" |& prepend "    "
    done

    echo "--> Adding kobopatch"
    cp -v $dl/kobopatch-* $dl/koboptch-* $dl/cssextract-* $temp/bin |& prepend "    "
    chmod +x $temp/bin/*

    echo "--> Merging source files"
    for patch_file_dir in src/versions/$version/*/; do
        patch_file="$(basename "$patch_file_dir")"
        printf "    %s\n" "$patch_file"

        touch $temp/src/$patch_file
        for part in $patch_file_dir/*; do
            printf "      %s\n" "$(basename "$part")"
            cat $part >> $temp/src/$patch_file
            printf "\n\n" >> $temp/src/$patch_file
        done
    done

    echo "--> Converting unix line breaks to dos"
    unix2dos $temp/src/* $temp/*.bat $temp/*.yaml $temp/*.txt |& prepend "    "

    echo "--> Creating zip"
    pushd $temp/..
    mv $version kobopatch_$version
    zip -r "$build/kobopatch_$version.zip" kobopatch_$version |& prepend "  "
    popd

    printf "\n"
    travis_fold_end "build.$version"
done

echo "Cleaning up"
rm -rf temp
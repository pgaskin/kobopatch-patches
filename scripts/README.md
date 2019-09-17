# scripts
These scripts are used to build and test the patches. They are automatically run
on every commit.

## build.go
`build.go` builds the kobopatch zips. To use it, you need Go 1.13+ installed,
and then you can run `go run ./scripts/build.go` from the root of the repository.

There are optional arguments you can pass to change the behavior of the build:

```
Usage: go run ./scripts/build.go [options] [versions...]

Options:
  -basedir string
        base dir for files (default ".")
  -dldir string
        directory under basedir for kobopatch download (default "dl")
  -kpbin string
        kobopatch binaries to download (default "kobopatch-darwin-64bit,cssextract-darwin-64bit,kobopatch-linux-32bit,cssextract-linux-32bit,kobopatch-linux-64bit,cssextract-linux-64bit,koboptch-windows.exe,cssextract-windows.exe")
  -kprepo string
        github repo for kobopatch (default "geek1011/kobopatch")
  -kpver string
        kobopatch version (default "v0.14.0")
  -outdir string
        directory under basedir for output (default "build")
  -skipbuild
        don't actually build the patches
  -skipdl
        don't download kobopatch (use this for parallell builds)
  -srcdir string
        directory under basedir for sources (default "src")
```

## test.sh
This script only runs on Linux. It tests all the patches to make sure they apply.
You can run it with `./scripts/test.sh` from the root of the repository. Optionally,
you can provide the version to test as an argument to the script.

# scripts
These scripts are used to build and test the patches. They are automatically run on every commit.

## build.go
`build.go` builds the kobopatch zips. To use it, you need Go 1.13+ installed, and then you can run `go run ./scripts/build` from the root of the repository.

```
Usage: go run ./scripts/build [options] [versions...]

Options:
  -basedir string
        base dir for files (default ".")
  -dldir string
        directory under basedir for kobopatch download (default "dl")
  -kpbin string
        kobopatch binaries to download (default "kobopatch-darwin-64bit,cssextract-darwin-64bit,kobopatch-linux-32bit,cssextract-linux-32bit,kobopatch-linux-64bit,cssextract-linux-64bit,koboptch-windows.exe,cssextract-windows.exe")
  -kprepo string
        github repo for kobopatch (default "pgaskin/kobopatch")
  -kpver string
        kobopatch version (default "v0.15.0")
  -outdir string
        directory under basedir for output (default "build")
  -skipbuild
        don't actually build the patches
  -skipdl
        don't download kobopatch (use this for parallel builds)
  -srcdir string
        directory under basedir for sources (default "src")
```

## test.go
`test.go` tests the patches. It currently only works on Linux/BSD/macOS. To use it, you need Go 1.13+ installed, and then you can run `go run ./scripts/test` from the root of the repository.

```
go run ./scripts/test [options] [versions...]

Options:
  -basedir string
        base dir for files (default ".")
  -kpbin string
        kobopatch binary path (default "kobopatch")
  -srcdir string
        directory under basedir for sources (default "src")
  -testdir string
        directory under basedir for testdata (default "testdata")
```
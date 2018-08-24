package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"

	"github.com/geek1011/kobopatch/patchfile"
	"github.com/geek1011/kobopatch/patchfile/kobopatch"
	"github.com/geek1011/kobopatch/patchlib"
	"github.com/spf13/pflag"
)

func main() {
	help := pflag.BoolP("help", "h", false, "show this help text")
	pflag.Parse()

	if *help || pflag.NArg() != 2 {
		fmt.Fprintf(os.Stderr, "Usage: test-patches [OPTIONS] PATH_TO_PATCH_FILE PATH_TO_BINARY\n\nOptions:\n")
		pflag.PrintDefaults()
		os.Exit(1)
	}

	pf, err := patchfile.ReadFromFile("kobopatch", pflag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not read patch file: %v\n", err)
		os.Exit(1)
	}
	ps := reflect.ValueOf(pf).Interface().(*kobopatch.PatchSet)

	buf, err := ioutil.ReadFile(pflag.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not read binary to patch: %v\n", err)
		os.Exit(1)
	}
	getBuf := func() []byte {
		nbuf := make([]byte, len(buf))
		copy(nbuf, buf)
		return nbuf
	}
	
	sortedNames := []string{}
	for name := range *ps {
		sortedNames = append(sortedNames, name)
	}
	sort.Strings(sortedNames)
	
	errs := map[string]error{}
	for _, name := range sortedNames {
		fmt.Printf(" -  %s", name)
		err := errors.New("no such patch")
		for pname, instructions := range *ps {
			for _, instruction := range instructions {
				if instruction.Enabled != nil {
					*instruction.Enabled = pname == name
					err = nil
					break
				}
			}
		}
		if err != nil {
			fmt.Printf("\r ✕  %s\n", name)
			errs[name] = err
			continue
		}
		out := os.Stdout
		os.Stdout = nil
		if err := ps.ApplyTo(patchlib.NewPatcher(getBuf())); err != nil {
			os.Stdout = out
			fmt.Printf("\r ✕  %s\n", name)
			errs[name] = err
			continue
		}
		os.Stdout = out
		if err := ps.SetEnabled(name, false); err != nil {
			fmt.Printf("\r ✕  %s\n", name)
			errs[name] = err
			continue
		}
		fmt.Printf("\r ✔  %s\n", name)
	}

	if len(errs) >= 1 {
		fmt.Printf("\nErrors:\n")
		for name, err := range errs {
			fmt.Printf(" ✕  %s: %v\n", name, err)
		}
	}

	fmt.Printf("\nSummary: %d total, %d errors, %d success\n", len(*ps), len(errs), len(*ps)-len(errs))
	if len(errs) > 1 {
		fmt.Printf("\nTest failed.\n")
		os.Exit(1)
	}
}

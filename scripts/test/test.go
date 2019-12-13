package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] [versions...]\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}

	basedir := flag.String("basedir", ".", "base dir for files")
	srcdir := flag.String("srcdir", "src", "directory under basedir for sources")
	testdir := flag.String("testdir", "testdata", "directory under basedir for testdata")
	kpbin := flag.String("kpbin", "kobopatch", "kobopatch binary path")
	flag.Parse()

	var err error
	for _, rdir := range []*string{srcdir, testdir} {
		if *rdir, err = filepath.Abs(filepath.Join(*basedir, *rdir)); err != nil {
			fmt.Fprintf(os.Stderr, "Error resolving path to %s\n", *rdir)
			os.Exit(1)
		}
	}

	logSect("scan versions")
	var vers []string
	if flag.NArg() > 0 {
		logItem("provided on command line, skipping scan")
		vers = flag.Args()
	} else {
		if vers, err = versions(*srcdir); err != nil {
			fmt.Fprintf(os.Stderr, "Error scanning versions: %v.\n", err)
			os.Exit(1)
		}
	}

	logSect("test patches")
	for _, version := range vers {
		if err := test(*kpbin, *basedir, *srcdir, *testdir, version); err != nil {
			fmt.Fprintf(os.Stderr, "Tests for %s failed (%v)! Stopping.\n", version, err)
			os.Exit(1)
		}
	}

	os.Exit(0)
}

func test(kpbin, basedir, srcdir, testdir, version string) error {
	logTask("test %s", version)

	logItem("scanning sources")
	fis, err := ioutil.ReadDir(filepath.Join(srcdir, "versions", version))
	if err != nil {
		return logErr(fmt.Errorf("scan files: %v", err))
	}

	sort.Slice(fis, func(i, j int) bool {
		return fis[i].Name() < fis[j].Name()
	})

	var patches [][]string
	for _, fi := range fis {
		if !fi.IsDir() {
			continue
		}

		fn := filepath.Join(srcdir, "versions", version, fi.Name())
		binfn := path.Join("usr", "local", "Kobo", strings.TrimSuffix(fi.Name(), filepath.Ext(fi.Name()))) // always forward slash separated

		sfis, err := ioutil.ReadDir(fn)
		if err != nil {
			return logErr(fmt.Errorf("bin %s: scan sources: %v", binfn, err))
		}

		sort.Slice(sfis, func(i, j int) bool {
			return sfis[i].Name() < sfis[j].Name()
		})

		for _, sfi := range sfis {
			sfn := filepath.Join(fn, sfi.Name())
			patches = append(patches, []string{sfn, binfn})
		}
	}

	logItem("generating config")
	cfg, err := cfg(version, patches)
	if err != nil {
		return logErr(fmt.Errorf("error generating config: %w", err))
	}

	logItem("running kobopatch")
	cmd := exec.Command(kpbin, "-tf", filepath.Join(testdir, version+".tar.xz"), "-")
	cmd.Dir = basedir
	cmd.Stdin = cfg
	if os.Getenv("DRONE") == "true" {
		// Drone doesn't handle CRs properly (used to overwrite the - with a checkmark or an X), so make them into LFs.
		cmd.Stdout = &cr2lf{os.Stdout}
	} else {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func cfg(version string, patches [][]string) (io.Reader, error) {
	cfgbuf := bytes.NewBuffer(nil)
	return cfgbuf, template.Must(template.New("").Parse(`
version: {{.Version}}
in: /dev/null
out: /dev/null
log: /dev/null
patchFormat: kobopatch
patches:{{range .Patches}}
  {{index . 0}}: {{index . 1}}
{{- end}}
`)).Execute(cfgbuf, struct {
		Version string
		Patches [][]string
	}{version, patches})
}

func versions(srcdir string) ([]string, error) {
	fis, err := ioutil.ReadDir(filepath.Join(srcdir, "versions"))
	if err != nil {
		return nil, err
	}
	var vers []string
	for _, fi := range fis {
		if fi.IsDir() {
			vers = append(vers, fi.Name())
		}
	}
	sort.Slice(vers, func(i, j int) bool {
		var a1, a2, a3 int
		if _, err := fmt.Sscanf(vers[i], "%d.%d.%d", &a1, &a2, &a3); err != nil {
			return vers[i] < vers[j]
		}
		var b1, b2, b3 int
		if _, err := fmt.Sscanf(vers[j], "%d.%d.%d", &b1, &b2, &b3); err != nil {
			return vers[i] < vers[j]
		}
		return !(a1 > b1 || (a1 == b1 && (a2 > b2 || (a2 == b2 && (a3 > b3 || a3 == b3)))))
	})
	return vers, nil
}

type cr2lf struct {
	w io.Writer
}

func (c *cr2lf) Write(buf []byte) (n int, err error) {
	for n, b := range buf {
		if b == '\r' {
			buf[n] = '\n'
		}
	}
	return c.w.Write(buf)
}

func logSect(format string, a ...interface{}) { fmt.Printf("### "+format+"\n", a...) }
func logTask(format string, a ...interface{}) { fmt.Printf(">>> "+format+"\n", a...) }
func logItem(format string, a ...interface{}) { fmt.Printf("--- "+format+"\n", a...) }
func logMesg(format string, a ...interface{}) { fmt.Printf("    "+format+"\n", a...) }
func logErr(err error) error {
	if err != nil {
		fmt.Fprintf(os.Stderr, "!!! %v\n", err)
	}
	return err
}

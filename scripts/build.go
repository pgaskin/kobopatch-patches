package main

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] [versions...]\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}

	basedir := flag.String("basedir", ".", "base dir for files")
	srcdir := flag.String("srcdir", "src", "directory under basedir for sources")
	outdir := flag.String("outdir", "build", "directory under basedir for output")
	dldir := flag.String("dldir", "dl", "directory under basedir for kobopatch download")
	kprepo := flag.String("kprepo", "geek1011/kobopatch", "github repo for kobopatch")
	kpver := flag.String("kpver", "v0.14.0", "kobopatch version")
	kpbin := flag.String("kpbin", strings.Join([]string{
		"kobopatch-darwin-64bit", "cssextract-darwin-64bit",
		"kobopatch-linux-32bit", "cssextract-linux-32bit",
		"kobopatch-linux-64bit", "cssextract-linux-64bit",
		"koboptch-windows.exe", "cssextract-windows.exe",
	}, ","), "kobopatch binaries to download")
	skipbuild := flag.Bool("skipbuild", false, "don't actually build the patches")
	skipdl := flag.Bool("skipdl", false, "don't download kobopatch (use this for parallell builds)")
	flag.Parse()

	var err error
	for _, rdir := range []*string{srcdir, outdir, dldir} {
		if *rdir, err = filepath.Abs(filepath.Join(*basedir, *rdir)); err != nil {
			fmt.Fprintf(os.Stderr, "Error resolving path to %s\n", *rdir)
		}
	}

	if !*skipdl {
		fmt.Printf("### download tools (%s)\n", *kpver)
		if errs := dl("dl", *kprepo, *kpver, strings.Split(*kpbin, ",")...); len(errs) != 0 {
			fmt.Fprintf(os.Stderr, "Error downloading kobopatch.\n")
			os.Exit(1)
		}
	}

	if *skipbuild {
		os.Exit(0)
	}

	fmt.Printf("### scan versions\n")
	var vers []string
	if flag.NArg() > 0 {
		fmt.Printf("--- provided on command line, skipping scan\n")
		vers = flag.Args()
	} else {
		if vers, err = versions(*srcdir); err != nil {
			fmt.Fprintf(os.Stderr, "Error scanning versions: %v.\n", err)
		}
	}

	fmt.Printf("### build patch zips\n")
	for _, version := range vers {
		if err := build(*srcdir, *dldir, *outdir, *kpver, version, strings.Split(*kpbin, ",")); err != nil {
			fmt.Fprintf(os.Stderr, "Error building %s: versions: %v.\n", version, err)
			os.Exit(1)
		}
	}

	os.Exit(0)
}

func build(srcdir, dldir, outdir, kpver, version string, kpbin []string) error {
	fmt.Printf(">>> build %s\n", version)
	kpf := filepath.Join(outdir, "kobopatch_"+version+".zip")

	fmt.Printf("--- create output file\n")
	os.MkdirAll(outdir, 0755)
	f, err := ioutil.TempFile(outdir, "tmp_*.zip")
	if err != nil {
		fmt.Printf("!!! create output file: %v\n", err)
		return fmt.Errorf("create output file: %v", err)
	}
	fmt.Printf("    created %s\n", f.Name())
	defer f.Close() // this will return an error (which is ignored) if called after the close at the end, but this is just for cleanup on an actual error
	defer os.Remove(f.Name())

	ss := sha1.New()
	zw := zip.NewWriter(io.MultiWriter(f, ss))

	// change compression for best time/size ratio
	zw.RegisterCompressor(zip.Deflate, func(w io.Writer) (io.WriteCloser, error) {
		// inefficient (writers aren't pooled)
		return flate.NewWriter(w, 3)
	})

	fmt.Printf("--- add template\n")
	if err := add(zw, ".", filepath.Join(srcdir, "template"), ".", map[string]string{
		"{{version}}": version,
	}); err != nil {
		fmt.Printf("!!! add template: %v\n", err)
		return fmt.Errorf("add template: %v", err)
	}

	fmt.Printf("--- add kobopatch\n")
	for _, bin := range kpbin {
		fmt.Printf("    add %s from %s\n", filepath.Join("bin", bin), filepath.Join(dldir, kpver, bin))
		f, err := os.Open(filepath.Join(dldir, kpver, bin))
		if err != nil {
			fmt.Printf("!!! open %s: %v\n", bin, err)
			return fmt.Errorf("open %s: %v", bin, err)
		}

		fi, err := f.Stat()
		if err != nil {
			f.Close()
			fmt.Printf("!!! stat %s: %v\n", bin, err)
			return fmt.Errorf("stat %s: %v", bin, err)
		}

		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			f.Close()
			fmt.Printf("!!! add %s: %v\n", bin, err)
			return fmt.Errorf("add %s: %v", bin, err)
		}
		fh.Name = filepath.Join("bin", bin)
		fh.Method = zip.Deflate

		w, err := zw.CreateHeader(fh)
		if err != nil {
			f.Close()
			fmt.Printf("!!! add %s: %v\n", bin, err)
			return fmt.Errorf("add %s: %v", bin, err)
		}

		if _, err := io.CopyN(w, f, fi.Size()); err != nil {
			f.Close()
			fmt.Printf("!!! write %s: %v\n", bin, err)
			return fmt.Errorf("write %s: %v", bin, err)
		}

		f.Close()
	}

	fmt.Printf("--- scanning for files to generate\n")
	fis, err := ioutil.ReadDir(filepath.Join(srcdir, "versions", version))
	if err != nil {
		fmt.Printf("!!! scan versions: %v\n", err)
		return fmt.Errorf("scan versions: %v", err)
	}

	for _, fi := range fis {
		if !fi.IsDir() {
			continue
		}
		fn := filepath.Join(srcdir, "versions", version, fi.Name())
		gfn := filepath.Join("src", fi.Name())
		fmt.Printf("--- generating %s from %s\n", gfn, fn)

		fmt.Printf("    scanning for source files\n")
		sfis, err := ioutil.ReadDir(fn)
		if err != nil {
			fmt.Printf("!!! scan sources: %v\n", err)
			return fmt.Errorf("generate %s: scan sources: %v", gfn, err)
		}

		fmt.Printf("    sorting source files\n")
		sort.Slice(sfis, func(i, j int) bool {
			return sfis[i].Name() < sfis[j].Name()
		})

		modtime := fi.ModTime() // for reproducibility
		bufw := bytes.NewBuffer(nil)
		for _, sfi := range sfis {
			sfn := filepath.Join(fn, sfi.Name())
			fmt.Printf("    merging %s from %s\n", sfi.Name(), sfn)
			sbuf, err := ioutil.ReadFile(sfn)
			if err != nil {
				fmt.Printf("!!! read source file: %v\n", err)
				return fmt.Errorf("generate %s: read source file %s: %v", gfn, sfi.Name(), err)
			}
			bufw.Write(sbuf)
			bufw.WriteRune('\n')
			if !bytes.HasSuffix(sbuf, []byte{'\n'}) {
				bufw.WriteRune('\n') // each part should be separated with a blank line
			}
			if sfi.ModTime().After(modtime) {
				modtime = sfi.ModTime()
			}
		}

		fmt.Printf("    converting unix line breaks to dos\n")
		buf := bytes.ReplaceAll(bufw.Bytes(), []byte{'\r', '\n'}, []byte{'\n'}) // convert dos to unix in case it is already dos
		buf = bytes.ReplaceAll(buf, []byte{'\n'}, []byte{'\r', '\n'})           // and then back to dos again

		fmt.Printf("    adding to zip (mod: %s)\n", modtime)
		zh := &zip.FileHeader{
			Name:               gfn,
			UncompressedSize:   uint32(len(buf)),
			UncompressedSize64: uint64(len(buf)),
		}
		zh.SetModTime(modtime)
		zh.SetMode(0644)

		if w, err := zw.CreateHeader(zh); err != nil {
			fmt.Printf("!!! add to zip: %v\n", err)
			return fmt.Errorf("generate %s: add to zip: %v", gfn, err)
		} else if _, err := io.CopyN(w, bytes.NewReader(buf), int64(len(buf))); err != nil {
			fmt.Printf("!!! add to zip: %v\n", err)
			return fmt.Errorf("generate %s: add to zip: %v", gfn, err)
		}

	}

	if err := zw.Close(); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	fmt.Printf("--- move output to target path %s\n", kpf)
	if err := os.Rename(f.Name(), kpf); err != nil {
		fmt.Printf("!!! move output: %v\n", err)
		return fmt.Errorf("move output %s: %v", f.Name(), err)
	}

	fmt.Printf("--- calculating sha1: %x\n", ss.Sum(nil))
	return nil
}

func dl(dldir, repo, tag string, files ...string) []error {
	var errs []error
	for _, file := range files {
		fmt.Printf(">>> download %s@%s/%s to %s/%s/%s\n", repo, tag, file, dldir, tag, file)
		url := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", repo, tag, file)
		fp := filepath.Join(dldir, tag, file)

		if fi, err := os.Stat(fp); err == nil && fi.Size() > 0 {
			fmt.Printf("--- already downloaded, skipping\n")
			continue
		}

		resp, err := http.Get(url)
		if err == nil && resp.StatusCode != 200 {
			err = fmt.Errorf("response status %s", resp.Status)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "!!! error: get: %v\n", err)
			errs = append(errs, fmt.Errorf("get %s@%s/%s: %w", repo, tag, file, err))
			continue
		}
		defer resp.Body.Close()

		if err := os.MkdirAll(filepath.Dir(fp), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "!!! error: mkdir: %v\n", err)
			errs = append(errs, fmt.Errorf("get %s@%s/%s: mkdir: %w", repo, tag, file, err))
			continue
		}

		f, err := ioutil.TempFile(filepath.Join(dldir, tag), "tmp_*")
		if err != nil {
			fmt.Fprintf(os.Stderr, "!!! error: create: %v\n", err)
			errs = append(errs, fmt.Errorf("get %s@%s/%s: create: %w", repo, tag, file, err))
			continue
		}
		defer f.Close()
		defer os.Remove(f.Name())

		if _, err := io.Copy(f, resp.Body); err != nil {
			fmt.Fprintf(os.Stderr, "\n!!! error: write: %v\n", err)
			errs = append(errs, fmt.Errorf("get %s@%s/%s: write: %w", repo, tag, file, err))
			continue
		}

		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "\n!!! error: close: %v\n", err)
			errs = append(errs, fmt.Errorf("get %s@%s/%s: close: %w", repo, tag, file, err))
			continue
		}

		if err := os.Rename(f.Name(), fp); err != nil {
			fmt.Fprintf(os.Stderr, "\n!!! error: rename: %v\n", err)
			errs = append(errs, fmt.Errorf("get %s@%s/%s: rename: %w", repo, tag, file, err))
			continue
		}

		if err := os.Chmod(fp, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "\n!!! error: chmod: %v\n", err)
			errs = append(errs, fmt.Errorf("get %s@%s/%s: chmod: %w", repo, tag, file, err))
			continue
		}

		if t, err := time.Parse(time.RFC1123, resp.Header.Get("Last-Modified")); err == nil {
			os.Chtimes(fp, t, t)
		}
	}
	return errs
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

func add(zw *zip.Writer, zipBase, base, dir string, replace map[string]string) error {
	fis, err := ioutil.ReadDir(filepath.Join(base, dir))
	if err != nil {
		return err
	}
	for _, fi := range fis {
		zh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return fmt.Errorf("add %s: %w", fi.Name(), err)
		}

		fp := filepath.Join(base, dir, fi.Name())
		zh.Name = filepath.ToSlash(filepath.Clean(filepath.Join(zipBase, dir, zh.Name)))
		if fi.IsDir() {
			zh.Name += "/"
			zh.Method = zip.Store
		} else {

			zh.Method = zip.Deflate
		}
		fmt.Printf("    add %s from %s\n", zh.Name, fp)

		var w io.Writer
		if replace == nil || len(replace) == 0 {
			w, err = zw.CreateHeader(zh)
			if err != nil {
				return fmt.Errorf("add %s: %w", fi.Name(), err)
			}
		}

		if fi.IsDir() {
			if err := add(zw, zipBase, base, filepath.Join(dir, fi.Name()), replace); err != nil {
				return fmt.Errorf("add %s: %v", fi.Name(), err)
			}
		} else if replace == nil || len(replace) == 0 {
			f, err := os.Open(fp)
			if err != nil {
				return fmt.Errorf("add %s: %w", fi.Name(), err)
			}
			if _, err := io.Copy(w, f); err != nil {
				f.Close()
				return fmt.Errorf("add %s: %w", fi.Name(), err)
			}
			f.Close()
		} else {
			buf, err := ioutil.ReadFile(fp)
			if err != nil {
				return fmt.Errorf("add %s: %w", fi.Name(), err)
			}
			for find, replace := range replace { // inefficient
				buf = bytes.ReplaceAll(buf, []byte(find), []byte(replace))
			}
			zh.UncompressedSize64 = uint64(len(buf))
			zh.UncompressedSize = uint32(zh.UncompressedSize64)
			w, err := zw.CreateHeader(zh)
			if err != nil {
				return fmt.Errorf("add %s: %w", fi.Name(), err)
			}
			if _, err := io.CopyN(w, bytes.NewReader(buf), int64(len(buf))); err != nil {
				return fmt.Errorf("add %s: %w", fi.Name(), err)
			}
		}
	}
	return nil
}

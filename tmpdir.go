package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	prefixes = []string{
		"go-build",
		"go-tool-dist-",
		"go-code-check",
		"go-link-",
		"go-sqlite3-test-",
		"cgo-gcc-input-",
		"check-log-test",
		"check-event-log-test",
		"check-windows-eventlog-test",
		"gom",
		"tomcat.",
		"tmp-",
		"zglob",
	}
)

func init() {
	register("TmpDir", cleanTmpDir)
}

func dirsize(name string) uint64 {
	dir, err := os.Open(name)
	if err != nil {
		return 0
	}
	fis, err := dir.Readdir(-1)
	if err != nil {
		return 0
	}
	dir.Close()
	var size uint64
	for _, fi := range fis {
		if fi.IsDir() {
			size += dirsize(filepath.Join(name, fi.Name()))
		} else {
			size += uint64(fi.Size())
		}
	}
	return size
}

func cleanTmpDir(dryrun, verbose bool) (string, error) {
	name := os.TempDir()
	dir, err := os.Open(name)
	defer dir.Close()

	fis, err := dir.Readdir(-1)
	if err != nil {
		return "", err
	}
	var size uint64
	for _, fi := range fis {
		for _, prefix := range prefixes {
			if strings.HasPrefix(fi.Name(), prefix) {
				p := filepath.Join(name, fi.Name())
				if verbose {
					log.Println("TmpDir:", p)
				}
				if fi.IsDir() {
					size += dirsize(p)
					if !dryrun {
						err = os.RemoveAll(p)
						if err != nil {
							return "", err
						}
					}
				} else {
					size += uint64(fi.Size())
					if !dryrun {
						err = os.Remove(p)
						if err != nil {
							return "", err
						}
					}
				}
			}
		}
	}
	result := "%.1f MB removed"
	if dryrun {
		result = "%.1f MB removable"
	}
	return fmt.Sprintf(result, float64(size)/1024/1024), nil
}

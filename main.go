package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"
)

var (
	version = "dev"
	commit  = "none"
)

// PrintOptions 表示のオプション
type PrintOptions struct {
	printSize bool
	printMd5  bool
}

func main() {

	if len(commit) > 7 {
		commit = commit[:7]
	}

	var printSize bool
	var printMd5 bool
	var help bool

	flag.BoolVarP(&printSize, "size", "s", false, "Print file size")
	flag.BoolVarP(&printMd5, "md5", "m", false, "Print md5 hash")
	flag.BoolVarP(&help, "help", "h", false, "Help")
	flag.Parse()
	flag.CommandLine.SortFlags = false
	flag.Usage = func() {
		fmt.Printf("filist v%s (%s)\n", version, commit)
		fmt.Fprintf(os.Stderr, "Usage of %s:\n%s [options] directory ...\noptions\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	if help {
		flag.Usage()
		os.Exit(0)
	}

	dirs := flag.Args()

	if len(dirs) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	err := print(dirs, PrintOptions{printSize: printSize, printMd5: printMd5})
	if err != nil {
		fmt.Println("\nError: ", err)
		os.Exit(1)
	}
}

func print(dirs []string, options PrintOptions) error {

	for _, dir := range dirs {
		err := printDir(dir, options)
		if err != nil {
			return err
		}
	}

	return nil
}

func printDir(dir string, options PrintOptions) error {

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		printFileInfo(path, info, options)

		return nil
	})

	return err
}

func printFileInfo(filePath string, info os.FileInfo, options PrintOptions) error {

	fmt.Print(filePath)

	if options.printSize {
		fmt.Printf("\t%v", info.Size())
	}

	if options.printMd5 {
		md5, err := calcMd5(filePath)
		if err != nil {
			return err
		}

		fmt.Printf("\t%v", md5)
	}

	fmt.Println()

	return nil
}

func calcMd5(filePath string) (string, error) {

	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

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

// Option 表示のオプション
type Option struct {
	optionalColumns []func(string, os.FileInfo) (string, error)
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

	var optionalColumns []func(string, os.FileInfo) (string, error)

	// オプションは指定順に表示したいので
	flag.Visit(func(f *flag.Flag) {
		switch f.Shorthand {
		case "s":
			optionalColumns = append(optionalColumns, calcSize)
		case "m":
			optionalColumns = append(optionalColumns, calcMd5)
		}
	})

	if help {
		flag.Usage()
		os.Exit(0)
	}

	dirs := flag.Args()

	if len(dirs) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	err := print(dirs, Option{optionalColumns: optionalColumns})
	if err != nil {
		fmt.Println("\nError: ", err)
		os.Exit(1)
	}
}

func print(dirs []string, option Option) error {

	for _, dir := range dirs {
		err := printDir(dir, option)
		if err != nil {
			return err
		}
	}

	return nil
}

func printDir(dir string, option Option) error {

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		printFileInfo(path, info, option)

		return nil
	})

	return err
}

func printFileInfo(filePath string, info os.FileInfo, option Option) error {

	fmt.Print(filePath)

	// オプション分を出力
	for _, column := range option.optionalColumns {
		value, err := column(filePath, info)
		if err != nil {
			return err
		}
		fmt.Printf("\t%s", value)
	}

	fmt.Println()

	return nil
}

func calcMd5(filePath string, info os.FileInfo) (string, error) {

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

func calcSize(filePath string, info os.FileInfo) (string, error) {

	return fmt.Sprint(info.Size()), nil
}

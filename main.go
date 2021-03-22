package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
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
	columns []func(string, string, os.FileInfo) (string, error)
}

func main() {

	if len(commit) > 7 {
		commit = commit[:7]
	}

	var help bool
	var printRelPath bool
	var printAbsPath bool

	flag.BoolVarP(&printRelPath, "rel", "r", false, "Print relative path (If neither 'rel' nor 'abs' is specified, 'rel' will be printed first column.)")
	flag.BoolVarP(&printAbsPath, "abs", "a", false, "Print absolute path")
	flag.BoolP("size", "s", false, "Print file size")
	flag.BoolP("mtime", "m", false, "Print modification time")
	flag.BoolP("md5", "M", false, "Print MD5 hash")
	flag.BoolP("sha1", "S", false, "Print SHA-1 hash")
	flag.BoolP("sha256", "", false, "Print SHA-256 hash")
	flag.BoolVarP(&help, "help", "h", false, "Help")
	flag.Parse()
	flag.CommandLine.SortFlags = false
	flag.Usage = func() {
		fmt.Printf("filist v%s (%s)\n", version, commit)
		fmt.Fprintf(os.Stderr, "Usage: %s [options] directory ...\noptions\n", os.Args[0])
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

	var columns []func(string, string, os.FileInfo) (string, error)

	if !printRelPath && !printAbsPath {
		// relとabsどちらも指定されていなかった場合、先頭にrelを表示
		columns = append(columns, getRelPath)
	}

	// オプションは指定順に表示したいので
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "rel":
			columns = append(columns, getRelPath)
		case "abs":
			columns = append(columns, getAbsPath)
		case "size":
			columns = append(columns, getSize)
		case "mtime":
			columns = append(columns, getMtime)
		case "md5":
			columns = append(columns, calcMd5)
		case "sha1":
			columns = append(columns, calcSha1)
		case "sha256":
			columns = append(columns, calcSha256)
		}
	})

	option := Option{
		columns: columns,
	}

	err := print(dirs, option)

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

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	err = filepath.Walk(absDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		printFileInfo(absDir, path, info, option)

		return nil
	})

	return err
}

func printFileInfo(baseDir string, filePath string, info os.FileInfo, option Option) error {

	for i, column := range option.columns {
		if i > 0 {
			fmt.Printf("\t")
		}
		value, err := column(baseDir, filePath, info)
		if err != nil {
			return err
		}
		fmt.Printf("%s", value)
	}

	fmt.Println()

	return nil
}

func getRelPath(baseDir string, filePath string, info os.FileInfo) (string, error) {

	return filepath.Rel(baseDir, filePath)
}

func getAbsPath(baseDir string, filePath string, info os.FileInfo) (string, error) {

	return filePath, nil
}

func getSize(baseDir string, filePath string, info os.FileInfo) (string, error) {

	return fmt.Sprint(info.Size()), nil
}

func getMtime(baseDir string, filePath string, info os.FileInfo) (string, error) {

	return info.ModTime().Format("2006-01-02T15:04:05.000000-07:00"), nil
}

func calcMd5(baseDir string, filePath string, info os.FileInfo) (string, error) {

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

func calcSha1(baseDir string, filePath string, info os.FileInfo) (string, error) {

	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func calcSha256(baseDir string, filePath string, info os.FileInfo) (string, error) {

	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

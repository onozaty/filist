package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"
)

var (
	Version = "dev"
	Commit  = "dev"
)

// Option 表示のオプション
type Option struct {
	includeDirectories bool
	excludeFiles       bool
	columns            []func(string, string, os.FileInfo) (string, error)
}

func main() {

	var help bool
	var printRelPath bool
	var printAbsPath bool
	var includeDirectories bool
	var excludeFiles bool

	flag.BoolVarP(&printRelPath, "rel", "r", false, "Print relative path (If neither 'rel' nor 'abs' is specified, 'rel' will be printed first column.)")
	flag.BoolVarP(&printAbsPath, "abs", "a", false, "Print absolute path")
	flag.BoolP("size", "s", false, "Print file size")
	flag.BoolP("mtime", "m", false, "Print modification time")
	flag.BoolP("md5", "M", false, "Print MD5 hash")
	flag.BoolP("sha1", "S", false, "Print SHA-1 hash")
	flag.BoolP("sha256", "", false, "Print SHA-256 hash")
	flag.BoolVarP(&includeDirectories, "include-dir", "", false, "Include directories")
	flag.BoolVarP(&excludeFiles, "exclude-file", "", false, "Exclude files")
	flag.BoolVarP(&help, "help", "h", false, "Help")
	flag.Parse()
	flag.CommandLine.SortFlags = false
	flag.Usage = func() {
		fmt.Printf("filist v%s (%s)\n\n", Version, Commit)
		fmt.Fprint(os.Stderr, "Usage: filist [flags] directory ...\n\nFlags\n")
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
		columns:            columns,
		includeDirectories: includeDirectories,
		excludeFiles:       excludeFiles,
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

		if absDir == path {
			return nil
		}

		if info.IsDir() {
			if option.includeDirectories {
				return printFileInfo(absDir, path, info, option)
			}
		} else {
			if !option.excludeFiles {
				return printFileInfo(absDir, path, info, option)
			}
		}

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

	relPath, err := filepath.Rel(baseDir, filePath)
	if err != nil {
		return "", err
	}

	if info.IsDir() {
		relPath = relPath + string(filepath.Separator)
	}

	return relPath, nil
}

func getAbsPath(baseDir string, filePath string, info os.FileInfo) (string, error) {

	if info.IsDir() {
		filePath = filePath + string(filepath.Separator)
	}

	return filePath, nil
}

func getSize(baseDir string, filePath string, info os.FileInfo) (string, error) {

	if info.IsDir() {
		return "", nil
	}

	return fmt.Sprint(info.Size()), nil
}

func getMtime(baseDir string, filePath string, info os.FileInfo) (string, error) {

	if info.IsDir() {
		return "", nil
	}

	return info.ModTime().Format("2006-01-02T15:04:05.000000-07:00"), nil
}

func calcMd5(baseDir string, filePath string, info os.FileInfo) (string, error) {

	if info.IsDir() {
		return "", nil
	}

	return calcHash(filePath, md5.New())
}

func calcSha1(baseDir string, filePath string, info os.FileInfo) (string, error) {

	if info.IsDir() {
		return "", nil
	}

	return calcHash(filePath, sha1.New())
}

func calcSha256(baseDir string, filePath string, info os.FileInfo) (string, error) {

	if info.IsDir() {
		return "", nil
	}

	return calcHash(filePath, sha256.New())
}

func calcHash(filePath string, hash hash.Hash) (string, error) {

	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

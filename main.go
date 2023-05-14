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

const (
	OK int = 0
	NG int = 1
)

func main() {
	exitCode := run(os.Args[1:], os.Stdout)
	os.Exit(exitCode)
}

func run(arguments []string, output io.Writer) int {

	var help bool
	var printRelPath bool
	var printAbsPath bool
	var includeDirectories bool
	var excludeFiles bool

	flagSet := flag.NewFlagSet("filist", flag.ContinueOnError)

	flagSet.BoolVarP(&printRelPath, "rel", "r", false, "Print relative path (If neither 'rel' nor 'abs' is specified, 'rel' will be printed first column.)")
	flagSet.BoolVarP(&printAbsPath, "abs", "a", false, "Print absolute path")
	flagSet.BoolP("size", "s", false, "Print file size")
	flagSet.BoolP("mtime", "m", false, "Print modification time")
	flagSet.BoolP("md5", "M", false, "Print MD5 hash")
	flagSet.BoolP("sha1", "S", false, "Print SHA-1 hash")
	flagSet.BoolP("sha256", "", false, "Print SHA-256 hash")
	flagSet.BoolVarP(&includeDirectories, "include-dir", "", false, "Include directories")
	flagSet.BoolVarP(&excludeFiles, "exclude-file", "", false, "Exclude files")
	flagSet.BoolVarP(&help, "help", "h", false, "Help")

	flagSet.SortFlags = false
	flagSet.Usage = func() {
		fmt.Fprintf(output, "filist v%s (%s)\n\n", Version, Commit)
		fmt.Fprint(output, "Usage: filist [flags] directory ...\n\nFlags\n")
		flagSet.PrintDefaults()
		fmt.Fprintln(output)
	}
	flagSet.SetOutput(output)

	if err := flagSet.Parse(arguments); err != nil {
		flagSet.Usage()
		fmt.Fprintln(output, err)
		return NG
	}

	if help {
		flagSet.Usage()
		return OK
	}

	dirs := flagSet.Args()

	if len(dirs) == 0 {
		flagSet.Usage()
		return NG
	}

	var columns []func(string, string, os.FileInfo) (string, error)

	if !printRelPath && !printAbsPath {
		// relとabsどちらも指定されていなかった場合、先頭にrelを表示
		columns = append(columns, getRelPath)
	}

	// オプションは指定順に表示したいので
	flagSet.Visit(func(f *flag.Flag) {
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
		fmt.Fprintln(output, err)
		return NG
	}

	return OK
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

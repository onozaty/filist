package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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
	level              int
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

func run(arguments []string, out io.Writer) int {

	var help bool
	var printRelPath bool
	var printAbsPath bool
	var includeDirectories bool
	var excludeFiles bool
	var level int

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
	flagSet.IntVarP(&level, "level", "l", 0, "Number of directory level (Default is unlimited)")
	flagSet.BoolVarP(&help, "help", "h", false, "Help")

	flagSet.SortFlags = false
	flagSet.Usage = func() {
		fmt.Fprintf(out, "filist v%s (%s)\n\n", Version, Commit)
		fmt.Fprint(out, "Usage: filist [flags] directory ...\n\nFlags\n")
		flagSet.PrintDefaults()
	}
	flagSet.SetOutput(out)

	if err := flagSet.Parse(arguments); err != nil {
		flagSet.Usage()
		fmt.Fprintf(out, "Error: %v", err)
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
		level:              level,
	}

	err := print(out, dirs, option)

	if err != nil {
		fmt.Fprintf(out, "Error: %v", err)
		return NG
	}

	return OK
}

func print(out io.Writer, dirs []string, option Option) error {

	for _, dir := range dirs {
		err := printDir(out, dir, option)
		if err != nil {
			return err
		}
	}

	return nil
}

func printDir(out io.Writer, dir string, option Option) error {

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	err = filepath.WalkDir(absDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if absDir == path {
			return nil
		}

		if d.IsDir() {
			if option.includeDirectories {

				info, err := d.Info()
				if err != nil {
					return err
				}

				if err := printFileInfo(out, absDir, path, info, option); err != nil {
					return err
				}
			}

			depth, err := getDepth(absDir, path)
			if err != nil {
				return err
			}

			if option.level != 0 && depth >= int(option.level) {
				// 指定レベル以上になったら、そのディレクトリ配下は見ない
				return filepath.SkipDir
			}

		} else {
			if !option.excludeFiles {

				info, err := d.Info()
				if err != nil {
					return err
				}

				return printFileInfo(out, absDir, path, info, option)
			}
		}

		return nil
	})

	return err
}

func getDepth(basePath string, path string) (int, error) {

	relPath, err := filepath.Rel(basePath, path)
	if err != nil {
		return 0, err
	}

	return len(strings.Split(relPath, string(filepath.Separator))), nil
}

func printFileInfo(out io.Writer, baseDir string, filePath string, info os.FileInfo, option Option) error {

	for i, column := range option.columns {
		if i > 0 {
			fmt.Fprintf(out, "\t")
		}
		value, err := column(baseDir, filePath, info)
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "%s", value)
	}

	fmt.Fprintln(out)

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

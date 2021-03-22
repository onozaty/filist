package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestRelPath(t *testing.T) {

	baseDir, tempFile, info, err := craeteTempFile([]byte("ABCDEFG"))

	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(tempFile)

	relPath, err := getRelPath(baseDir, tempFile, info)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if relPath != filepath.Base(tempFile) {
		t.Fatal("failed test\n", relPath)
	}
}

func TestAbsPath(t *testing.T) {

	baseDir, tempFile, info, err := craeteTempFile([]byte("ABCDEFG"))

	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(tempFile)

	relPath, err := getAbsPath(baseDir, tempFile, info)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if relPath != tempFile {
		t.Fatal("failed test\n", relPath)
	}
}

func TestGetSize(t *testing.T) {

	baseDir, tempFile, info, err := craeteTempFile([]byte("ABCDEFG"))

	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(tempFile)

	size, err := getSize(baseDir, tempFile, info)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if size != "7" {
		t.Fatal("failed test\n", size)
	}
}

func TestGetMtime(t *testing.T) {

	baseDir, tempFile, info, err := craeteTempFile([]byte("ABCDEFG"))

	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(tempFile)

	mtime, err := getMtime(baseDir, tempFile, info)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if len(mtime) != len("2021-03-14T23:24:00.750627+09:00") {
		t.Fatal("failed test\n", mtime)
	}
}

func TestCalcMd5(t *testing.T) {

	baseDir, tempFile, info, err := craeteTempFile([]byte("ABCDEFG"))

	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(tempFile)

	md5, err := calcMd5(baseDir, tempFile, info)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if md5 != "bb747b3df3130fe1ca4afa93fb7d97c9" {
		t.Fatal("failed test\n", md5)
	}
}

func TestCalcMd5_empty(t *testing.T) {

	baseDir, tempFile, info, err := craeteTempFile([]byte{})

	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(tempFile)

	md5, err := calcMd5(baseDir, tempFile, info)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if md5 != "d41d8cd98f00b204e9800998ecf8427e" {
		t.Fatal("failed test\n", md5)
	}
}

func TestCalcSha1(t *testing.T) {

	baseDir, tempFile, info, err := craeteTempFile([]byte("ABCDEFG"))

	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(tempFile)

	md5, err := calcSha1(baseDir, tempFile, info)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if md5 != "93be4612c41d23af1891dac5fd0d535736ffc4e3" {
		t.Fatal("failed test\n", md5)
	}
}

func TestCalcSha1_empty(t *testing.T) {

	baseDir, tempFile, info, err := craeteTempFile([]byte{})

	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(tempFile)

	md5, err := calcSha1(baseDir, tempFile, info)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if md5 != "da39a3ee5e6b4b0d3255bfef95601890afd80709" {
		t.Fatal("failed test\n", md5)
	}
}

func TestCalcSha256(t *testing.T) {

	baseDir, tempFile, info, err := craeteTempFile([]byte("ABCDEFG"))

	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(tempFile)

	md5, err := calcSha256(baseDir, tempFile, info)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if md5 != "e9a92a2ed0d53732ac13b031a27b071814231c8633c9f41844ccba884d482b16" {
		t.Fatal("failed test\n", md5)
	}
}

func TestCalcSha256_empty(t *testing.T) {

	baseDir, tempFile, info, err := craeteTempFile([]byte{})

	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(tempFile)

	md5, err := calcSha256(baseDir, tempFile, info)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if md5 != "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
		t.Fatal("failed test\n", md5)
	}
}

func TestPrintDir(t *testing.T) {

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout := os.Stdout
	os.Stdout = w

	tempDir, err := os.MkdirTemp("", "filist")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.RemoveAll(tempDir)

	file1, err := os.Create(filepath.Join(tempDir, "file1"))
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	file1.Write([]byte("a"))

	// サブディレクトリにもファイル配置
	sub1 := filepath.Join(tempDir, "sub1")
	err = os.Mkdir(sub1, 0777)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	file2, err := os.Create(filepath.Join(sub1, "file2"))
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	file2.Write([]byte("abc"))

	err = printDir(tempDir, Option{columns: []func(string, string, os.FileInfo) (string, error){getRelPath}})
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	os.Stdout = stdout
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)

	if buf.String() != "file1\nsub1\\file2\n" {
		t.Fatal("failed test\n", buf.String())
	}
}

func TestPrintDir_option(t *testing.T) {

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout := os.Stdout
	os.Stdout = w

	tempDir, err := os.MkdirTemp("", "filist")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.RemoveAll(tempDir)

	file1, err := os.Create(filepath.Join(tempDir, "file1"))
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	file1.Write([]byte("a"))

	err = printDir(tempDir, Option{columns: []func(string, string, os.FileInfo) (string, error){getRelPath, getSize, calcMd5}})
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	os.Stdout = stdout
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)

	if buf.String() != "file1\t1\t0cc175b9c0f1b6a831c399e269772661\n" {
		t.Fatal("failed test\n", buf.String())
	}
}

func craeteTempFile(contents []byte) (string, string, os.FileInfo, error) {

	tempFile, err := os.CreateTemp("", "filist")

	if err != nil {
		return os.TempDir(), "", nil, err
	}
	defer tempFile.Close()

	tempFile.Write(contents)
	info, err := tempFile.Stat()
	if err != nil {
		return os.TempDir(), "", nil, err
	}

	return os.TempDir(), tempFile.Name(), info, nil
}

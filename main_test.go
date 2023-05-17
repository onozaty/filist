package main

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {

	loc, _ := time.LoadLocation("UTC")
	time.Local = loc

	code := m.Run()
	os.Exit(code)
}

func TestRun(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	setupFiles(t, temp)

	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			temp,
		},
		out,
	)

	// ASSERT
	require.Equal(t, OK, exitCode)

	expected := allLines(
		line("1.txt"),
		line(filepath.Join("a", "a.txt")),
		line(filepath.Join("a", "b.txt")),
		line(filepath.Join("a", "xxx", "x.txt")),
		line(filepath.Join("x", "y", "z", "テスト.txt")),
	)
	assert.Equal(t, expected, out.String())
}

func TestRun_All(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	setupFiles(t, temp)

	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			temp,
			"-r",
			"-a",
			"-s",
			"-m",
			"-M",
			"-S",
			"--sha256",
		},
		out,
	)

	// ASSERT
	require.Equal(t, OK, exitCode)

	expected := allLines(
		line(
			"1.txt",
			filepath.Join(temp, "1.txt"),
			"0",
			"2020-01-01T00:00:00.000000+00:00",
			"d41d8cd98f00b204e9800998ecf8427e",
			"da39a3ee5e6b4b0d3255bfef95601890afd80709",
			"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		),
		line(
			filepath.Join("a", "a.txt"),
			filepath.Join(temp, "a", "a.txt"),
			"1",
			"2020-12-21T11:12:21.000000+00:00",
			"9dd4e461268c8034f5c8564e155c67a6",
			"11f6ad8ec52a2984abaafd7c3b516503785c2072",
			"2d711642b726b04401627ca9fbac32f5c8530fb1903cc4db02258717921a4881",
		),
		line(
			filepath.Join("a", "b.txt"),
			filepath.Join(temp, "a", "b.txt"),
			"10",
			"2020-12-20T00:00:00.000000+00:00",
			"336311a016184326ddbdd61edd4eeb52",
			"ff9ee043d85595eb255c05dfe32ece02a53efbb2",
			"fc11d6f28e59d3cc33c0b14ceb644bf0902ebd63d61218dffe9e7dac7c254542",
		),
		line(
			filepath.Join("a", "xxx", "x.txt"),
			filepath.Join(temp, "a", "xxx", "x.txt"),
			"20",
			"2019-01-01T12:34:56.000000+00:00",
			"baf1da0e2b9065ab5edd36ca00ed1826",
			"d02e53411e8cb4cd709778f173f7bc9a3455f8ed",
			"d4fc1db665446507dc51b0c9392dd9649291581bfe1b48e241b2b08032b3b647",
		),
		line(
			filepath.Join("x", "y", "z", "テスト.txt"),
			filepath.Join(temp, "x", "y", "z", "テスト.txt"),
			"100",
			"2021-03-28T00:12:34.000000+00:00",
			"aed563ecafb4bcc5654c597a421547b2",
			"50e483690ec481f4af7f6fb524b2b99eb1716565",
			"09ecb6ebc8bcefc733f6f2ec44f791abeed6a99edf0cc31519637898aebd52d8",
		),
	)
	assert.Equal(t, expected, out.String())
}

func TestRun_All_IncludeDir(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	setupFiles(t, temp)

	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			temp,
			"-a",
			"-s",
			"-m",
			"-M",
			"-S",
			"--sha256",
			"--include-dir",
			"-r",
		},
		out,
	)

	// ASSERT
	require.Equal(t, OK, exitCode)

	expected := allLines(
		line(
			filepath.Join(temp, "1.txt"),
			"0",
			"2020-01-01T00:00:00.000000+00:00",
			"d41d8cd98f00b204e9800998ecf8427e",
			"da39a3ee5e6b4b0d3255bfef95601890afd80709",
			"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			"1.txt",
		),
		line(
			filepath.Join(temp, "a")+string(filepath.Separator),
			"",
			"",
			"",
			"",
			"",
			"a"+string(filepath.Separator),
		),
		line(
			filepath.Join(temp, "a", "a.txt"),
			"1",
			"2020-12-21T11:12:21.000000+00:00",
			"9dd4e461268c8034f5c8564e155c67a6",
			"11f6ad8ec52a2984abaafd7c3b516503785c2072",
			"2d711642b726b04401627ca9fbac32f5c8530fb1903cc4db02258717921a4881",
			filepath.Join("a", "a.txt"),
		),
		line(
			filepath.Join(temp, "a", "b.txt"),
			"10",
			"2020-12-20T00:00:00.000000+00:00",
			"336311a016184326ddbdd61edd4eeb52",
			"ff9ee043d85595eb255c05dfe32ece02a53efbb2",
			"fc11d6f28e59d3cc33c0b14ceb644bf0902ebd63d61218dffe9e7dac7c254542",
			filepath.Join("a", "b.txt"),
		),
		line(
			filepath.Join(temp, "a", "xxx")+string(filepath.Separator),
			"",
			"",
			"",
			"",
			"",
			filepath.Join("a", "xxx")+string(filepath.Separator),
		),
		line(
			filepath.Join(temp, "a", "xxx", "x.txt"),
			"20",
			"2019-01-01T12:34:56.000000+00:00",
			"baf1da0e2b9065ab5edd36ca00ed1826",
			"d02e53411e8cb4cd709778f173f7bc9a3455f8ed",
			"d4fc1db665446507dc51b0c9392dd9649291581bfe1b48e241b2b08032b3b647",
			filepath.Join("a", "xxx", "x.txt"),
		),
		line(
			filepath.Join(temp, "a", "xxx", "yyy")+string(filepath.Separator),
			"",
			"",
			"",
			"",
			"",
			filepath.Join("a", "xxx", "yyy")+string(filepath.Separator),
		),
		line(
			filepath.Join(temp, "a", "xxx", "zzz")+string(filepath.Separator),
			"",
			"",
			"",
			"",
			"",
			filepath.Join("a", "xxx", "zzz")+string(filepath.Separator),
		),
		line(
			filepath.Join(temp, "x")+string(filepath.Separator),
			"",
			"",
			"",
			"",
			"",
			"x"+string(filepath.Separator),
		),
		line(
			filepath.Join(temp, "x", "y")+string(filepath.Separator),
			"",
			"",
			"",
			"",
			"",
			filepath.Join("x", "y")+string(filepath.Separator),
		),
		line(
			filepath.Join(temp, "x", "y", "z")+string(filepath.Separator),
			"",
			"",
			"",
			"",
			"",
			filepath.Join("x", "y", "z")+string(filepath.Separator),
		),
		line(
			filepath.Join(temp, "x", "y", "z", "テスト.txt"),
			"100",
			"2021-03-28T00:12:34.000000+00:00",
			"aed563ecafb4bcc5654c597a421547b2",
			"50e483690ec481f4af7f6fb524b2b99eb1716565",
			"09ecb6ebc8bcefc733f6f2ec44f791abeed6a99edf0cc31519637898aebd52d8",
			filepath.Join("x", "y", "z", "テスト.txt"),
		),
	)
	assert.Equal(t, expected, out.String())
}

func TestRun_IncludeDir_ExcludeFile(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	setupFiles(t, temp)

	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			temp,
			"-a",
			"--include-dir",
			"--exclude-file",
		},
		out,
	)

	// ASSERT
	require.Equal(t, OK, exitCode)

	expected := allLines(
		line(
			filepath.Join(temp, "a")+string(filepath.Separator),
		),
		line(
			filepath.Join(temp, "a", "xxx")+string(filepath.Separator),
		),
		line(
			filepath.Join(temp, "a", "xxx", "yyy")+string(filepath.Separator),
		),
		line(
			filepath.Join(temp, "a", "xxx", "zzz")+string(filepath.Separator),
		),
		line(
			filepath.Join(temp, "x")+string(filepath.Separator),
		),
		line(
			filepath.Join(temp, "x", "y")+string(filepath.Separator),
		),
		line(
			filepath.Join(temp, "x", "y", "z")+string(filepath.Separator),
		),
	)
	assert.Equal(t, expected, out.String())
}

func TestRun_ExcludeFile(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	setupFiles(t, temp)

	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			temp,
			"--exclude-file",
		},
		out,
	)

	// ASSERT
	require.Equal(t, OK, exitCode)
	assert.Equal(t, "", out.String())
}

func TestRun_Level1(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	setupFiles(t, temp)

	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			temp,
			"-a",
			"-r",
			"--include-dir",
			"-l", "1",
		},
		out,
	)

	// ASSERT
	require.Equal(t, OK, exitCode)

	expected := allLines(
		line(
			filepath.Join(temp, "1.txt"),
			"1.txt",
		),
		line(
			filepath.Join(temp, "a")+string(filepath.Separator),
			"a"+string(filepath.Separator),
		),
		line(
			filepath.Join(temp, "x")+string(filepath.Separator),
			"x"+string(filepath.Separator),
		),
	)
	assert.Equal(t, expected, out.String())
}

func TestRun_Level2(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	setupFiles(t, temp)

	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			temp,
			"-a",
			"-r",
			"--include-dir",
			"-l", "2",
		},
		out,
	)

	// ASSERT
	require.Equal(t, OK, exitCode)

	expected := allLines(
		line(
			filepath.Join(temp, "1.txt"),
			"1.txt",
		),
		line(
			filepath.Join(temp, "a")+string(filepath.Separator),
			"a"+string(filepath.Separator),
		),
		line(
			filepath.Join(temp, "a", "a.txt"),
			filepath.Join("a", "a.txt"),
		),
		line(
			filepath.Join(temp, "a", "b.txt"),
			filepath.Join("a", "b.txt"),
		),
		line(
			filepath.Join(temp, "a", "xxx")+string(filepath.Separator),
			filepath.Join("a", "xxx")+string(filepath.Separator),
		),
		line(
			filepath.Join(temp, "x")+string(filepath.Separator),
			"x"+string(filepath.Separator),
		),
		line(
			filepath.Join(temp, "x", "y")+string(filepath.Separator),
			filepath.Join("x", "y")+string(filepath.Separator),
		),
	)
	assert.Equal(t, expected, out.String())
}

func TestRun_Level3(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	setupFiles(t, temp)

	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			temp,
			"-a",
			"-r",
			"--include-dir",
			"-l", "3",
		},
		out,
	)

	// ASSERT
	require.Equal(t, OK, exitCode)

	expected := allLines(
		line(
			filepath.Join(temp, "1.txt"),
			"1.txt",
		),
		line(
			filepath.Join(temp, "a")+string(filepath.Separator),
			"a"+string(filepath.Separator),
		),
		line(
			filepath.Join(temp, "a", "a.txt"),
			filepath.Join("a", "a.txt"),
		),
		line(
			filepath.Join(temp, "a", "b.txt"),
			filepath.Join("a", "b.txt"),
		),
		line(
			filepath.Join(temp, "a", "xxx")+string(filepath.Separator),
			filepath.Join("a", "xxx")+string(filepath.Separator),
		),
		line(
			filepath.Join(temp, "a", "xxx", "x.txt"),
			filepath.Join("a", "xxx", "x.txt"),
		),
		line(
			filepath.Join(temp, "a", "xxx", "yyy")+string(filepath.Separator),
			filepath.Join("a", "xxx", "yyy")+string(filepath.Separator),
		),
		line(
			filepath.Join(temp, "a", "xxx", "zzz")+string(filepath.Separator),
			filepath.Join("a", "xxx", "zzz")+string(filepath.Separator),
		),
		line(
			filepath.Join(temp, "x")+string(filepath.Separator),
			"x"+string(filepath.Separator),
		),
		line(
			filepath.Join(temp, "x", "y")+string(filepath.Separator),
			filepath.Join("x", "y")+string(filepath.Separator),
		),
		line(
			filepath.Join(temp, "x", "y", "z")+string(filepath.Separator),
			filepath.Join("x", "y", "z")+string(filepath.Separator),
		),
	)
	assert.Equal(t, expected, out.String())
}

func TestRun_MultiDir(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	setupFiles(t, temp)

	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-a",
			filepath.Join(temp, "a", "xxx"),
			filepath.Join(temp, "x", "y", "z"),
		},
		out,
	)

	// ASSERT
	require.Equal(t, OK, exitCode)

	expected := allLines(
		line(filepath.Join(temp, "a", "xxx", "x.txt")),
		line(filepath.Join(temp, "x", "y", "z", "テスト.txt")),
	)
	assert.Equal(t, expected, out.String())
}

func TestRun_DirNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	setupFiles(t, temp)

	out := new(bytes.Buffer)

	targetDir := filepath.Join(temp, "___") // 存在しない

	// ACT
	exitCode := run(
		[]string{
			targetDir,
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)
	assert.Contains(t, out.String(), "Error: ") // 実況環境で異なるので
	assert.Contains(t, out.String(), targetDir)
}

func TestRun_Help(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	setupFiles(t, temp)

	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-h",
		},
		out,
	)

	// ASSERT
	require.Equal(t, OK, exitCode)

	expected := `filist vdev (dev)

Usage: filist [flags] directory ...

Flags
  -r, --rel            Print relative path (If neither 'rel' nor 'abs' is specified, 'rel' will be printed first column.)
  -a, --abs            Print absolute path
  -s, --size           Print file size
  -m, --mtime          Print modification time
  -M, --md5            Print MD5 hash
  -S, --sha1           Print SHA-1 hash
      --sha256         Print SHA-256 hash
      --include-dir    Include directories
      --exclude-file   Exclude files
  -l, --level int      Number of directory level (Default is unlimited)
  -h, --help           Help
`
	assert.Equal(t, expected, out.String())
}

func TestRun_NoArgs(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	setupFiles(t, temp)

	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	expected := `filist vdev (dev)

Usage: filist [flags] directory ...

Flags
  -r, --rel            Print relative path (If neither 'rel' nor 'abs' is specified, 'rel' will be printed first column.)
  -a, --abs            Print absolute path
  -s, --size           Print file size
  -m, --mtime          Print modification time
  -M, --md5            Print MD5 hash
  -S, --sha1           Print SHA-1 hash
      --sha256         Print SHA-256 hash
      --include-dir    Include directories
      --exclude-file   Exclude files
  -l, --level int      Number of directory level (Default is unlimited)
  -h, --help           Help
`
	assert.Equal(t, expected, out.String())
}

func TestRelPath(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	filePath, info := setupFile(t, temp, "hoge.txt", "ABCDEFG", "")

	// ACT
	result, err := getRelPath(temp, filePath, info)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, "hoge.txt", result)
}

func TestAbsPath(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	filePath, info := setupFile(t, temp, "hoge.txt", "ABCDEFG", "")

	// ACT
	result, err := getAbsPath(temp, filePath, info)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, filePath, result)
}

func TestGetSize(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	filePath, info := setupFile(t, temp, "hoge.txt", "ABCDEFG", "")

	// ACT
	result, err := getSize(temp, filePath, info)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, "7", result)
}

func TestGetMtime(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	filePath, info := setupFile(t, temp, "hoge.txt", "ABCDEFG", "2011-01-02T12:13:14")

	// ACT
	result, err := getMtime(temp, filePath, info)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, "2011-01-02T12:13:14.000000+00:00", result)
}

func TestCalcMd5(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	filePath, info := setupFile(t, temp, "hoge.txt", "ABCDEFG", "")

	// ACT
	result, err := calcMd5(temp, filePath, info)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, "bb747b3df3130fe1ca4afa93fb7d97c9", result)
}

func TestCalcMd5_empty(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	filePath, info := setupFile(t, temp, "hoge.txt", "", "")

	// ACT
	result, err := calcMd5(temp, filePath, info)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, "d41d8cd98f00b204e9800998ecf8427e", result)
}

func TestCalcSha1(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	filePath, info := setupFile(t, temp, "hoge.txt", "ABCDEFG", "")

	// ACT
	result, err := calcSha1(temp, filePath, info)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, "93be4612c41d23af1891dac5fd0d535736ffc4e3", result)
}

func TestCalcSha1_empty(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	filePath, info := setupFile(t, temp, "hoge.txt", "", "")

	// ACT
	result, err := calcSha1(temp, filePath, info)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, "da39a3ee5e6b4b0d3255bfef95601890afd80709", result)
}

func TestCalcSha256(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	filePath, info := setupFile(t, temp, "hoge.txt", "ABCDEFG", "")

	// ACT
	result, err := calcSha256(temp, filePath, info)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, "e9a92a2ed0d53732ac13b031a27b071814231c8633c9f41844ccba884d482b16", result)
}

func TestCalcSha256_empty(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	filePath, info := setupFile(t, temp, "hoge.txt", "", "")

	// ACT
	result, err := calcSha256(temp, filePath, info)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", result)
}

func setupDir(t *testing.T, dir string) {

	err := os.MkdirAll(dir, 0777)
	require.NoError(t, err)
}

func setupFile(t *testing.T, dir string, name string, contents string, mtimeStr string) (string, fs.FileInfo) {

	setupDir(t, dir)

	filePath := filepath.Join(dir, name)

	file, err := os.Create(filePath)
	require.NoError(t, err)
	defer file.Close()

	_, err = file.Write([]byte(contents))
	require.NoError(t, err)

	mtime := time.Now()
	if mtimeStr != "" {
		mtime, err = time.Parse(time.RFC3339, mtimeStr+"Z")
		require.NoError(t, err)
	}

	err = os.Chtimes(filePath, mtime, mtime)
	require.NoError(t, err)

	stat, err := file.Stat()
	require.NoError(t, err)

	return filePath, stat
}

func setupFiles(t *testing.T, root string) {

	setupFile(t, root, "1.txt", strings.Repeat("x", 0), "2020-01-01T00:00:00")
	setupFile(t, filepath.Join(root, "a"), "a.txt", strings.Repeat("x", 1), "2020-12-21T11:12:21")
	setupFile(t, filepath.Join(root, "a"), "b.txt", strings.Repeat("x", 10), "2020-12-20T00:00:00")
	setupFile(t, filepath.Join(root, "a/xxx"), "x.txt", strings.Repeat("x", 20), "2019-01-01T12:34:56")
	setupDir(t, filepath.Join(root, "a/xxx/yyy"))
	setupDir(t, filepath.Join(root, "a/xxx/zzz"))
	setupFile(t, filepath.Join(root, "x/y/z"), "テスト.txt", strings.Repeat("x", 100), "2021-03-28T00:12:34")
}

func line(items ...string) string {
	return strings.Join(items, "\t") + "\n"
}
func allLines(lines ...string) string {
	return strings.Join(lines, "")
}

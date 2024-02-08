package fsutils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
)

var (
	ErrNotFile      = errors.New("path is not a file")
	ErrNotDir       = errors.New("path is not a directory")
	ErrDirNotExists = errors.New("directory not exists")
)

// Returns if given path exists and
// any os error occurs (except os.ErrNotExist).
// Before using the value check error first.
func Exists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil // not exists
		}
		return false, err // error
	}
	return true, nil // exists
}

// Returns if a file exists in given path and
// any os error occurs (except os.ErrNotExist).
// Returns ErrNotFile if the path is an existing directory.
// Before using the value check error first.
// For directory checks better use ExistsDir.
func ExistsFile(path string) (bool, error) {
	var inf fs.FileInfo
	inf, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // not exists
		}
		return false, err // error
	}
	if inf.IsDir() {
		return false, ErrNotFile //exists but not file
	}
	return true, nil // exists and file
}

// Returns if a file exists in given path and
// any os error occurs (except os.ErrNotExist).
// Returns ErrNotFile if the path is an existing directory.
// Before using the value check error first.
// For directory checks better use ExistsDir.
func ExistsDir(path string) (bool, error) {
	var inf fs.FileInfo
	inf, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // not exists
		}
		return false, err // error
	}
	if !inf.IsDir() {
		return false, ErrNotDir //exists but not file
	}
	return true, nil // exists and file
}

// TouchAll mimics os.MkdirAll.
// Creates a file if it not exists in given path.
// Returns corresponding error if given dir not exists or not a dir at all.
func TouchAll(p string) error {
	dir := path.Dir(p) // split paths as dir and file
	ex, err := ExistsDir(dir)
	if err != nil {
		return err
	}
	if !ex {
		return ErrDirNotExists // dir not exists dont proceed further
	}

	ex, err = ExistsFile(p)
	if err != nil {
		return err
	}
	if ex {
		return nil // file exists
	}
	_, err = os.Create(p) // create file
	return err
}

// Base 10 humanized byte strings. Rounding done at number side only.
// Example: 999.9 kB will be shown as 999.9 kB
// but 999.91 kB will be shown as 1000.kB not 1.0 MB
func HumanizedSizeB10(b int64) string {
	var unit int64 = 1000

	if b < unit {
		return fmt.Sprintf("%d B", b)
	}

	div, exp := unit, 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPEZY"[exp])
}

// Base 2 humanized byte strings. Rounding done at number side only.
// Example: 1023.9 kB will be shown as 1023.9 kB
// but 1023.91 kB will be shown as 1024.kB not 1.0 MB
func HumanizedSizeB2(b int64) string {
	var unit int64 = 1024

	if b < unit {
		return fmt.Sprintf("%d B", b)
	}

	var div int64 = unit * unit
	var i int = 0

	for k := 0; k < 8; k++ {
		if b < div {
			i = k
			break
		}
		div *= unit
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div/unit), "kMGTPEZY"[i])
}

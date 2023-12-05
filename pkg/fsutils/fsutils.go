package fsutils

import (
	"errors"
	"fmt"
	"os"
)

func Exists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func ExistFile(path string) (bool, error) {
	if ex, err := Exists(path); err != nil || !ex {
		return ex, err
	}

	inf, err := os.Stat(path)
	if err != nil {
		return true, err
	}
	return !inf.IsDir(), err
}

func ExistsDir(path string) (bool, error) {
	if ex, err := Exists(path); err != nil || !ex {
		return ex, err
	}

	inf, err := os.Stat(path)
	if err != nil {
		return true, err
	}
	return inf.IsDir(), err
}

func TouchAll(path string) error {
	ex, err := Exists(path)
	if err != nil {
		return err
	}
	if ex {
		return nil
	}
	_, err = os.Create(path)
	return err

}

func HumanizedSize(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

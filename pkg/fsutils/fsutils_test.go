package fsutils_test

import (
	"math"
	"os"
	"path"
	"testing"

	"github.com/kmg7/fson/pkg/fsutils"
	"github.com/stretchr/testify/assert"
)

type testDir struct {
	Path            string
	ExistingDir     string
	ExistingFile    string
	NotExistingDir  string
	NotExistingFile string
}

func TestHumanizedSizeB10(t *testing.T) {
	//Arrange
	t.Parallel()
	tests := map[string]struct {
		b    int64
		want string
	}{
		"B-smallest":  {b: 1, want: "1 B"},
		"B-biggest":   {b: 999, want: "999 B"},
		"kB-smallest": {b: 1000, want: "1.0 kB"},
		"kB-biggest":  {b: sizeBiggest(1000, 1), want: "1000.0 kB"},
		"MB-smallest": {b: sizeSmallest(1000, 2), want: "1.0 MB"},
		"MB-biggest":  {b: sizeBiggest(1000, 2), want: "1000.0 MB"},
		"GB-smallest": {b: sizeSmallest(1000, 3), want: "1.0 GB"},
		"GB-biggest":  {b: sizeBiggest(1000, 3), want: "1000.0 GB"},
	}

	for name, test := range tests {
		tr := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := fsutils.HumanizedSizeB10(tr.b)
			assert.Equal(t, tr.want, got)
		})
	}
}

func TestHumanizedSizeB2(t *testing.T) {
	//Arrange
	t.Parallel()

	tests := map[string]struct {
		b    int64
		want string
	}{
		"B-smallest":   {b: 1, want: "1 B"},
		"B-biggest":    {b: 1023, want: "1023 B"},
		"kiB-smallest": {b: 1024, want: "1.0 kiB"},
		"kiB-biggest":  {b: sizeBiggest(1024, 1), want: "1024.0 kiB"},
		"MiB-smallest": {b: sizeSmallest(1024, 2), want: "1.0 MiB"},
		"MiB-biggest":  {b: sizeBiggest(1024, 2), want: "1024.0 MiB"},
		"GiB-smallest": {b: sizeSmallest(1024, 3), want: "1.0 GiB"},
		"GiB-biggest":  {b: sizeBiggest(1024, 3), want: "1024.0 GiB"},
	}

	for name, test := range tests {
		tr := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := fsutils.HumanizedSizeB2(tr.b)
			assert.Equal(t, tr.want, got)
		})
	}
}
func TestExist(t *testing.T) {
	//Arrange
	t.Parallel()
	temp, err := arrangeTestDir()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer func() {
		if temp != nil || temp.Path != "" {
			err := os.RemoveAll(temp.Path)
			if err != nil {
				t.Logf("Test suite was unable to delete temp directory. Path: '%v', Error '%v'", temp.Path, err.Error())
			}
		}
	}()

	tests := map[string]struct {
		p       string
		wantErr error
		want    bool
	}{
		"existing-file":          {p: temp.ExistingFile, want: true},
		"existing-directory":     {p: temp.ExistingDir, want: true},
		"not-existing-file":      {p: temp.NotExistingFile, want: false},
		"not-existing-directory": {p: temp.NotExistingDir, want: false},
	}

	for name, test := range tests {
		tr := test
		t.Run(name, func(t *testing.T) {
			got, gotErr := fsutils.Exists(tr.p)
			if tr.wantErr != nil {
				assert.Equal(t, tr.wantErr, gotErr)
				return
			}
			assert.Nil(t, gotErr)
			assert.Equal(t, tr.want, got)
		})
	}

}

func TestExistsFile(t *testing.T) {
	//Arrange
	t.Parallel()
	temp, err := arrangeTestDir()
	defer func() {
		if temp != nil || temp.Path != "" {
			err := os.RemoveAll(temp.Path)
			if err != nil {
				t.Logf("Test suite was unable to delete temp directory. Path: '%v', Error '%v'", temp.Path, err.Error())
			}
		}
	}()
	if err != nil {
		t.Fatal(err.Error())
	}

	tests := map[string]struct {
		p       string
		wantErr error
		want    bool
	}{
		"existing-file":          {p: temp.ExistingFile, want: true},
		"existing-directory":     {p: temp.ExistingDir, want: false, wantErr: fsutils.ErrNotFile},
		"not-existing-file":      {p: temp.NotExistingFile, want: false},
		"not-existing-directory": {p: temp.NotExistingDir, want: false},
	}

	for name, test := range tests {
		tr := test
		t.Run(name, func(t *testing.T) {
			got, gotErr := fsutils.ExistsFile(tr.p)
			if tr.wantErr != nil {
				assert.Equal(t, tr.wantErr, gotErr)
				return
			}
			assert.Nil(t, gotErr)
			assert.Equal(t, tr.want, got)
		})
	}

}

func TestExistsDir(t *testing.T) {
	//Arrange
	t.Parallel()
	temp, err := arrangeTestDir()
	defer func() {
		if temp != nil || temp.Path != "" {
			err := os.RemoveAll(temp.Path)
			if err != nil {
				t.Logf("Test suite was unable to delete temp directory. Path: '%v', Error '%v'", temp.Path, err.Error())
			}
		}
	}()
	if err != nil {
		t.Fatal(err.Error())
	}

	tests := map[string]struct {
		p       string
		wantErr error
		want    bool
	}{
		"existing-file":          {p: temp.ExistingFile, want: false, wantErr: fsutils.ErrNotDir},
		"existing-directory":     {p: temp.ExistingDir, want: true},
		"not-existing-file":      {p: temp.NotExistingFile, want: false},
		"not-existing-directory": {p: temp.NotExistingDir, want: false},
	}

	for name, test := range tests {
		tr := test
		t.Run(name, func(t *testing.T) {
			got, gotErr := fsutils.ExistsDir(tr.p)
			if tr.wantErr != nil {
				assert.Equal(t, tr.wantErr, gotErr)
				return
			}
			assert.Nil(t, gotErr)
			assert.Equal(t, tr.want, got)
		})
	}

}
func TestTouchAll(t *testing.T) {
	//Arrange
	t.Parallel()
	temp, err := arrangeTestDir()
	defer func() {
		if temp != nil || temp.Path != "" {
			err := os.RemoveAll(temp.Path)
			if err != nil {
				t.Logf("Test suite was unable to delete temp directory. Path: '%v', Error '%v'",
					temp.Path, err.Error())
			}
		}
	}()
	if err != nil {
		t.Fatal(err.Error())
	}

	tests := map[string]struct {
		p    string
		want error
	}{
		"existing-file": { // dir exists file exists
			p:    temp.ExistingFile,
			want: nil},
		"not-existing-file": { // dir exists file not exists
			p:    temp.NotExistingFile,
			want: nil},
		"not-existing-directory": { // dir not exists file also not exists
			p:    path.Join(temp.NotExistingDir, "file.file"),
			want: fsutils.ErrDirNotExists},
	}

	for name, test := range tests {
		tr := test
		t.Run(name, func(t *testing.T) {
			got := fsutils.TouchAll(tr.p)
			assert.Equal(t, tr.want, got)
		})
	}

}

func BenchmarkExists(b *testing.B) {
	temp, err := arrangeTestDir()
	defer func() {
		if temp != nil || temp.Path != "" {
			err := os.RemoveAll(temp.Path)
			if err != nil {
				b.Logf("Test suite was unable to delete temp directory. Path: '%v', Error '%v'",
					temp.Path, err.Error())
			}
		}
	}()
	if err != nil {
		b.Fatal(err.Error())
	}
	benchs := map[string]string{
		"existing-file":          temp.ExistingFile,
		"existing-directory":     temp.ExistingDir,
		"not-existing-file":      temp.NotExistingFile,
		"not-existing-directory": temp.NotExistingDir,
	}

	for name, p := range benchs {
		bp := p
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				fsutils.Exists(bp)
			}
		})
	}

}

func BenchmarkExistsFile(b *testing.B) {
	temp, err := arrangeTestDir()
	defer func() {
		if temp != nil || temp.Path != "" {
			err := os.RemoveAll(temp.Path)
			if err != nil {
				b.Logf("Test suite was unable to delete temp directory. Path: '%v', Error '%v'",
					temp.Path, err.Error())
			}
		}
	}()
	if err != nil {
		b.Fatal(err.Error())
	}
	benchs := map[string]string{
		"existing-file":          temp.ExistingFile,
		"existing-directory":     temp.ExistingDir,
		"not-existing-file":      temp.NotExistingFile,
		"not-existing-directory": temp.NotExistingDir,
	}

	for name, p := range benchs {
		bp := p
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				fsutils.ExistsFile(bp)
			}
		})
	}

}
func BenchmarkExistsDir(b *testing.B) {
	temp, err := arrangeTestDir()
	defer func() {
		if temp != nil || temp.Path != "" {
			err := os.RemoveAll(temp.Path)
			if err != nil {
				b.Logf("Test suite was unable to delete temp directory. Path: '%v', Error '%v'",
					temp.Path, err.Error())
			}
		}
	}()
	if err != nil {
		b.Fatal(err.Error())
	}
	benchs := map[string]string{
		"existing-file":          temp.ExistingFile,
		"existing-directory":     temp.ExistingDir,
		"not-existing-file":      temp.NotExistingFile,
		"not-existing-directory": temp.NotExistingDir,
	}

	for name, p := range benchs {
		bp := p
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				fsutils.ExistsDir(bp)
			}
		})
	}
}
func BenchmarkTouchAll(b *testing.B) {
	temp, err := arrangeTestDir()
	defer func() {
		if temp != nil || temp.Path != "" {
			err := os.RemoveAll(temp.Path)
			if err != nil {
				b.Logf("Test suite was unable to delete temp directory. Path: '%v', Error '%v'",
					temp.Path, err.Error())
			}
		}
	}()
	if err != nil {
		b.Fatal(err.Error())
	}
	benchs := map[string]string{
		"existing-file":          temp.ExistingFile,
		"existing-directory":     temp.ExistingDir,
		"not-existing-file":      temp.NotExistingFile,
		"not-existing-directory": temp.NotExistingDir,
	}

	for name, p := range benchs {
		bp := p
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				fsutils.TouchAll(bp)
			}
		})
	}
}

func arrangeTestDir() (*testDir, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	temp, err := os.MkdirTemp(wd, "temp_*")
	if err != nil {
		return nil, err
	}
	td := testDir{
		Path:            temp,
		ExistingDir:     path.Join(temp, "ExistingDir"),
		ExistingFile:    path.Join(temp, "Existing.File"),
		NotExistingDir:  path.Join(temp, "NotExistingDir"),
		NotExistingFile: path.Join(temp, "NotExisting.File"),
	}

	if err := os.Mkdir(td.ExistingDir, 0666); err != nil {
		return &td, err
	}
	if _, err := os.Create(td.ExistingFile); err != nil {
		return &td, err
	}
	return &td, nil
}

func sizeBiggest(base, pow float64) int64 {
	return int64(math.Pow(base, pow+1) - 1)
}

func sizeSmallest(base, pow float64) int64 {
	return int64(math.Pow(base, pow))
}

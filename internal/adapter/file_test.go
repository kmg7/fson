package adapter_test

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/kmg7/fson/internal/adapter"
	"github.com/stretchr/testify/assert"
)

type TestFile struct {
	Field1 string `json:"field1"`
	path   string
}

func (t *TestFile) FilePath() string {
	return t.path
}

func TestParseAndSave(t *testing.T) {
	p, err := arrangeTestDir()
	if err != nil {
		t.Fatal("Test suite failed to arrange test dir check current directory")
	}
	defer func() {
		if p != "" {
			err := os.RemoveAll(p)
			if err != nil {
				t.Logf("Test suite was unable to delete temp directory. Path: '%v', Error '%v'",
					p, err.Error())
			}
		}
	}()

	t.Run("JSON", func(t *testing.T) {
		fa := adapter.File{
			Parse:   json.Marshal,
			Unparse: json.Unmarshal,
		}
		tf := &TestFile{
			Field1: "testing",
			path:   path.Join(p, "t.json"),
		}
		err := fa.ParseAndSave(tf)
		assert.Nil(t, err)

		data, err := os.ReadFile(tf.FilePath())
		assert.Nil(t, err)

		af := &TestFile{}
		err = json.Unmarshal(data, af)
		assert.Nil(t, err)
		assert.Equal(t, tf.Field1, af.Field1)
	})

}

func TestReadAndParse(t *testing.T) {
	p, err := arrangeTestDir()
	if err != nil {
		t.Fatal("Test suite failed to arrange test dir check current directory")
	}
	defer func() {
		if p != "" {
			err := os.RemoveAll(p)
			if err != nil {
				t.Logf("Test suite was unable to delete temp directory. Path: '%v', Error '%v'",
					p, err.Error())
			}
		}
	}()

	t.Run("NotExistingFile", func(t *testing.T) {
		fa := adapter.File{
			Parse:   json.Marshal,
			Unparse: json.Unmarshal,
		}
		tf := &TestFile{
			Field1: "testing",
			path:   path.Join(p, "t.json"),
		}
		err := fa.ReadAndParse(tf)
		assert.NotNil(t, err)
		// assert.Equal(t, os.ErrNotExist.Error(), err.Error())
	})

	t.Run("JSON", func(t *testing.T) {
		fa := adapter.File{
			Parse:   json.Marshal,
			Unparse: json.Unmarshal,
		}
		filePath := path.Join(p, "t.json")
		tf := &TestFile{
			Field1: "testing",
			path:   filePath,
		}

		file, err := os.Create(filePath)
		if err != nil {
			t.Fatal(err.Error())
		}
		defer file.Close()

		// parse file to intended format
		data, err := json.Marshal(tf)
		if err != nil {
			t.Fatal(err.Error())
		}

		// write parsed bytes to file
		if _, err := file.Write(data); err != nil {
			t.Fatal(err.Error())
		}

		af := &TestFile{path: filePath}
		err = fa.ReadAndParse(af)
		assert.Nil(t, err)
		assert.Equal(t, tf.Field1, af.Field1)
	})
}

// Dont forget to delete test dir.
func arrangeTestDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return os.MkdirTemp(wd, "temp_*")
}

// Since logger.Fatal calls os.Exit directly it fails the test.
// I tried runnin the test in a subprocess but it messes up with coverprofile etc
// So i left t.Fatal untested.

package logger_test

import (
	"io"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/kmg7/fson/internal/logger"
	"github.com/stretchr/testify/assert"
)

type TestLoggerWriter struct {
	logs string
}

func (t *TestLoggerWriter) Write(p []byte) (n int, err error) {
	t.logs = t.logs + string(p)
	return len(p), nil
}

func TestLogIOWriter(t *testing.T) {
	testStdout := &TestLoggerWriter{}
	testLog := logger.New(logger.Options{
		Stdout: testStdout,
	})
	t.Run("Info", func(t *testing.T) {
		ls1 := "LogSample1"
		testLog.Info(ls1)
		assert.Contains(t, testStdout.logs, ls1)
		assert.Contains(t, strings.ToLower(testStdout.logs), "info")
	})

	t.Run("Warn", func(t *testing.T) {
		ls2 := "LogSample2"
		testLog.Warn(ls2)
		assert.Contains(t, testStdout.logs, ls2)
		assert.Contains(t, strings.ToLower(testStdout.logs), "warn")
	})

	// t.Run("Debug", func(t *testing.T) {
	// 	ls3 := "LogSample3"
	// 	testLog.Debug(ls3)
	// 	assert.Contains(t, testStdout.logs, ls3)
	// 	assert.Contains(t, strings.ToLower(testStdout.logs), "debug")
	// })

	t.Run("Error", func(t *testing.T) {
		ls4 := "LogSample4"
		testLog.Error(ls4)
		assert.Contains(t, testStdout.logs, ls4)
		assert.Contains(t, strings.ToLower(testStdout.logs), "error")
	})

	// it tests if it panics but prevents other tests running.
	// t.Run("Fatal", func(t *testing.T) {
	// ls5 := "LogSample5"
	// 	assert.Panics(t, func() {
	// 		testLog.Fatal(ls5)
	// 	})
	// 	assert.Contains(t, testStdout.logs, ls5)
	// 	assert.Contains(t, strings.ToLower(testStdout.logs), "fatal")
	// })
}

func TestFileLogging(t *testing.T) {
	var logFile *os.File
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal("Cannot arrange")
	}
	td, err := os.MkdirTemp(wd, "temp_test_*")
	if err != nil {
		t.Fatal("Cannot arrange")
	}

	// delete arranged stuff
	defer func() {
		if err = os.RemoveAll(td); err != nil {
			t.Logf("Cannot delete temp test dir. Path: %v", td)
		}
		t.Log("Tempdir deleted")
	}()

	// open a file for testing json file logs
	filePath := path.Join(td, "file.log")
	logFile, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0740)
	if err != nil {
		t.Fatal("Cannot arrange")
	}

	fileLog := logger.New(logger.Options{
		Files: []io.Writer{
			logFile,
		},
	})

	t.Run("Error", func(t *testing.T) {
		ls1 := "Testing"
		fileLog.Error(ls1)
		b, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		logs := string(b)
		assert.Contains(t, logs, ls1)
		assert.Contains(t, strings.ToLower(logs), "error")

	})

}

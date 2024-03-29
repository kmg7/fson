package adapter

import (
	"os"
)

type File struct {
	Parse   Parser   //ex json.Marshall
	Unparse Unparser //ex json.Unmarshall
}

// It utilizes os.Create and file.Write methods. Any occuring error
// will be returned directly. Same goes for Parser errors.
func (fa *File) ParseAndSave(f AdapterFile) error {
	// create file
	file, err := os.Create(f.FilePath())
	if err != nil {
		return err
	}
	defer file.Close()

	// parse file to intended format
	data, err := fa.Parse(f)
	if err != nil {
		return err
	}

	// write parsed bytes to file
	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}

// It utilizes os.ReadFile method. Any error occurs
// will be returned untouched.
func (fa *File) ReadAndParse(f AdapterFile) error {
	// read the file
	data, err := os.ReadFile(f.FilePath())
	if err != nil {
		return err
	}

	// unparse the file
	return fa.Unparse(data, f)
}

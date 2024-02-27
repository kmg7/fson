package adapter

type Parser = func(v any) ([]byte, error)
type Unparser = func(data []byte, v any) error

type FileAdapter interface {
	ParseAndSave(f AdapterFile) error
	ReadAndParse(f AdapterFile) error
}

type AdapterFile interface {
	FilePath() string
}

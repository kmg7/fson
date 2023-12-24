package scout

/*
* scout will stream downloading files etc
* make use of filehandlers methods
* manage current files and directories user added
*
 */

import (
	"os"
	"path"
	"time"

	"github.com/kmg7/fson/pkg/fsutils"
)

type FileToServe struct {
	Name    string
	Path    string
	Sizehr  string
	Mdfdate time.Time
	Headers map[string]string
}

type DirToServe struct {
	Name    string
	Path    string
	Mdfdate time.Time
	Files   []*FileToServe
	Subdirs []*DirToServe
}

var serve *DirToServe

func orDone(done, c <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

func StartScouting() {

}

func ReadDir(p string) (*DirToServe, error) {
	dir := DirToServe{Path: p}
	dirs, err := os.ReadDir(p)
	if err != nil {
		return &dir, err
	}

	for _, d := range dirs {
		sdi, err := d.Info() //* pull information about dir
		if err != nil {
			continue
		}
		if sdi.IsDir() { //*if this is a sub directory read it
			if sd, err := ReadDir(path.Join(p, sdi.Name())); err == nil {
				dir.Subdirs = append(dir.Subdirs, sd)
			} // *read dir and append to parent dir
			continue
		}
		f := FileToServe{
			Name:    sdi.Name(),
			Path:    path.Join(path.Join(p, sdi.Name())),
			Mdfdate: sdi.ModTime(),
			Sizehr:  fsutils.HumanizedSize(sdi.Size()),
		}
		dir.Files = append(dir.Files, &f)

	}
	return &dir, nil
}

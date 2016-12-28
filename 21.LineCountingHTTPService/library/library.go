package library

import (
	"io/ioutil"
	"path/filepath"
	"os"
	"github.com/8tomat8/GoCourse/21.LineCountingHTTPService/common"
)

var inst *library

type book struct {
	Title string `json:"title"`
	Lines uint `json:"lines"`
}

type library struct {
	Path string
}

func (l *library) GetBook(name string) (b book, err error) {
	file, err := os.Open(l.Path + name)
	defer file.Close()
	if err != nil {
		return b, err
	}
	lines, err := common.CountLines(file)
	b.Lines = lines
	b.Title = name
	return b, nil
}

func New(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic("Could not read directory!")
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".txt" {
			break
		}
		panic("There is no txt files in folder!")
	}
	inst = &library{Path:dir}
}

func Get() *library {
	return inst
}

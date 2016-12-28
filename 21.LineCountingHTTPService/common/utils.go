package common

import (
	"io"
	"bufio"
	"os"
	"io/ioutil"
)

func CountLines(r io.Reader) (uint, error) {
	sc := bufio.NewScanner(r)
	var lines uint
	for sc.Scan() {
		lines++
	}
	return lines, sc.Err()
}

func GetFiles(dir string) (files []os.FileInfo, err error) {
	files, err = ioutil.ReadDir(dir)
	if err != nil {
		panic("Could not read directory!")
	}
	return
}



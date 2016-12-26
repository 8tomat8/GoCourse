package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"io/ioutil"
	"path/filepath"
)

func CountLines(r io.Reader) (int, error) {
	sc := bufio.NewScanner(r)
	var lines int
	for sc.Scan() {
		lines++
	}
	return lines, sc.Err()
}

func CountFileLines(path string) (int, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return 0, err
	}
	count, err := CountLines(f)
	if err != nil {
		return 0, err
	}
	return count, nil

}

func main() {
	dir := os.Args[1]

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic("Could not read directory!")
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".txt" {
			continue
		}
		count, err := CountFileLines(fmt.Sprintf("%v/%v", dir, file.Name()))
		if err != nil {
			log.Fatalln(err)
		} else {
			fmt.Print(file.Name(), "\t", count, "\n")
		}
	}
}

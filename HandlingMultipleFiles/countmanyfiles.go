package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
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
	files := os.Args[1:]

	for _, file := range files {
		count, err := CountFileLines(file)
		if err != nil {
			log.Fatalln(err)
		} else {
			fmt.Print(file, "\t", count, "\n")
		}
	}
}

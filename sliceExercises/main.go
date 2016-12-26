package main

import "fmt"

func main() {

	i := make([]int, 5)               // Declare a variable called i which is a slice of 5 int.
	f := make([]float64, 9)           // Declare a variable called f which is a slice of 9 float64.
	s := []string{"a", "b", "c", "d"} // Declare a variable called s which is a slice of 4 string.

	fmt.Println(len(i), len(f), len(s))
}

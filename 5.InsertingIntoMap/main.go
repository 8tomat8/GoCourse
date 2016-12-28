package main

import (
	"fmt"
)

func main() {
	// Maps can have keys of almost any type.
	//
	// Create a map from the names of days, eg "Monday", "Tuesday", to their
	// numbers. 1 for Monday, 2 for Tuesday.

	// uncomment and complete this declaration
	daynumber := map[string]int{
		"Monday":    1,
		"Tuesday":   2,
		"Wednesday": 3,
		"Thurthday": 4,
		"Friday":    5,
		"Saturday":  6,
		"Sunday":    7,
	}

	// The correct answer should print 3 6
	fmt.Println(daynumber["Wednesday"], daynumber["Saturday"])

	// answer is in src/maps/maps2a/main.go
}

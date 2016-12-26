package simplestrings

import "testing"

const weekdays = "Monday Tuesday Wednesday Thursday Friday"

func TestContains(t *testing.T) {
	var rv bool

	// test that Tuesday is a weekday
	rv = Contains(weekdays, "Tuesday")
	if rv == false {
		t.Error("Tuesday is not a weekday")
	}

	// test that Sunday is not a weekday
	rv = Contains(weekdays, "Sunday")
	if rv == true {
		t.Error("Tuesday is a weekday")
	}

	// test that the string Monday is not found in the empty string
	rv = Contains("", "Monday")
	if rv == true {
		t.Error("String Monday found in the empty string!")
	}

}

func TestIndex(t *testing.T) {
	var rv int

	// test that an empty search string returns 0
	rv = Index(weekdays, "")
	if rv != 0 {
		t.Error("An empty search string does not returns 0!")
	}
}


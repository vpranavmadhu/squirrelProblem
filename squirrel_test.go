package main

import "testing"

func TestPhi(t *testing.T) {
	count := Counts{
		n00: 1,
		n01: 1,
		n10: 1,
		n11: 1,
	}

	expected := 0.0
	result := phi(count)

	if result != expected {
		t.Errorf("Expected %f but  got %f", expected, result)
	}

	result = phi(Counts{n00: 10, n01: 0, n10: 0, n11: 10})
	expected = 1.0
	if result != expected {
		t.Errorf("Expected %f but  got %f", expected, result)
	}

}

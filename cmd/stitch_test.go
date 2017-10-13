package cmd

import (
	"reflect"
	"testing"
)

// TestGetPermutation ensures that we get the correct permutations
func TestGetPermutation(t *testing.T) {
	testing := [][]string{
		[]string{"0"},
		[]string{"0", "1"},
		[]string{"0"},
		[]string{"0", "1"}}
	expecting := [][]string{
		[]string{"0", "0", "0", "0"},
		[]string{"0", "0", "0", "1"},
		[]string{"0", "1", "0", "0"},
		[]string{"0", "1", "0", "1"}}

	//
	result := getPermutation(testing)

	//
	if !reflect.DeepEqual(expecting, result) {
		t.Errorf("Expecting %v, got %v", expecting, result)
	}
}

// TestGetPermutationEmptySlice ensures we get an empty slice returned when an
// empty slice is provided
func TestGetPermutationEmptySlice(t *testing.T) {
	testing := [][]string{}
	expecting := [][]string{}

	//
	result := getPermutation(testing)

	//
	if len(result) != 0 {
		t.Errorf("Expecting %v, got %v", len(expecting), len(result))
	}
}

package main

import (
	"fmt"
	"sort"
)

func Example_sortMapByValue() {
	m := map[string]int{
		"blah": 5,
		"boo":  9,
		"yah":  1,
	}

	sm := SortMapByValue(m)

	for _, v := range sm {
		fmt.Println(v.Value, v.Key)
	}

	// Output:
	// 9 boo
	// 5 blah
	// 1 yah
}

func Example_reverseAdapter() {
	s := []int{5, 2, 6, 3, 1, 4} // Unsorted.

	sort.Sort(sort.IntSlice(s))
	fmt.Println(s)

	sort.Sort(reverseAdapter{Interface: sort.IntSlice(s)})
	fmt.Println(s)

	// Output:
	// [1 2 3 4 5 6]
	// [6 5 4 3 2 1]
}

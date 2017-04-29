package main

import "fmt"

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

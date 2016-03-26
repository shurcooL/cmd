// rune_stats prints counts of total and unique runes from stdin.
// It's helpful for finding non-ASCII characters in files.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func printRuneStats(b []byte) {
	r := []rune(string(b))
	fmt.Printf("Total runes: %v\n", len(r))

	m := map[rune]int{}
	for _, v := range r {
		m[v]++
	}
	fmt.Printf("Total unique runes: %v\n\n", len(m))

	sm := SortMapByKey(m, true)

	for _, v := range sm {
		fmt.Printf("%q (%v)\t%v\n", v.Key, v.Key, v.Value)
	}
}

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}
	printRuneStats(b)
}

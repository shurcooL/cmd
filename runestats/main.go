// runestats prints counts of total and unique runes from stdin.
// It's helpful for finding non-ASCII characters in files.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

func main() {
	flag.Parse()

	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	printRuneStats(stdin)
	return nil
}

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

// Pair is a data structure to hold a key/value pair.
type Pair struct {
	Key   string
	Value int
}

// PairList is a slice of Pairs that implements sort.Interface to sort by Pair.Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// SortMapByValue turns a map into a PairList, then sorts and returns it.
func SortMapByValue(m map[string]int) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(reverseAdapter{Interface: p})
	return p
}

// RuneIntPair is a data structure to hold a key/value pair.
type RuneIntPair struct {
	Key   rune
	Value int
}

// RuneIntPairList is a slice of RuneIntPair that implements sort.Interface to sort by RuneIntPair.Value.
type RuneIntPairList []RuneIntPair

func (p RuneIntPairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p RuneIntPairList) Len() int           { return len(p) }
func (p RuneIntPairList) Less(i, j int) bool { return p[i].Key < p[j].Key }

// SortMapByKey sorts map by key.
func SortMapByKey(m map[rune]int, reverse bool) RuneIntPairList {
	sm := make(RuneIntPairList, len(m))
	i := 0
	for k, v := range m {
		sm[i] = RuneIntPair{k, v}
		i++
	}
	if !reverse {
		sort.Sort(sm)
	} else {
		sort.Sort(reverseAdapter{Interface: sm})
	}
	return sm
}

// reverseAdapter is a reverse sort.Interface adapter.
type reverseAdapter struct {
	sort.Interface
}

// Less returns the opposite of the embedded implementation's Less method.
func (r reverseAdapter) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}

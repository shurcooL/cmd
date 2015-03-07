// Command rune_stats prints counts of total and unique runes from stdin.
package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/shurcooL/go/gists/gist5408736"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}
	gist5408736.PrintRuneStats(string(b))
}

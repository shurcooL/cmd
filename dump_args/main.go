// Dumps the command-line arguments.
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/shurcooL/go-goon"
)

func main() {
	out := goon.SdumpExpr(os.Args[0])  // Program name.
	out += goon.SdumpExpr(os.Args[1:]) // Program arguments.
	out += goon.SdumpExpr(os.Getwd())  // Current working directory.

	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	out += "### Stdin ###\n" + string(stdin)

	fmt.Println(out)

	err = ioutil.WriteFile(filepath.Join(os.TempDir(), "dump_args.txt"), []byte(out), 0644)
	if err != nil {
		panic(err)
	}
}

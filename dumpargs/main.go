// dumpargs dumps the command-line arguments.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Purposefully avoid flag.Parse() here, because we're dealing with
	// dumping low level command-line arguments without any processing.

	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	out := fmt.Sprintf("os.Args[0]:  %#q\n", os.Args[0])  // Program name.
	out += fmt.Sprintf("os.Args[1:]: %#q\n", os.Args[1:]) // Program arguments.
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	out += fmt.Sprintf("os.Getwd():  %#q\n", wd) // Current working directory.

	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	out += "### Stdin ###\n" + string(stdin)

	fmt.Println(out)

	// Write a copy of output to "dumpargs.txt" in temp folder, in case stdout is hard to see.
	err = ioutil.WriteFile(filepath.Join(os.TempDir(), "dumpargs.txt"), []byte(out), 0644)
	return err
}

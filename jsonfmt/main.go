// jsonfmt pretty-prints JSON from stdin.
package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func run() error {
	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	var out bytes.Buffer
	err = json.Indent(&out, in, "", "\t")
	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, &out)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

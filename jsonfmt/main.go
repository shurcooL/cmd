// jsonfmt pretty-prints JSON from stdin.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	flag.Parse()

	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	var out bytes.Buffer
	err = json.Indent(&out, bytes.TrimSpace(in), "", "\t") // Need to trim space in order to consistently print single newline at the end. This is because of bug in json.Indent that couldn't be resolved in a better way due to Go 1 guarantee, see https://github.com/golang/go/issues/13520#issuecomment-162544193.
	if err != nil {
		return err
	}
	err = out.WriteByte('\n')
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, &out)
	return err
}

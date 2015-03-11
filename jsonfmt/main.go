// Pretty-prints JSON from stdin.
package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}

	var out bytes.Buffer
	err = json.Indent(&out, in, "", "\t")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = io.Copy(os.Stdout, &out)
	if err != nil {
		log.Fatalln(err)
	}
}

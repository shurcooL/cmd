// Pretty-prints JSON from stdin.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return
	}

	var out bytes.Buffer
	err = json.Indent(&out, in, "", "\t")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return
	}

	_, err = io.Copy(os.Stdout, &out)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return
	}
}

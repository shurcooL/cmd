// Chef client command-line tool.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/marpaia/chef-golang"
)

func chefConnect() *chef.Chef {
	c, err := chef.Connect()
	if err != nil {
		panic(err)
	}
	c.SSLNoVerify = true
	return c
}

func main() {
	flag.Parse()

	switch args := flag.Args(); {
	case len(args) == 1:
		c := chefConnect()

		results, err := c.Search("node", "role:"+args[0])
		if err != nil {
			panic(err)
		}

		for _, row := range results.Rows {
			var node chef.Node
			err := json.Unmarshal(row, &node)
			if err != nil {
				panic(err)
			}

			fmt.Println(node.Name)
		}
	case len(args) == 2 && args[0] == "ipaddress":
		c := chefConnect()

		results, err := c.Search("node", "ipaddress:"+args[1])
		if err != nil {
			panic(err)
		}

		for _, row := range results.Rows {
			var node chef.Node
			err := json.Unmarshal(row, &node)
			if err != nil {
				panic(err)
			}

			fmt.Println(node.Name)
		}
	default:
		flag.Usage()
		os.Exit(2)
	}
}

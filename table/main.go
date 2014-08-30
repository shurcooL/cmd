// Chef client command-line tool.
package main

import (
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
	args := flag.Args()

	switch {
	case len(args) == 1:
		c := chefConnect()

		results, err := c.Search("node", "role:"+args[0])
		if err != nil {
			panic(err)
		}

		for _, row := range results.Rows {
			row := row.(map[string]interface{})

			fmt.Println(row["name"])
		}
	case len(args) == 2 && args[0] == "ipaddress":
		c := chefConnect()

		results, err := c.Search("node", "ipaddress:"+args[1])
		if err != nil {
			panic(err)
		}

		for _, row := range results.Rows {
			row := row.(map[string]interface{})

			fmt.Println(row["name"])
		}
	default:
		flag.PrintDefaults()
		os.Exit(2)
	}
}

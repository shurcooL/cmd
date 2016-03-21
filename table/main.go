// table is a chef client command-line tool.
// It's similar to knife, but easier to install and run.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/marpaia/chef-golang"
)

func chefConnect() *chef.Chef {
	c, err := chef.Connect()
	if err != nil {
		log.Fatalln(err)
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
			log.Fatalln(err)
		}

		for _, row := range results.Rows {
			var node chef.Node
			err := json.Unmarshal(row, &node)
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Println(node.Name)
		}
	case len(args) == 2 && args[0] == "ipaddress":
		c := chefConnect()

		results, err := c.Search("node", "ipaddress:"+args[1])
		if err != nil {
			log.Fatalln(err)
		}

		for _, row := range results.Rows {
			var node chef.Node
			err := json.Unmarshal(row, &node)
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Println(node.Name)
		}
	case len(args) == 2 && args[0] == "role" && args[1] == "list":
		c := chefConnect()

		roles, err := c.GetRoles()
		if err != nil {
			log.Fatalln(err)
		}

		for role := range roles {
			fmt.Println(role)
		}
	default:
		flag.Usage()
		os.Exit(2)
	}
}

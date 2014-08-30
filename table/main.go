// Chef client command-line tool.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/marpaia/chef-golang"
	"github.com/shurcooL/go/gists/gist7651991"
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
	case len(args) == 2 && args[0] == "ip":
		c := chefConnect()

		nodes, err := c.GetNodes()
		if err != nil {
			panic(err)
		}

		/*for n, _ := range nodes {
			node, ok, err := c.GetNode(n)
			if err != nil {
				panic(err)
			}
			if ok && node.Info.IPAddress == args[1] {
				fmt.Println(node.Name)
			} else if ok {
				fmt.Println("not matched", node.Name, node.Info.IPAddress)
			} else {
				fmt.Println("not found!", n)
			}
		}*/

		inChan := make(chan interface{})
		go func() { // This needs to happen in the background because sending input will be blocked on reading output.
			for n, _ := range nodes {
				inChan <- n
			}
			close(inChan)
		}()
		reduceFunc := func(in interface{}) interface{} {
			n := in.(string)
			node, ok, err := c.GetNode(n)
			if err != nil {
				panic(err)
			}
			if ok && node.Info.IPAddress == args[1] {
				return node.Name
			} /* else if ok {
				fmt.Println("not matched", node.Name, node.Info.IPAddress)
			} else {
				fmt.Println("not found!", n)
			}*/
			return nil
		}
		outChan := gist7651991.GoReduce(inChan, 256, reduceFunc)

		for out := range outChan {
			fmt.Println(out)
		}
	default:
		flag.PrintDefaults()
		os.Exit(2)
	}
}

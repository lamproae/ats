package main

import (
	//"telnet"
	"log"
	"os"
	//"expect"
	"product"
	"dut"
	"cli"
)


func checkErr(err error) {
	if err != nil {
		log.Fatalln("Error:", err)
	}
}

func main() {
	if len(os.Args) != 2{
		log.Printf("Usage: %s PRODUCTNAME", os.Args[0])
		return
	}

	prod, err := product.New(os.Args[1])
	if err != nil {
		log.Fatal("Product not exist check product list: product/product.json")
	}

	c := cli.New(prod)
	d := dut.New(os.Args[1], c)
	c = d.Cli

	c.RunCommand(cli.Command{"config", "show interface", "#"})
	log.Println(c.CommandResult())
	c.RunCommand(cli.Command{"config", "show vlan", "#"})
	log.Println(c.CommandResult())
	c.RunCommand(cli.Command{"config", "show port", "#"})
	log.Println(c.CommandResult())
	c.RunCommand(cli.Command{"config", "show system", "#"})
	log.Println(c.CommandResult())
	c.RunCommand(cli.Command{"bridge", "show vlan", "#"})
	log.Println(c.CommandResult())
	c.RunCommand(cli.Command{"bridge", "port disable 8", "#"})
	log.Println(c.CommandResult())
	c.RunCommand(cli.Command{"bridge", "vlan create 400", "#"})
	log.Println(c.CommandResult())
	c.RunCommand(cli.Command{"config", "show vlan", "#"})
	log.Println(c.CommandResult())
	c.RunCommand(cli.Command{"enable", "show port", "#"})
	log.Println(c.CommandResult())
	c.RunCommand(cli.Command{"bridge", "vlan create 800", "#"})
	log.Println(c.CommandResult())
	c.RunCommand(cli.Command{"enable", "show port", "#"})
	log.Println(c.CommandResult())
}

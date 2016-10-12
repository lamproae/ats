package main

import (
	"telnet"
	"log"
	"os"
	"expect"
	"product"
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

	dst := prod.Address() + ":23"
	t, err := telnet.NewClient(dst)
	checkErr(err)
	t.SetUnixWriteMode(true)

	var data []byte

	expect.Expect(t, prod.Hostname()+" login: ")
	expect.Sendln(t, prod.Username())
	expect.Expect(t, "Password: ")
	expect.Sendln(t, prod.Password())
	expect.Expect(t, prod.Hostname()+">")
	expect.Sendln(t, "sh ver")
	expect.Expect(t, prod.Hostname()+">")
	expect.Sendln(t, "enable")
	expect.Expect(t, prod.Hostname()+"#")
	expect.Sendln(t, "terminal length 0")
	expect.Expect(t, prod.Hostname()+"#")
	expect.Sendln(t, "show interface")
	data, err = t.ReadBytes('#')
	os.Stdout.Write(data)
	expect.Sendln(t, "show port")
	data, err = t.ReadUntil(prod.Hostname()+"#")
	os.Stdout.Write(data)
//	expect.Sendln(t, "show running-config")
//	data, err = t.ReadUntil(prod.Hostname()+"#")
	checkErr(err)
//	os.Stdout.Write(data)
	os.Stdout.WriteString("\n")
}

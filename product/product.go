package product
import (
	"os"
	"log"
	"encoding/json"
)

var productList map[string]*Product

type ProductList struct {
	Products []string 	`products`
}

type DefaultConfig struct {
	Hostname string		`hostname`
	Username string		`username`
	Password string		`password`
	Address string		`address`
	Initmode string		`initmode`
}

type PortInfo struct {
	Ethernet int	`Ethernet`
	GEPON int	`GEPON`
	Trunk int	`Trunk`
}

type SysInfo struct {
	ModelName string	`modelName`
	MemorySize string	`memorySize`
	Port PortInfo		`port`
}

type Product struct {
	Product string	`product`
	Config DefaultConfig `config`
	Sysinfo SysInfo	`sysinfo`
}

func (p *Product) Username() string {
	return p.Config.Username
}

func (p *Product) Password() string {
	return p.Config.Password
}

func (p *Product) Address() string {
	return p.Config.Address
}

func (p *Product) Hostname() string {
	return p.Config.Hostname
}

func (p *Product) Initmode() string {
	return p.Config.Initmode
}

func New(name string) (*Product, error) {
	if _, ok := productList[name]; !ok {
		file, err := os.Open("product/"+name+".json")
		if err != nil {
			log.Fatal(err)
		}

		data := make([]byte, 1000)
		count, err := file.Read(data)
		if err != nil {
			log.Fatal(err)
		}

		config := data[:count]

		var product Product
		err = json.Unmarshal(config, &product)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(product)

		return &product, nil
	}

	return productList[name], nil
}

func ProductDB() map[string]*Product {
	return productList
}

func getProductList() ProductList {
	file, err := os.Open("product/product.json")
	if err != nil {
		log.Fatal(err)
	}

	data := make([]byte, 1000)
	count, err := file.Read(data)
	if err != nil {
		log.Fatal(err)
	}

	p := data[:count]

	var prodlists ProductList
	err = json.Unmarshal(p, &prodlists)
	return prodlists
}

func init () {
	productList = make(map[string]*Product, 100)
	lists := getProductList()
	for _, p := range lists.Products {
		if _, ok := productList[p]; ok {
			log.Print("Conflict Model %s!\n", p)
			continue
		}

		_, err := os.Stat("product/"+p+".json")
		if err != nil {
			log.Panic("There is no configuration file for ", p)
		}
		product, err := New(p)
		if err != nil {
			log.Panic("Create device failed: ", p)
		}
		productList[p] = product 
	}

	log.Println("Supported Products: ")
	log.Printf("Supported Products: %v", productList)
}

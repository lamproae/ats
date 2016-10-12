package atscase

import (
	"cli"
	"log"
	"dut"
	"os"
	"io/ioutil"
	"encoding/json"
)

type Operation struct {
	Name string 	`name`	//Operation on which device
	Commands []cli.Command	`commands`
	Expected []cli.Expected	`expected`
}

type Step struct {
	ops map[string]Operation
	opsList []Operation

	// This operation mast in sequence
	//ops []Operation
}

type Case struct {
	duts map[string]*dut.DUT
	Name string		`name`
	Parent string		`parent`
	Description string	`description`
	Duts []string		`duts`
	precondition []Step
	subCases []Step
	postcondition []Step
}

func New (name string) *Case {
	_, err := os.Stat("cases/L2/Vlan/"+name+".json") 
	if err != nil {
		log.Println("No json file for case: ", name)
		return nil
	}

	data, err := ioutil.ReadFile("cases/L2/Vlan/"+name+".json")
	if err != nil {
		log.Println("Read json file error: ", name)
		return nil
	}

	var c Case
	err = json.Unmarshal(data, &c)
	if err != nil {
		log.Println("Parse json file error: ", name)
	}

	log.Println(c)
	c.BuildCase()
	return &c
}

type StepHelper struct {
	Name string	`name`
	Operations []Operation `operations`
}

func (c *Case) BuildCase() {
	c.duts = make(map[string]*dut.DUT, len(c.Duts))
	for _, d := range c.Duts {
		c.duts[d] = dut.New(d)
	}

	data, err := ioutil.ReadFile("cases/L2/Vlan/"+"sub.json")
	if err != nil {
		log.Println("Read file error ", err.Error())
	}
	var shs []StepHelper
	err = json.Unmarshal(data, &shs)
	if err != nil {
		log.Println("Pares json file error ", err.Error())
	}

	log.Println(shs)

	c.subCases = make([]Step, 0, len(shs))
	for _, s := range shs {
		var step Step
		step.ops = make (map[string]Operation, len(s.Operations))
		step.opsList = make ([]Operation, 0, len(s.Operations))
		for _, o := range s.Operations {
			if _, ok := c.duts[o.Name]; ok {
				step.ops[o.Name] = o
				step.opsList = append(step.opsList, o)
			} else {
				log.Fatal("Unkown dut: ", o.Name)
			}
		}
		c.subCases = append(c.subCases, step)
	}

	log.Println(c)
}

func (c *Case) RunSteps(steps []Step) (expect *cli.Expected, result string, err error) {
	for _, s := range steps {
		for _, o := range s.opsList {
			for _, command := range o.Commands {
				c.duts[o.Name].Cli.RunCommand(command)
			}

			for _, e := range o.Expected {
				r, err := c.duts[o.Name].Cli.Assert(e)
				if err != nil {
					return &e, r, err
				}
			}
		}
	}

	return nil, "", nil
}

func (c *Case) CheckPrecondition() error {
	log.Println("----------------Checking Preconditions---------------")
	expect, r, err := c.RunSteps(c.precondition)
	if err != nil {
		log.Printf("Precondition Check failed: %s expected: %s but get: %s\n", expect.GetCommand(), expect.GetExpected(), r)
		return err
	}
	return nil
}

func (c *Case) CheckPostcondition() error {
	log.Println("----------------Checking Postconditions---------------")
	expect, r, err := c.RunSteps(c.postcondition)
	if err != nil {
		log.Printf("Postcondition Check failed: %s expected: %s but get: %s\n", expect.GetCommand(), expect.GetExpected(), r)
		return err
	}
	return nil
}

func (c *Case) RunSubCases() error {
	log.Println("----------------Running Test---------------")
	expect, r, err := c.RunSteps(c.subCases)
	if err != nil {
		log.Printf("Test Check failed: %s expected: %s but get: %s\n", expect.GetCommand(), expect.GetExpected(), r)
		return err
	}
	return nil
}

func (c *Case) Run() error {
	log.Printf("Running case: %s\n", c.Name)
	err := c.CheckPrecondition() 
	if err != nil {
		return err
	}
	err = c.RunSubCases()
	if err != nil {
		return err
	}
	err = c.CheckPostcondition()
	if err != nil {
		return err
	}

	return nil
}

func init () {

}

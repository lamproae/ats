package dut

import (
	"cli"
)

type DUT struct {
	Name string
	Cli *cli.Cli
}

func New(name string) *DUT {
	cli := cli.New(name)
	//cli.SetHostname(name)

	dut := DUT {
		Name : name,
		Cli : cli,
	}

	return &dut
}


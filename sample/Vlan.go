package main

import (
	"cli"
	"dut"
	"fmt"
	"log"
	"net"
	"regexp"
	"time"
	//	"strconv"
	//	"strings"
)

type Interface struct {
	Name               string
	IsL3               bool
	IsVlan             bool
	IsLoopback         bool
	IsAdminUP          bool
	IsRunning          bool
	IsBroadcastCapable bool
	IsMulticastCapable bool
	IfIndex            int64
	IfMetric           int64
	IfMTU              int64
	Bandwidth          int64
	Description        string
	IPAddress          map[string]*net.IP
	IPConnected        map[string]*net.IPNet
	IPv6Address        map[string]*net.IP
	IPv6Connected      map[string]*net.IPNet
	MACAddress         map[string]*net.HardwareAddr
	Configuration      []string
	ShowInterface      string
	InputPackets       int64
	InputBytes         int64
	InputDropped       int64
	OutputPackets      int64
	OutputBytes        int64
	OutputDropped      int64
}

func (i *Interface) String() string {
	var result string
	result += fmt.Sprintf("Interface name: %s\n", i.Name)
	result += fmt.Sprintf("\tIsL3: %v\n", i.IsL3)
	result += fmt.Sprintf("\tIsVlan: %v\n", i.IsL3)
	result += fmt.Sprintf("\tIsLoopback: %v\n", i.IsLoopback)
	result += fmt.Sprintf("\tIsAdminUP: %v\n", i.IsAdminUP)
	result += fmt.Sprintf("\tIsRunning: %v\n", i.IsRunning)
	result += fmt.Sprintf("\tIsBroadcastCapable: %v\n", i.IsBroadcastCapable)
	result += fmt.Sprintf("\tIsMulticastCapable: %v\n", i.IsMulticastCapable)
	result += fmt.Sprintf("\tIfIndex: %v\n", i.IfIndex)
	result += fmt.Sprintf("\tIfMetric: %v\n", i.IfMetric)
	result += fmt.Sprintf("\tIfMTU: %v\n", i.IfMTU)
	result += fmt.Sprintf("\tBandwidth: %v\n", i.Bandwidth)
	result += fmt.Sprintf("\tDescription: %v\n", i.Description)
	result += fmt.Sprintf("\tIPAddress: %v\n", i.IPAddress)
	result += fmt.Sprintf("\tIPConnected: %v\n", i.IPConnected)
	result += fmt.Sprintf("\tIPv6Address: %v\n", i.IPv6Address)
	result += fmt.Sprintf("\tIPv6Connected: %v\n", i.IPv6Connected)
	result += fmt.Sprintf("\tMACAddress: %v\n", i.MACAddress)
	result += fmt.Sprintf("\tConfiguration: %v\n", i.Configuration)
	//result += fmt.Sprintf("\tShowInterface: %v\n", i.ShowInterface)
	result += fmt.Sprintf("\tInputPackets: %v\n", i.InputPackets)
	result += fmt.Sprintf("\tInputBytes: %v\n", i.InputBytes)
	result += fmt.Sprintf("\tInputDropped: %v\n", i.InputDropped)
	result += fmt.Sprintf("\tOutputPackets: %v\n", i.OutputPackets)
	result += fmt.Sprintf("\tOutputBytes: %v\n", i.OutputBytes)
	result += fmt.Sprintf("\tOutputDropped: %v\n", i.OutputDropped)
	result += fmt.Sprintf("==========================================")
	return result
}

var InterfaceBandwidthMap = map[string]int64{
	"10m":  10000,
	"100M": 100000,
	"1g":   1000000,
	"10g":  10000000,
	"100g": 100000000,
}

var InterfaceDB map[string]*Interface
var MatchAllInterface = regexp.MustCompile(`(?P<interface>interface[[:space:][:word:]/\-_\.:]+)\!`)
var MatchInterfaceLoopback = regexp.MustCompile(`Hardware is Loopback`)
var MatchInterfaceMACAddress = regexp.MustCompile(`Current HW addr:[[:space:]]+(?P<mac>[[:word:]]{2}\:[[:word:]]{2}\:[[:word:]]{2}\:[[:word:]]{2}\:[[:word:]]{2}\:[[:word:]]{2})`)
var MatchInterfaceIPAddress = regexp.MustCompile(`inet[[:space:]]+(?P<ip>[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+/[0-9]+)[[:space:]]+`)
var MatchInterfaceIPv6Address = regexp.MustCompile(`inet6[[:space:]]+(?P<ipv6>[[:word:]\:]+/[0-9]+)[[:space:]]+`)

//var MatchInterfaceIndexMetricMTU = regexp.MustCompile(`index (?P<index>[0-9]+) metric (?P<metric>[0-9]+) (?P<mtu>[0-9]+)`)
var MatchInterfaceIndexMetricMTU = regexp.MustCompile(`index[[:space:]]+(?P<index>[0-9]+)[[:space:]]+metric[[:space:]]+(?P<metric>[0-9]+)[[:space:]]+mtu[[:space:]]+(?P<mtu>[0-9]+)[[:space:]]+`)
var MatchInterfaceBandwidth = regexp.MustCompile(`Bandwidth[[:space:]]+(?P<bandwidth>[[:word:]]+)`)
var MatchInterfaceInputStatistics = regexp.MustCompile(`input packets (?P<packets>[0-9,]+), bytes (?P<bytes>[0-9,]+), dropped (?P<bytes>[0-9,]+)`)
var MatchInterfaceOutputStatistics = regexp.MustCompile(`output packets (?P<packets>[0-9,]+), bytes (?P<bytes>[0-9,]+), dropped (?P<bytes>[0-9,]+)`)
var MatchInterfaceStatus = regexp.MustCompile(`\<(?P<state>[A-Z,]+)\>`)

func main() {
	c := dut.New("V8500_SFU").Cli
	c.RunCommand(cli.Command{"enable", "configure terminal", "#"})
	ticker := time.NewTicker(time.Second / 1)
	for net := 0; net <= 2000000; net++ {
		<-ticker.C
		c.RunCommand(cli.Command{"config", "ip route 123.1.1.0/24 20.1.1.254", "#"})
		<-ticker.C
		c.RunCommand(cli.Command{"config", "no ip route 123.1.1.0/24 20.1.1.254", "#"})
	}
	//	c.RunCommand(cli.Command{"enable", "show running-config", "#"})

	log.Println(c.CommandResult())
}

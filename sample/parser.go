package main

import (
	"cli"
	"dut"
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
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
	c.RunCommand(cli.Command{"enable", "show running-config", "#"})
	log.Println(c.CommandResult())
	interfaces := MatchAllInterface.FindAllStringSubmatch(c.CommandResult(), -1)
	InterfaceDB = make(map[string]*Interface, len(interfaces))
	for _, i := range interfaces {
		//log.Println(string(i[1]))
		// Show interface
		config := strings.Split(string(i[1]), "\r")
		ifname := config[0]
		isL3 := false
		isVlan := false
		if strings.Contains(ifname, "interface vlan") || strings.Contains(ifname, "interface br") || strings.Contains(ifname, "interface lo") {
			isL3 = true
			if strings.Contains(ifname, "interface vlan") || strings.Contains(ifname, "interface br") {
				isVlan = true
			}
		}

		c.RunCommand(cli.Command{"enable", "show " + ifname, "#"})
		var newInterface = Interface{
			Name:               ifname,
			IsLoopback:         false,
			IsL3:               isL3,
			IsVlan:             isVlan,
			IsRunning:          false,
			IsMulticastCapable: false,
			IsBroadcastCapable: false,
			Configuration:      config[1:],
			ShowInterface:      c.CommandResult(),
		}

		// Get Interface IPv4 address configuration
		log.Println(newInterface.ShowInterface)
		ips := MatchInterfaceIPAddress.FindAllStringSubmatch(newInterface.ShowInterface, -1)
		newInterface.IPAddress = make(map[string]*net.IP, len(ips))
		newInterface.IPConnected = make(map[string]*net.IPNet, len(ips))
		for _, ip := range ips {
			log.Println(ip[1])
			IP, IPNet, err := net.ParseCIDR(ip[1])
			if err != nil {
				log.Fatal("Invalid ip address: ", ip[1])
			}
			newInterface.IPAddress[ip[1]] = &IP
			newInterface.IPConnected[ip[1]] = IPNet
		}

		// Get Interface IPv6 address configuration
		ip6s := MatchInterfaceIPv6Address.FindAllStringSubmatch(newInterface.ShowInterface, -1)
		newInterface.IPv6Address = make(map[string]*net.IP, len(ip6s))
		newInterface.IPv6Connected = make(map[string]*net.IPNet, len(ip6s))
		for _, ip6 := range ip6s {
			log.Println(ip6[1])
			IP, IPNet, err := net.ParseCIDR(ip6[1])
			if err != nil {
				log.Fatal("Invalid ipv6 address: ", ip6[1])
			}
			newInterface.IPv6Address[ip6[1]] = &IP
			newInterface.IPv6Connected[ip6[1]] = IPNet
		}

		// Get interface ifindex/Metric/MTU
		imt := MatchInterfaceIndexMetricMTU.FindStringSubmatch(newInterface.ShowInterface)
		if imt == nil {
			log.Fatal("no index/metric/mtu for interface: ", newInterface.Name)
		}
		log.Println(imt[1], imt[2], imt[3])
		ifindex, err := strconv.ParseInt(imt[1], 10, 64)
		if err != nil {
			log.Fatal("Invalid ifindex: ", imt[1])
		}
		newInterface.IfIndex = ifindex

		ifmetric, err := strconv.ParseInt(imt[2], 10, 64)
		if err != nil {
			log.Fatal("Invalid ifmetric: ", imt[2])
		}
		newInterface.IfMetric = ifmetric

		ifmtu, err := strconv.ParseInt(imt[3], 10, 64)
		if err != nil {
			log.Fatal("Invalid ifmtu: ", imt[3])
		}
		newInterface.IfMTU = ifmtu

		// Get interace Bandwitdh
		bands := MatchInterfaceBandwidth.FindStringSubmatch(newInterface.ShowInterface)
		if bands == nil {
			log.Fatal("Invalid bandwidth: ", bands[1])
		}
		if b, ok := InterfaceBandwidthMap[bands[1]]; !ok {
			log.Fatal("Unknown bandwidth: ", bands[1])
		} else {
			newInterface.Bandwidth = b
		}

		// Get interface input statistics
		inputs := MatchInterfaceInputStatistics.FindStringSubmatch(newInterface.ShowInterface)
		if inputs == nil {
			log.Fatal("No input statistics for interface: ", ifname)
		}

		inputPackets, err := strconv.ParseInt(strings.Replace(inputs[1], ",", "", -1), 10, 64)
		if err != nil {
			log.Fatal("Parse input packets error")
		}
		newInterface.InputPackets = inputPackets

		inputBytes, err := strconv.ParseInt(strings.Replace(inputs[2], ",", "", -1), 10, 64)
		if err != nil {
			log.Fatal("Parse input bytes error")
		}
		newInterface.InputBytes = inputBytes

		inputDropped, err := strconv.ParseInt(strings.Replace(inputs[2], ",", "", -1), 10, 64)
		if err != nil {
			log.Fatal("Parse input Dropped error")
		}
		newInterface.InputDropped = inputDropped

		// Get interface output statistics
		outputs := MatchInterfaceInputStatistics.FindStringSubmatch(newInterface.ShowInterface)
		if outputs == nil {
			log.Fatal("No output statistics for interface: ", ifname)
		}

		outputPackets, err := strconv.ParseInt(strings.Replace(outputs[1], ",", "", -1), 10, 64)
		if err != nil {
			log.Fatal("Parse output packets error")
		}
		newInterface.InputPackets = outputPackets

		outputBytes, err := strconv.ParseInt(strings.Replace(outputs[2], ",", "", -1), 10, 64)
		if err != nil {
			log.Fatal("Parse output bytes error")
		}
		newInterface.InputBytes = outputBytes

		outputDropped, err := strconv.ParseInt(strings.Replace(outputs[3], ",", "", -1), 10, 64)
		if err != nil {
			log.Fatal("Parse output Dropped error")
		}
		newInterface.InputDropped = outputDropped

		// Get interface Status
		status := MatchInterfaceStatus.FindStringSubmatch(newInterface.ShowInterface)
		if status == nil {
			log.Fatal("Cannot find the interface status of: ", ifname)
		}
		log.Println(status)
		if strings.Contains(status[1], "UP") {
			newInterface.IsAdminUP = true
		} else {
			newInterface.IsAdminUP = false
		}

		if strings.Contains(status[1], "RUNNING") {
			newInterface.IsRunning = true
		} else {
			newInterface.IsRunning = false
		}

		if strings.Contains(status[1], "BROADCAST") {
			newInterface.IsBroadcastCapable = true
		} else {
			newInterface.IsBroadcastCapable = false
		}

		if strings.Contains(status[1], "MULTICAST") {
			newInterface.IsMulticastCapable = true
		} else {
			newInterface.IsMulticastCapable = false
		}

		// Check if interface is Loopback interface
		loopback := MatchInterfaceLoopback.FindStringSubmatch(newInterface.ShowInterface)
		if loopback != nil {
			newInterface.IsLoopback = true
		}

		// Get Interface MAC Address
		if newInterface.IsLoopback == false {
			macs := MatchInterfaceMACAddress.FindStringSubmatch(newInterface.ShowInterface)
			mac, err := net.ParseMAC(macs[1])
			if err != nil {
				log.Fatal("Invalid MAC address: ", macs[1])
			}
			newInterface.MACAddress = map[string]*net.HardwareAddr{
				macs[1]: &mac,
			}
		}

		InterfaceDB[ifname] = &newInterface
	}

	for _, intf := range InterfaceDB {
		fmt.Println(intf)
	}
}

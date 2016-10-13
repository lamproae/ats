package analyzer

import (
	//	"log"
	//"regexp"
)

var delimeter = '!'

//key->command : value->mode
//in running-config file the basic mode is config mode
var ModeSwitchCommand = map[string]string {
	"enable" : "enable",
	"config  terminal" : "config",
	"bridge" : "bridge",
	"interface" : "interface",
	"router ospf" : "ospf",
	"router bgp" : "bgp",
	"address-family ipv4" : "address family ipv4",
	"address-family ipv6" : "address family ipv6",
	"router-map" : "route-map",
	"ip dhcp pool" : "ip dhcp pool",
	"ip dhcp option" : "ip dhcp option",
	"ip dhcp class" : "ip dhcp class",
	"ipv6 dhcp pool" : "ipv6 dhcp pool",
	"ipv6 dhcp option" : "ipv6 dhcp option",
	"ipv6 dhcp class" : "ipv6 dhcp class",
	"access-list" : "access-list",
	"policy" : "policy",
	"policy-map" : "policy-map",
	"policer" : "policer",
	"epon-profile" : "epon-profile",
	"eonu-profile" : "eonu-profile",
}

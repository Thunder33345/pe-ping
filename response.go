package ping

import "time"

//Response is the parsed response
//raw format for reference:
//MCPE;server name;protocol;version;playercount;maxplayers;server ID;LAN text;gamemode name;gamemode ID;ipv4 port;ipv6 port
type Response struct {
	//Latency is how long it took for the ping to be responded, injected: not part of the protocol
	Latency time.Duration
	//ResponseLength is how many elements are responded by the server, injected: not part of the protocol
	ResponseLength int
	//Bellow is parsed response
	Name         string
	Protocol     string //probably int
	Version      string
	Players      string //probably int
	MaxPlayers   string //probably int
	ServerID     string //probably int
	LanText      string
	GameModeName string
	GameModeId   string //probably int
	PortV4       string //probably int
	PortV6       string //probably int
	//Extra is where the remainder of the splits goes in, if any
	Extra string
}

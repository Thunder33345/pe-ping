package ping

import (
	"github.com/sandertv/go-raknet"
	"strings"
	"time"
)

func readSlice(s []string, i int) string {
	if i < 0 {
		return ""
	}
	if len(s)-1 >= i {
		return s[i]
	}
	return ""
}

func ParseResponse(b []byte, req DataRequirement) (Response, error) {
	s := string(b)
	if !strings.HasPrefix(s, prefix) {
		return Response{}, errInvalidPrefixRespond
	}
	s = s[len(prefix):]
	sp := strings.SplitN(s, ";", splitCount)

	resp := Response{}

	if len(sp) < int(req) {
		return Response{}, errInsufficientDataResponse
	}
	resp.ResponseLength = len(sp)
	if sp[len(sp)-1] == "" {
		resp.ResponseLength--
	}

	resp.Name = sp[0]
	resp.Protocol = readSlice(sp, 1)
	resp.Version = readSlice(sp, 2)
	resp.Players = readSlice(sp, 3)
	resp.MaxPlayers = readSlice(sp, 4)
	resp.ServerID = readSlice(sp, 5)
	resp.LanText = readSlice(sp, 6)
	resp.GameModeName = readSlice(sp, 7)
	resp.GameModeID = readSlice(sp, 8)
	resp.PortV4 = readSlice(sp, 9)
	resp.PortV6 = readSlice(sp, 10)
	resp.Extra = readSlice(sp, 11)

	return resp, nil
}

func RawPing(addr string, dialer raknet.Dialer, timeout time.Duration) ([]byte, time.Duration, error) {
	start := time.Now()
	b, e := dialer.PingTimeout(addr, timeout)
	if e != nil {
		return []byte{}, 0, e
	}
	dur := time.Now().Sub(start)
	return b, dur, nil
}

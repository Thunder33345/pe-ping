package ping

import (
	"errors"
	"github.com/sandertv/go-raknet"
	"time"
)

//prefix of the raknet ping protocol, if prefix isn't present, the response will be deemed invalid.
const prefix = "MCPE;"

//splitCount is how many times should response be split into, it should be 11(standard)+1(extra)
const splitCount = 12

var (
	errInvalidPrefixRespond     = errors.New("query responded with invalid prefix")
	errInsufficientDataResponse = errors.New("query responded with insufficient data")
)

type DataRequirement int

const (
	//ReqNone requires nothing to be present
	ReqNone DataRequirement = 0
	//ReqPlayers requires server to respond with at least player count
	ReqPlayers DataRequirement = 4
	//ReqAll requires all data to be present
	//This shouldn't be used as it's unlikely to be implemented by most servers
	ReqAll DataRequirement = 11
	//ReqExtra requires extra data to be present, which normally shouldn't exist
	ReqExtra DataRequirement = 12
)

type Raknet struct {
	//Dialer the dialer that will be used
	Dialer raknet.Dialer
	//TimeOut the maximum timeout to wait
	TimeOut time.Duration
	//Requirement is the minimum amount of data that is required, or else error is returned
	//otherwise ignored if it's <=0
	Requirement DataRequirement
}

func NewRaknet(dialer raknet.Dialer, timeout time.Duration, requirement DataRequirement) *Raknet {
	return &Raknet{
		Dialer:      dialer,
		TimeOut:     timeout,
		Requirement: requirement,
	}
}

func (r *Raknet) Ping(addr string) (Response, error) {
	b, d, e := RawPing(addr, r.Dialer, r.TimeOut)
	if e != nil {
		return Response{}, e
	}

	resp, e := ParseResponse(b, r.Requirement)
	if e != nil {
		return Response{}, e
	}
	resp.Latency = d

	return resp, nil
}
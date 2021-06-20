package ping

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
)

func Test_readSlice(t *testing.T) {
	ss := []string{"1", "2", "3", "4"}
	set := []struct {
		key  int
		want string
	}{
		{key: -1, want: ""},
		{key: 0, want: "1"},
		{key: 2, want: "3"},
		{key: 3, want: "4"},
		{key: 10, want: ""},
	}
	for i, s := range set {
		t.Run(fmt.Sprintf("Test #%d", i), func(t *testing.T) {
			a := assert.New(t)
			result := readSlice(ss, s.key)
			a.Equal(s.want, result)
		})
	}
}

func TestParseResponse_prefix(t *testing.T) {
	sets := []struct {
		str     string
		success bool
		name    string
	}{
		{str: "abc;def", success: false},
		{str: "MCPE;abc", success: true, name: "abc"},
		{str: "MCPE;", success: true, name: ""},
		{str: "MC;hello", success: false},
		{str: "MC;", success: false},
		{str: "MC", success: false},
		{str: "MCPE", success: false},
	}
	for i, set := range sets {
		t.Run(fmt.Sprintf("Test #%d", i), func(t *testing.T) {
			a := assert.New(t)
			res, err := ParseResponse([]byte(set.str), ReqNone)
			if set.success {
				a.Nil(err)
				a.NotNil(res)
				a.Equal(set.name, res.Name)
			} else {
				a.Empty(res)
				a.Equal(errInvalidPrefixRespond, err)
			}
		})
	}
}

func TestParseResponse_requirement(t *testing.T) {
	sets := []struct {
		//in
		str string
		req int
		//out
		success bool
		name    string
		rl      int
	}{
		{str: "MCPE;abc", req: -1, success: true, name: "abc", rl: 1},
		{str: "MCPE;abc", req: 1, success: true, name: "abc", rl: 1},
		{str: "MCPE;abc", req: 2, success: false},
		{str: "MCPE;abc", req: 3, success: false},
		{str: "MCPE;acd;def;ghi", req: 1, success: true, name: "acd", rl: 3},
		{str: "MCPE;acd;def;ghi;", req: 2, success: true, name: "acd", rl: 3},
		{str: "MCPE;acd;def;ghi;;", req: 3, success: true, name: "acd", rl: 4},
		{str: "MCPE;acd;def;ghi", req: 5, success: false, name: "acd"},
		{str: "MCPE;acd;def;ghi", req: 6, success: false, name: "acd"},
	}
	for i, set := range sets {
		t.Run(fmt.Sprintf("Test #%d", i), func(t *testing.T) {
			a := assert.New(t)
			res, err := ParseResponse([]byte(set.str), DataRequirement(set.req))
			if set.success {
				a.Nil(err)
				a.NotNil(res)
				a.Equal(set.name, res.Name)
				a.Equal(set.rl, res.ResponseLength)
			} else {
				a.Empty(res)
				a.Equal(err, errInsufficientDataResponse)
			}
		})
	}
}

func TestParseResponse_responseLength(t *testing.T) {
	sets := []struct {
		//in
		str string
		//out
		rl int
	}{
		{str: "MCPE;", rl: 0},
		{str: "MCPE;;", rl: 1},
		{str: "MCPE;abc", rl: 1},
		{str: "MCPE;abc;", rl: 1},
		{str: "MCPE;abc;;", rl: 2},
		{str: "MCPE;abc;;;", rl: 3},
		{str: "MCPE;abc;def", rl: 2},
		{str: "MCPE;abc;def;", rl: 2},
		{str: "MCPE;abc;def;;", rl: 3},
	}
	for i, set := range sets {
		t.Run(fmt.Sprintf("Test #%d", i), func(t *testing.T) {
			a := assert.New(t)
			res, err := ParseResponse([]byte(set.str), ReqNone)
			a.Nil(err)
			a.Equal(set.rl, res.ResponseLength)
		})
	}
}
func TestParseResponse_mock(t *testing.T) {
	for i := 1; i < 50; i++ {
		t.Run(fmt.Sprintf("Mock Test %0d", i), func(t *testing.T) {
			a := assert.New(t)
			in, out := mockResponseAndInput()
			resp, err := ParseResponse([]byte(in), ReqNone)
			a.Nil(err)
			a.Equal(out, resp)
		})
	}
}

func BenchmarkParse(b *testing.B) {
	in, _ := mockResponseAndInput()
	for i := 0; i < b.N; i++ {
		_, _ = ParseResponse([]byte(in), ReqNone)
	}
}

func mockResponseAndInput() (input string, output Response) {
	input = "MCPE;"
	input, output.Name = addRandString(input)
	input, output.Protocol = addRandString(input)
	input, output.Version = addRandString(input)
	input, output.Players = addRandString(input)
	input, output.MaxPlayers = addRandString(input)
	input, output.ServerID = addRandString(input)
	input, output.LanText = addRandString(input)
	input, output.GameModeName = addRandString(input)
	input, output.GameModeID = addRandString(input)
	input, output.PortV4 = addRandString(input)
	input, output.PortV6 = addRandString(input)
	e := randString(10) + ";" + randString(10) + ";"
	input, output.Extra = input+e, e
	output.ResponseLength = 12
	return
}

func addRandString(org string) (string, string) {
	switch rand.Intn(5) {
	case 1:
		v := ""
		return org + v + ";", v
	case 2:
		v := strconv.Itoa(rand.Intn(10000))
		return org + v + ";", v
	}

	v := randString(rand.Intn(8))
	return org + v + ";", v
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

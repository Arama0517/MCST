package lib

import (
	"github.com/ybbus/jsonrpc"
)

func KillAria2c() {
	err := aria2Process.Kill()
	if err != nil {
		panic(err)
	}
}

func GetAria2cClient() jsonrpc.RPCClient {
	return jsonrpc.NewClient("ws://localhost:6800/jsonrpc")
}

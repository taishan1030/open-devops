package rpc

import "fmt"

func (*Server) Ping(input string, output *string) error {
	fmt.Println(input)
	*output = "receive data"
	return nil
}

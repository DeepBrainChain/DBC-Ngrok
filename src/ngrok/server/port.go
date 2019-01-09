package server

import (
	"errors"
	"fmt"
)
var portMode int
var portPool []int

const (
	_ = iota
	PortRawMode
	PortServerMode
)

func init() {
	portMode = PortServerMode

	portPool = make([]int, 0)

	for i:=20000; i<20100; i++ {
		portPool = append(portPool, i)
	}
}

func AllocPort() (port int, err error){

	err = nil
	if len(portPool) == 0 {
		return -1, errors.New("no port available")
	}

	port = portPool[0]
	portPool=portPool[1:]

	return
}

func RemovePort(port int) error {
	for i, v := range portPool {
		if v == port {
			//  port
			portPool= append(portPool[:i], portPool[i+1:]...)
			return nil
		}
	}

	return errors.New(fmt.Sprintf("no found port %d from portPool", port))

}

func FreePort(port int) {

	for _, v := range portPool {
		if v == port {
			// duplicate port
			return
		}
	}

	portPool = append(portPool, port)
}

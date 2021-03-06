package client

import "github.com/pebbe/zmq4"

// LogSign Used for logging client messages
const LogSign string = "Client"

// client A struct to represent the client structure
type client struct {
	id           int
	ip           string
	port         string
	notifyPort   string
	trackerIP    string
	trackerPorts []string
	socket       *zmq4.Socket
}

// NewClient A constructor function for the client type
func NewClient(_id int, _ip string, _port string, _notifyPort string, _trackerIP string, _trackerPorts []string) client {
	clientObj := client{
		id:           _id,
		ip:           _ip,
		port:         _port,
		notifyPort:   _notifyPort,
		trackerIP:    _trackerIP,
		trackerPorts: _trackerPorts,
	}

	return clientObj
}

package main

import (
	"github.com/75912001/goz/ztcp"
)

//User
type User struct {
	ztcp.PeerConn
	UID UserID
}

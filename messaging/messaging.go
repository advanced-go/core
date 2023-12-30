package messaging

import (
	"github.com/advanced-go/core/runtime"
)

const (
	addLoc      = PkgPath + ":Add"
	sendCtrlLoc = PkgPath + ":SendCtrl"
	sendDataLoc = PkgPath + ":SendData"
)

var HostExchange = NewExchange()

// Add - add a mailbox
func Add(m *Mailbox) runtime.Status {
	status := HostExchange.Add(m)
	if !status.OK() {
		status.AddLocation(addLoc)
	}
	return status
}

// SendCtrl - send to command channel
func SendCtrl(msg Message) runtime.Status {
	status := HostExchange.SendCtrl(msg)
	if !status.OK() {
		status.AddLocation(sendCtrlLoc)
	}
	return status
}

// SendData - send to data channel
func SendData(msg Message) runtime.Status {
	status := HostExchange.SendData(msg)
	if !status.OK() {
		status.AddLocation(sendDataLoc)
	}
	return status
}

// Shutdown - send a shutdown message to all directory entries
//func shutdown() {
//HostDirectory.Shutdown()
//}

func shutdownHost(msg Message) runtime.Status {
	//TO DO: authentication and implementation

	return runtime.StatusOK()
}

package host

import (
	"github.com/advanced-go/core/messaging"
	"github.com/advanced-go/core/runtime"
)

func shutdownHost(msg messaging.Message) runtime.Status {
	//TO DO: authentication and implementation
	return runtime.StatusOK()
}

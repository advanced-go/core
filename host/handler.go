package host

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/access"
	"github.com/advanced-go/core/controller"
	"net/http"
)

var (
	httpHandlerProxy = NewProxy()
	hostCtrl         *controller.Controller
	authHandler      HttpHandlerFunc //func(w http.ResponseWriter,r *http.Request)
	okFunc           = func(code int) bool { return code == http.StatusOK }
)

func SetHostController(ctrl *controller.Controller) {
	if ctrl != nil {
		hostCtrl = ctrl
		if ctrl.RouteName == "" {
			hostCtrl.RouteName = controller.HostRouteName
		}
	}
}

func SetAuthHandler(h func(w http.ResponseWriter, r *http.Request), ok func(int) bool) {
	if h != nil {
		authHandler = h
		if ok != nil {
			okFunc = ok
		}
	}
}

// RegisterHandler - add a path and Http handler to the proxy
// TO DO : panic on duplicate handler and pattern combination
func RegisterHandler(path string, handler HttpHandlerFunc) error {
	if len(path) == 0 {
		return errors.New("error: path is empty")
	}
	if handler == nil {
		return errors.New(fmt.Sprintf("error: handler for path %v is nil", path))
	}
	h := handler
	if authHandler != nil {
		h = NewConditionalIntermediary(authHandler, handler, okFunc)
	}
	if hostCtrl != nil {
		h = newIngressControllerIntermediary(hostCtrl, h, access.IngressTraffic)
	}
	err := httpHandlerProxy.Register(path, h)
	return err
}

// HttpHandler - handler for messaging
func HttpHandler(w http.ResponseWriter, r *http.Request) {
	if r == nil || w == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	nid, _, ok := UprootUrn(r.URL.Path)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	handler := httpHandlerProxy.LookupByNID(nid)
	if handler == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	handler(w, r)
}

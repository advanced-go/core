package host

import (
	"fmt"
	"github.com/advanced-go/core/controller"
	"net/http"
)

func NewEgressControllerIntermediary(ctrl *controller.Controller) HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if ctrl == nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error: controller is nil")
			return
		}
	}
}

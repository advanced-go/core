package resiliency

import (
	"context"
	"errors"
	"net/http"
)

type controllerWrapper struct {
	rt http.RoundTripper
}

// RoundTrip - implementation of the RoundTrip interface for a transport, also logs an access entry
func (w *controllerWrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	//var start = time.Now().UTC()
	var resp *http.Response
	var err error
	//var statusFlags = ""

	// !panic
	if w == nil || w.rt == nil {
		return nil, errors.New("invalid handler round tripper configuration : http.RoundTripper is nil")
	}
	/*
		ctrl, _ := controller.EgressLookup(req)
		if ctrl == nil {
			resp, err, statusFlags = w.exchange(nil, req)
			controller.LogHttpEgress(ctrl, start, time.Since(start), req, resp, statusFlags)
			return resp, err
		}
		if rlc := ctrl.RateLimiter(); rlc.IsEnabled() && !rlc.Allow() {
			resp = &http.Response{Request: req, StatusCode: rlc.StatusCode()}
			controller.LogHttpEgress(ctrl, start, time.Since(start), req, resp, controller.RateLimitFlag)
			return resp, nil
		}
		ctrl.UpdateHeaders(req)
		if pc := ctrl.Proxy(); pc.IsEnabled() && len(pc.Uri()) > 0 {
			req.URL = pc.BuildUri(req.URL)
			if req.URL != nil {
				req.Host = req.URL.Host
			}
			for _, header := range pc.Headers() {
				req.Header.Add(header.Name, header.Value)
			}
		}
		resp, err, statusFlags = w.exchange(ctrl.Timeout(), req)
		if err != nil {
			return resp, err
		}
		controller.LogHttpEgress(ctrl, start, time.Since(start), req, resp, statusFlags)


	*/
	return resp, err
}

/*
func (w *controllerWrapper) exchange(tc Timeout, req *http.Request) (resp *http.Response, err error, statusFlags string) {
	if !tc.Duration == 0 {
		resp, err = w.rt.RoundTrip(req)
		return
	}
	ctx, cancel := context.WithTimeout(req.Context(), tc.Duration)
	//defer cancel()
	req = req.Clone(ctx)
	resp, err = w.rt.RoundTrip(req)
	if w.deadlineExceeded(err) {
		resp = &http.Response{Request: req, StatusCode: tc.StatusCode}
		err = nil
		statusFlags = UpstreamTimeoutFlag
		cancel()
	}
	return
}


*/

func (w *controllerWrapper) deadlineExceeded(err error) bool {
	return err != nil && err == context.DeadlineExceeded //errors.As(err, de)
}

// ControllerWrapTransport - provides a RoundTrip wrapper that applies egress controllers
func ControllerWrapTransport(client *http.Client) {
	if client == nil || client == http.DefaultClient {
		if http.DefaultClient.Transport == nil {
			http.DefaultClient.Transport = &controllerWrapper{rt: http.DefaultTransport}
		} else {
			http.DefaultClient.Transport = ControllerWrapRoundTripper(http.DefaultClient.Transport)
		}
	} else {
		if client.Transport == nil {
			client.Transport = &controllerWrapper{rt: http.DefaultTransport}
		} else {
			client.Transport = ControllerWrapRoundTripper(client.Transport)
		}
	}
}

// ControllerWrapRoundTripper - provides a RoundTrip wrapper that applies egress controllers
func ControllerWrapRoundTripper(rt http.RoundTripper) http.RoundTripper {
	return &controllerWrapper{rt: rt}
}

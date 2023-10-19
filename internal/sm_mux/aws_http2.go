package sm_mux

import (
	"context"
	"errors"
	log "log/slog"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/gorilla/mux"
	"github.com/jackmcguire1/riot-global-rankings/internal/utils"
)

type GorillaMuxAdapterV2 struct {
	RequestAccessorV2 core.RequestAccessorV2
	RequestAccessor   core.RequestAccessor
	router            *mux.Router
}

func NewV2(router *mux.Router) *GorillaMuxAdapterV2 {
	return &GorillaMuxAdapterV2{
		router: router,
	}
}

// Proxy receives an API Gateway proxy event, transforms it into an http.Request
// object, and sends it to the mux.Router for routing.
// It returns a proxy response object generated from the http.ResponseWriter.
func (h *GorillaMuxAdapterV2) Proxy(event interface{}) (interface{}, error) {
	log.With("event", utils.ToJSON(event)).Info("lambda invoked")

	x := &core.SwitchableAPIGatewayRequest{}
	x.UnmarshalJSON(utils.ToJsonBytes(event))

	if x.Version2() != nil {
		req, err := h.RequestAccessorV2.EventToRequest(*x.Version2())
		return h.proxyInternalV2(req, err)
	}
	if x.Version1() != nil {
		req, err := h.RequestAccessor.EventToRequest(*x.Version1())
		return h.proxyInternal(req, err)
	}

	return nil, core.NewLoggedError("Could not convert proxy event to request: %v", errors.New("Unable to determine version "))
}

// ProxyWithContext receives context and an API Gateway proxy event,
// transforms them into an http.Request object, and sends it to the mux.Router for routing.
// It returns a proxy response object generated from the http.ResponseWriter.
func (h *GorillaMuxAdapterV2) ProxyWithContext(ctx context.Context, event interface{}) (interface{}, error) {
	x := &core.SwitchableAPIGatewayRequest{}
	x.UnmarshalJSON(utils.ToJsonBytes(event))

	log.With("event", utils.ToJSON(event)).Debug("lambda invoked")
	if x.Version2() != nil {
		req, err := h.RequestAccessorV2.EventToRequestWithContext(ctx, *x.Version2())
		return h.proxyInternalV2(req, err)
	}

	return nil, core.NewLoggedError("Could not convert proxy event to request: %v", errors.New("Unable to determine version "))
}

func (h *GorillaMuxAdapterV2) proxyInternalV2(req *http.Request, err error) (*core.SwitchableAPIGatewayResponse, error) {
	if err != nil {
		timeout := core.GatewayTimeoutV2()

		return core.NewSwitchableAPIGatewayResponseV2(&timeout), core.NewLoggedError("Could not convert proxy event to request: %v", err)
	}

	w := core.NewProxyResponseWriterV2()

	var resp events.APIGatewayV2HTTPResponse
	if req.Method == http.MethodOptions {
		resp = events.APIGatewayV2HTTPResponse{
			Body:       "{}",
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Methods":     "OPTIONS,GET, POST, PUT, DELETE",
				"Access-Control-Allow-Headers":     "Content-Type,Authorization,X-Requested-With,Origin,Accept",
				"Access-Control-Allow-Credentials": "true",
				"Content-Type":                     "application/json",
			},
		}

		goto result
	}

	h.router.ServeHTTP(http.ResponseWriter(w), req)

	resp, err = w.GetProxyResponse()
	if err != nil {
		timeout := core.GatewayTimeoutV2()
		return core.NewSwitchableAPIGatewayResponseV2(&timeout), core.NewLoggedError("Error while generating proxy response: %v", err)
	}

result:
	log.
		With("response", utils.ToJSON(resp)).
		Debug("HTTP RESPONSE")

	return core.NewSwitchableAPIGatewayResponseV2(&resp), nil
}

func (h *GorillaMuxAdapterV2) proxyInternal(req *http.Request, err error) (*core.SwitchableAPIGatewayResponse, error) {
	if err != nil {
		timeout := core.GatewayTimeout()
		return core.NewSwitchableAPIGatewayResponseV1(&timeout), core.NewLoggedError("Could not convert proxy event to request: %v", err)
	}

	if req.Method == http.MethodOptions {
		resp := events.APIGatewayProxyResponse{
			Body:       "{}",
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Methods":     "OPTIONS,GET, POST, PUT, DELETE",
				"Access-Control-Allow-Headers":     "Content-Type,Authorization,X-Requested-With,Origin,Accept",
				"Access-Control-Allow-Credentials": "true",
				"Content-Type":                     "application/json",
			},
		}

		return core.NewSwitchableAPIGatewayResponseV1(&resp), nil
	}

	w := core.NewProxyResponseWriter()
	h.router.ServeHTTP(http.ResponseWriter(w), req)

	resp, err := w.GetProxyResponse()
	if err != nil {
		timeout := core.GatewayTimeout()
		return core.NewSwitchableAPIGatewayResponseV1(&timeout), core.NewLoggedError("Error while generating proxy response: %v", err)
	}

	return core.NewSwitchableAPIGatewayResponseV1(&resp), nil
}

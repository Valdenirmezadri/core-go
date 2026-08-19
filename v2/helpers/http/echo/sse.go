package helperecho

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v4"
)

// WriteSSE writes one event in SSE wire format and flushes.
func (h *helper) WriteSSE(e echo.Context, ID uint, kind string, data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w := e.Response()
	_, err = fmt.Fprintf(w, "id: %d\nevent: %s\ndata: %s\n\n", ID, kind, b)
	if err != nil {
		return err
	}
	e.Response().Flush()
	return nil
}

// WriteSSEPong writes a minimal pong event without JSON marshaling.
func (h *helper) WriteSSEPong(e echo.Context) error {
	w := e.Response()
	_, err := fmt.Fprintf(w, "event: ka\n\n")
	if err != nil {
		return err
	}
	e.Response().Flush()
	return nil
}

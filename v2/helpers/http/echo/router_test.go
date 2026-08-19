package helperecho

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestMuxDispatchPorInstancia(t *testing.T) {
	mux := NewMux()

	// registros antes do Attach (tenants sobem antes do echo) ficam em buffer
	mux.For(1).GET("/x", func(c echo.Context) error { return c.String(200, "t1") })
	mux.For(2).GET("/x", func(c echo.Context) error { return c.String(200, "t2") })

	e := echo.New()
	current := uint(1)
	mux.Attach(e.Group(""), func(echo.Context) (uint, error) { return current, nil })

	// registro depois do Attach entra direto
	mux.For(1).GET("/y", func(c echo.Context) error { return c.String(200, "y1") })

	get := func(path string) (int, string) {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		return rec.Code, rec.Body.String()
	}

	if code, body := get("/x"); code != 200 || body != "t1" {
		t.Fatalf("tenant 1: got %d %q", code, body)
	}
	current = 2
	if code, body := get("/x"); code != 200 || body != "t2" {
		t.Fatalf("tenant 2: got %d %q", code, body)
	}
	// tenant 2 não registrou /y — 404
	if code, _ := get("/y"); code != http.StatusNotFound {
		t.Fatalf("tenant 2 /y: got %d, want 404", code)
	}
	current = 1
	if code, body := get("/y"); code != 200 || body != "y1" {
		t.Fatalf("tenant 1 /y: got %d %q", code, body)
	}

	// Remove descarta os handlers do tenant
	mux.Remove(1)
	if code, _ := get("/x"); code != http.StatusNotFound {
		t.Fatalf("após Remove: got %d, want 404", code)
	}
}

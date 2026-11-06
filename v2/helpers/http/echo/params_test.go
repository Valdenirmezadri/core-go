package helperecho

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"testing"

	"github.com/Valdenirmezadri/core-go/v2/context/corectx"
	"github.com/Valdenirmezadri/core-go/v2/i18n"
	"github.com/labstack/echo/v4"
)

type nopLogger struct{}

func (nopLogger) Warningf(string, ...interface{}) {}

type user struct{ lang string }

func (u user) User() corectx.User { return corectx.User{Lang: u.lang} }

func newContext(method, target, body string) echo.Context {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	return echo.New().NewContext(req, httptest.NewRecorder())
}

func newHelper(t *testing.T) Helper {
	t.Helper()
	h, err := New(func() i18n.Logger { return nopLogger{} })
	if err != nil {
		t.Fatal(err)
	}
	return h
}

func TestParamToStrArr(t *testing.T) {
	h := newHelper(t)
	ctx := user{"pt-BR"}

	got, err := h.RequiredParamToStrArr(ctx, newContext(http.MethodPost, "/", "ids=a&ids=+b+&ids="), "ids")
	if err != nil || !slices.Equal(got, []string{"a", "b"}) {
		t.Fatalf("form: got %v, %v", got, err)
	}

	got, err = h.ReqParamStrArr(ctx, newContext(http.MethodDelete, "/?ids=a&ids=b", ""), "ids")
	if err != nil || !slices.Equal(got, []string{"a", "b"}) {
		t.Fatalf("query: got %v, %v", got, err)
	}

	_, err = h.RequiredParamToStrArr(ctx, newContext(http.MethodPost, "/", "ids="), "ids")
	if err == nil || err.Error() != "parâmetro 'ids' é obrigatório" {
		t.Fatalf("required sem valor: got %v", err)
	}

	got, err = h.ParamStrArr(ctx, newContext(http.MethodPost, "/", ""), "ids")
	if err != nil || got == nil || len(got) != 0 {
		t.Fatalf("sem required deve voltar [] vazio: got %#v, %v", got, err)
	}
}

func TestParamToUintArr(t *testing.T) {
	h := newHelper(t)
	ctx := user{"en"}

	got, err := h.ReqParamUintArr(ctx, newContext(http.MethodPost, "/", "ids=1&ids=0&ids=2"), "ids")
	if err != nil || !slices.Equal(got, []uint{1, 2}) {
		t.Fatalf("got %v, %v", got, err)
	}

	_, err = h.RequiredParamToUintArr(ctx, newContext(http.MethodPost, "/", "ids=0"), "ids")
	if err == nil || err.Error() != "parameter 'ids' is required" {
		t.Fatalf("required só com 0: got %v", err)
	}

	if _, err = h.ParamUintArr(ctx, newContext(http.MethodPost, "/", "ids=x"), "ids"); err == nil {
		t.Fatal("valor inválido deveria dar erro")
	}

	got, err = h.ParamToUintArr(ctx, newContext(http.MethodPost, "/", ""), "ids")
	if err != nil || got == nil || len(got) != 0 {
		t.Fatalf("sem required deve voltar [] vazio: got %#v, %v", got, err)
	}
}

func TestParamToArrStrArr(t *testing.T) {
	h := newHelper(t)
	ctx := user{"pt-BR"}
	body := url.Values{"groups": {`[["a","A b"],["c","D"]]`}}.Encode()

	got, err := h.ReqParamArrStrArr(ctx, newContext(http.MethodPost, "/", body), "groups")
	if err != nil || len(got) != 2 || !slices.Equal(got[1], []string{"c", "D"}) {
		t.Fatalf("got %v, %v", got, err)
	}

	_, err = h.RequiredParamToArrStrArr(ctx, newContext(http.MethodPost, "/", "groups=%5B%5D"), "groups")
	if err == nil || err.Error() != "parâmetro 'groups' é obrigatório" {
		t.Fatalf("required sem valor: got %v", err)
	}

	if _, err = h.RequiredParamToArrStrArr(ctx, newContext(http.MethodPost, "/", "groups=%5B%5Ba"), "groups"); err == nil {
		t.Fatal("json inválido deveria dar erro")
	}

	got, err = h.ParamArrStrArr(ctx, newContext(http.MethodPost, "/", ""), "groups")
	if err != nil || got == nil || len(got) != 0 {
		t.Fatalf("sem required deve voltar [] vazio: got %#v, %v", got, err)
	}
}

func TestParamToArrUintArr(t *testing.T) {
	h := newHelper(t)
	ctx := user{"pt-BR"}
	body := func(value string) string { return url.Values{"ids": {value}}.Encode() }

	got, err := h.ReqParamArrUintArr(ctx, newContext(http.MethodPost, "/", body("[[1,2],[3]]")), "ids")
	if err != nil || len(got) != 2 || !slices.Equal(got[0], []uint{1, 2}) || !slices.Equal(got[1], []uint{3}) {
		t.Fatalf("got %v, %v", got, err)
	}

	_, err = h.RequiredParamToArrUintArr(ctx, newContext(http.MethodPost, "/", body("[]")), "ids")
	if err == nil || err.Error() != "parâmetro 'ids' é obrigatório" {
		t.Fatalf("required sem valor: got %v", err)
	}

	for _, value := range []string{"[[-1]]", `[["a"]]`} {
		if _, err := h.RequiredParamToArrUintArr(ctx, newContext(http.MethodPost, "/", body(value)), "ids"); err == nil {
			t.Fatalf("%q deveria dar erro", value)
		}
	}

	got, err = h.ParamArrUintArr(ctx, newContext(http.MethodPost, "/", ""), "ids")
	if err != nil || got == nil || len(got) != 0 {
		t.Fatalf("sem required deve voltar [] vazio: got %#v, %v", got, err)
	}
}

func TestParamArrJSON(t *testing.T) {
	h := newHelper(t)
	ctx := user{"pt-BR"}
	body := func(value string) string { return url.Values{"ids": {value}}.Encode() }

	strs, err := h.RequiredParamToStrArr(ctx, newContext(http.MethodPost, "/", body("[\n\"a\",\n\" b \", \"\"\n]")), "ids")
	if err != nil || !slices.Equal(strs, []string{"a", "b"}) {
		t.Fatalf("got %v, %v", strs, err)
	}

	uints, err := h.RequiredParamToUintArr(ctx, newContext(http.MethodPost, "/", body("[1,2,0]")), "ids")
	if err != nil || !slices.Equal(uints, []uint{1, 2}) {
		t.Fatalf("got %v, %v", uints, err)
	}

	for _, value := range []string{`["a",]`, `[["a"]]`} {
		if _, err := h.RequiredParamToStrArr(ctx, newContext(http.MethodPost, "/", body(value)), "ids"); err == nil {
			t.Fatalf("%q deveria dar erro", value)
		}
	}
}

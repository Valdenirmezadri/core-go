package helperecho

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func context(method, target, body string) echo.Context {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	return echo.New().NewContext(req, httptest.NewRecorder())
}

func TestParamToStrArr(t *testing.T) {
	h := &helperEcho{}

	got, err := h.RequiredParamToStrArr(context(http.MethodPost, "/", "ids=a&ids=+b+&ids="), "ids")
	if err != nil || !slices.Equal(got, []string{"a", "b"}) {
		t.Fatalf("form: got %v, %v", got, err)
	}

	got, err = h.RequiredParamToStrArr(context(http.MethodDelete, "/?ids=a&ids=b", ""), "ids")
	if err != nil || !slices.Equal(got, []string{"a", "b"}) {
		t.Fatalf("query: got %v, %v", got, err)
	}

	if _, err = h.RequiredParamToStrArr(context(http.MethodPost, "/", "ids="), "ids"); err == nil {
		t.Fatal("required sem valor deveria dar erro")
	}

	got, err = h.ParamToStrArr(context(http.MethodPost, "/", ""), "ids")
	if err != nil || got == nil || len(got) != 0 {
		t.Fatalf("sem required deve voltar [] vazio: got %#v, %v", got, err)
	}
}

func TestParamToUintArr(t *testing.T) {
	h := &helperEcho{}

	got, err := h.RequiredParamToUintArr(context(http.MethodPost, "/", "ids=1&ids=0&ids=2"), "ids")
	if err != nil || !slices.Equal(got, []uint{1, 2}) {
		t.Fatalf("got %v, %v", got, err)
	}

	if _, err = h.RequiredParamToUintArr(context(http.MethodPost, "/", "ids=0"), "ids"); err == nil {
		t.Fatal("required só com 0 deveria dar erro")
	}

	if _, err = h.ParamToUintArr(context(http.MethodPost, "/", "ids=x"), "ids"); err == nil {
		t.Fatal("valor inválido deveria dar erro")
	}

	got, err = h.ParamToUintArr(context(http.MethodPost, "/", ""), "ids")
	if err != nil || got == nil || len(got) != 0 {
		t.Fatalf("sem required deve voltar [] vazio: got %#v, %v", got, err)
	}
}

func TestParamToArrStrArr(t *testing.T) {
	h := &helperEcho{}
	body := url.Values{"groups": {`[["a","A b"],["c","D"]]`}}.Encode()

	got, err := h.RequiredParamToArrStrArr(context(http.MethodPost, "/", body), "groups")
	if err != nil || len(got) != 2 || !slices.Equal(got[1], []string{"c", "D"}) {
		t.Fatalf("got %v, %v", got, err)
	}

	for name, value := range map[string]string{"vazio": "", "array vazio": "[]", "json inválido": "[[a]"} {
		if _, err := h.RequiredParamToArrStrArr(context(http.MethodPost, "/", url.Values{"groups": {value}}.Encode()), "groups"); err == nil {
			t.Fatalf("%s deveria dar erro", name)
		}
	}

	got, err = h.ParamToArrStrArr(context(http.MethodPost, "/", ""), "groups")
	if err != nil || got == nil || len(got) != 0 {
		t.Fatalf("sem required deve voltar [] vazio: got %#v, %v", got, err)
	}
}

func TestParamToArrUintArr(t *testing.T) {
	h := &helperEcho{}
	body := func(value string) string { return url.Values{"ids": {value}}.Encode() }

	got, err := h.RequiredParamToArrUintArr(context(http.MethodPost, "/", body("[[1,2],[3]]")), "ids")
	if err != nil || len(got) != 2 || !slices.Equal(got[0], []uint{1, 2}) || !slices.Equal(got[1], []uint{3}) {
		t.Fatalf("got %v, %v", got, err)
	}

	for _, value := range []string{"", "[]", "[[-1]]", `[["a"]]`} {
		if _, err := h.RequiredParamToArrUintArr(context(http.MethodPost, "/", body(value)), "ids"); err == nil {
			t.Fatalf("%q deveria dar erro", value)
		}
	}

	got, err = h.ParamToArrUintArr(context(http.MethodPost, "/", ""), "ids")
	if err != nil || got == nil || len(got) != 0 {
		t.Fatalf("sem required deve voltar [] vazio: got %#v, %v", got, err)
	}
}

func TestParamArrJSON(t *testing.T) {
	h := &helperEcho{}
	body := func(value string) string { return url.Values{"ids": {value}}.Encode() }

	strs, err := h.RequiredParamToStrArr(context(http.MethodPost, "/", body("[\n\"a\",\n\" b \", \"\"\n]")), "ids")
	if err != nil || !slices.Equal(strs, []string{"a", "b"}) {
		t.Fatalf("got %v, %v", strs, err)
	}

	uints, err := h.RequiredParamToUintArr(context(http.MethodPost, "/", body("[1,2,0]")), "ids")
	if err != nil || !slices.Equal(uints, []uint{1, 2}) {
		t.Fatalf("got %v, %v", uints, err)
	}

	for _, value := range []string{`["a",]`, `[["a"]]`} {
		if _, err := h.RequiredParamToStrArr(context(http.MethodPost, "/", body(value)), "ids"); err == nil {
			t.Fatalf("%q deveria dar erro", value)
		}
	}

	if _, err := h.RequiredParamToStrArr(context(http.MethodPost, "/", body("[]")), "ids"); err == nil {
		t.Fatal("array vazio deveria dar erro no required")
	}
}

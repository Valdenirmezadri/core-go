package helperecho

import (
	"mime/multipart"

	"github.com/Valdenirmezadri/core-go/context/corectx"
	"github.com/labstack/echo/v4"
)

// RequiredFormFile lê um arquivo do multipart; ausente ou form inválido vira
// erro traduzido (FormRequired).
func (h *helper) RequiredFormFile(ctx corectx.UserContext, c echo.Context, param string) (*multipart.FileHeader, error) {
	return h.formFile(ctx, c, param, true)
}

// FormFile lê um arquivo do multipart, se presente. Sem arquivo (ou request
// sem multipart) retorna nil sem erro.
func (h *helper) FormFile(ctx corectx.UserContext, c echo.Context, param string) (*multipart.FileHeader, error) {
	return h.formFile(ctx, c, param, false)
}

func (h *helper) formFile(ctx corectx.UserContext, c echo.Context, param string, required bool) (*multipart.FileHeader, error) {
	header, err := c.FormFile(param)
	if err != nil {
		if required {
			return nil, h.errT(ctx, "FormRequired", param)
		}

		return nil, nil
	}

	return header, nil
}

// Alias
func (h *helper) ReqFormFile(ctx corectx.UserContext, c echo.Context, param string) (*multipart.FileHeader, error) {
	return h.RequiredFormFile(ctx, c, param)
}

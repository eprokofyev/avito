package response

import (
	"avito/internal/pkg/models"
	"github.com/mailru/easyjson"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func Respond(ctx *routing.Context, status int, body map[string]interface{}) error {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	response := &models.Response{
		Status: status,
		Body:   body,
	}

	data, err := easyjson.Marshal(response)
	if err != nil {
		return err
	}

	_, err = ctx.Write(data)
	if err != nil {
		return err
	}

	return nil
}



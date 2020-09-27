package balance

import (
	"avito/internal/pkg/database"
	"avito/internal/pkg/errors"
	"avito/internal/pkg/models"
	"avito/internal/pkg/response"
	"avito/internal/pkg/utils"
	routing "github.com/qiangxue/fasthttp-routing"
	"log"
	"strconv"
	"strings"
)

type Handler struct {
	db database.IDatabase
}

func NewHandler(usecase database.IDatabase) *Handler {
	return &Handler{
		usecase,
	}
}

func (h *Handler) Transfer(ctx *routing.Context) error {
	input := &models.Transfer{}

	err := input.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		return errors.BadData.Wrap(err, "handler Transfer: unmarshal data error")
	}

	if input.SenderID == input.RecipientID {
		return errors.BadUserData.New("handler Transfer: user's ids are incorrect")
	}

	err = h.db.Transfer(input)
	if err != nil {
		return errors.Wrap(err, "handler Transfer: db query is not correct")
	}

	return response.Respond(ctx, 201, map[string]interface{} {"message": "OK"})
}

func (h *Handler) GetBalance(ctx *routing.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return errors.BadData.New( "handler GetBalance: bad id")
	}

	balance, err := h.db.GetBalance(uint(id))
	if err != nil {
		return errors.Wrap(err, "handler GetBalance: db query is not correct")
	}

	balance.Currency = "RUB"
	if currency := string(ctx.QueryArgs().Peek("currency")); currency != "" {
		value, err := utils.Exchange(balance.Currency, balance.Total, currency)
		if err == nil {
			log.Println("Exchange service doesn't work", err)
			balance.Currency = currency
			balance.Total = value
		}
	}

	return response.Respond(ctx, 200, map[string]interface{} {"balance": balance})
}

func (h *Handler) GetTransferList(ctx *routing.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return errors.BadData.New( "handler GetTransferList: bad id")
	}

	sort := string(ctx.QueryArgs().Peek("sort"))
	if sort != "amount" && sort != "date" {
		sort = "date"
	}

	limit := ctx.QueryArgs().GetUintOrZero("limit")
	if limit == 0 {
		limit = 20
	}

	offset := ctx.QueryArgs().GetUintOrZero("offset")

	order := string(ctx.QueryArgs().Peek("order"))
	if strings.ToUpper(order) != "DESC" && strings.ToUpper(order) != "ASC" {
		order = "DESC"
	}

	list, err := h.db.GetListTransfer(uint(id), sort, order, limit, offset)
	if err != nil {
		return errors.Wrap( err,"handler GetTransferList: db query error")
	}

	return response.Respond(ctx, 200, map[string]interface{} {"list": list})
}
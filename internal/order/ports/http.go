package ports

import (
	"html/template"
	"log"
	"net/http"

	"github.com/MarinDmitrii/WB-L0/internal/order/builder"
	"github.com/labstack/echo/v4"
)

type HttpOrderHandler struct {
	app *builder.Application
}

func NewHttpOrderHandler(app *builder.Application) HttpOrderHandler {
	return HttpOrderHandler{app: app}
}

func (h HttpOrderHandler) GetOrderByUID(ctx echo.Context, orderUID string) error {
	order, err := h.app.GetOrderByUID.Execute(ctx.Request().Context(), orderUID)
	if err != nil {
		log.Println("здесь косяк!")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	tmpl, err := template.ParseFiles("./web/index.html")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = tmpl.Execute(ctx.Response().Writer, order)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func CustomRegisterHandlers(router EchoRouter, si ServerInterface) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET("/orders/:order_uid", wrapper.GetOrderByUID)
}

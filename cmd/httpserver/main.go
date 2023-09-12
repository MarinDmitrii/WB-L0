package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/MarinDmitrii/WB-L0/cmd/publisher"
	orderBuilder "github.com/MarinDmitrii/WB-L0/internal/order/builder"
	orderPorts "github.com/MarinDmitrii/WB-L0/internal/order/ports"
)

type HttpServer struct {
	orderPorts.HttpOrderHandler
	orderPorts.NatsOrderHandler
}

type Application struct {
	httpServer *http.Server
}

func (a *Application) Run(addr string, debug bool) error {
	router := echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Static("/orders", "web/index.html")

	ctx := context.Background()

	orderApp, orderCleanup := orderBuilder.NewApplication(ctx)
	orderHttpHandler := orderPorts.NewHttpOrderHandler(orderApp)
	orderPorts.CustomRegisterHandlers(router, orderHttpHandler)
	defer orderCleanup()

	a.httpServer = &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    180 * time.Second,
		WriteTimeout:   180 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Server is running...")

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	doneCh := make(chan struct{})
	orderNatsHandler := orderPorts.NewNatsOrderHandler(orderApp)
	orderNatsHandler.NatsSubscriber(ctx, doneCh)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	close(doneCh)

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func main() {
	go publisher.NatsPublisher()

	app := &Application{}
	err := app.Run(":9090", false)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"fmt"
	"log"
	"os"

	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
)

type Service struct {
	Hostname string
}

func main() {

	log.SetPrefix("[APPL] ")

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("failed to get hostname, %v", err)
	}

	service := &Service{
		Hostname: hostname,
	}

	viewEngine := iris.HTML("./views", ".html").Layout("layout.html")
	app := iris.New()
	if os.Getenv("APP_MODE") == "debug" {
		log.Printf("Debug mode is enabled")
		app.Logger().SetLevel("debug")
		viewEngine.Reload(true)
	}
	app.Use(logger.New())
	app.RegisterView(viewEngine)
	app.Get("/", getIndex)
	app.Get("/hostname", getHostname)
	app.Get("/health", getHealth)
	err = app.Listen(
		fmt.Sprintf(":%s", os.Getenv("PORT")),
		iris.WithOtherValue("service", service),
	)
	if err != nil {
		log.Fatalf("failed to listen, %v", err)
	}
}

func GetService(ctx iris.Context) *Service {
	svc, ok := ctx.Application().ConfigurationReadOnly().GetOther()["service"].(*Service)
	if !ok {
		panic("failed to get service")
	}
	return svc
}

func getIndex(ctx iris.Context) {
	svc := GetService(ctx)
	type IndexView struct {
		Title    string
		Hostname string
	}
	ctx.View("index.html", &IndexView{
		Title:    "golang-sample-webapp",
		Hostname: svc.Hostname,
	})
}

func getHostname(ctx iris.Context) {
	svc := GetService(ctx)
	ctx.Text(svc.Hostname)
}

func getHealth(ctx iris.Context) {
	ctx.Text("OK")
}

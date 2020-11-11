package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"golang.org/x/net/context"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func newApp() *iris.Application {
	app := iris.New()
	app.Logger().SetLevel("debug")
	requestLogger := logger.New(logger.Config{
		// Status displays status code
		Status: true,
		// IP displays request's remote address
		IP: true,
		// Method displays the http method
		Method: true,
		// Path displays the request path
		Path: true,
		// Query appends the url query to the Path.
		Query: true,

		// if !empty then its contents derives from `ctx.Values().Get("logger_message")
		// will be added to the logs.
		MessageContextKeys: []string{"logger_message"},

		// if !empty then its contents derives from `ctx.GetHeader("User-Agent")
		MessageHeaderKeys: []string{"User-Agent"},
	})
	// cross-domain
	//crs := cors.New(cors.Options{
	//	AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
	//	AllowCredentials: true,
	//	AllowedMethods:   []string{"PUT", "GET", "OPTIONS"},
	//	//AllowedHeaders:   []string{"*"},
	//	MaxAge: int((24 * time.Hour).Seconds()),
	//})
	app.Use(requestLogger)
	//app.Use(crs)
	app.AllowMethods(iris.MethodOptions)
	app.HandleDir("/", "./res", iris.DirOptions{
		IndexName: "index.html",
	})
	//party := app.Party("/", crs).AllowMethods(iris.MethodOptions)
	//{
	app.Get("/health", healthCheck)
	app.Get("/exam", genExam)
	//}

	return app
}

func main() {
	app := newApp()
	iris.RegisterOnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		// close all hosts
		app.Shutdown(ctx)
	})

	app.Listen(":8080", iris.WithOptimizations,
		iris.WithCharset("UTF-8"),
	)
}

func healthCheck(ctx iris.Context) {
	ctx.WriteString("OK")
}

func genExam(ctx iris.Context) {
	rand.Seed(time.Now().UnixNano())
	examRange := ctx.URLParamIntDefault("range", 10)
	examCount := ctx.URLParamIntDefault("count", 100)
	examSymbol := ctx.URLParamDefault("symbol", "1,2")
	examDiff := ctx.URLParamIntDefault("diff", 0)
	symbols := strings.Split(examSymbol, ",")
	array := make([]string, examCount)
	for i := 0; i < examCount; i++ {
		index := rand.Intn(len(symbols))
		for {
			num1 := rand.Intn(examRange)
			num2 := rand.Intn(examRange)
			if num1 == 0 || num2 == 0 {
				continue
			}
			var result int
			switch symbols[index] {
			case "1":
				if examDiff == 1 {
					if num2%10+num1%10 > 10 {
						result = num1 + num2
					} else {
						continue
					}
				} else {
					result = num1 + num2
				}
			case "2":
				if num1 >= num2 {
					result = num1 - num2
				} else {
					continue
				}
			case "3":
				result = num1 * num2
			case "4":
				if num1 >= num2 && num1%num2 == 0 {
					result = num1 / num2
				} else {
					continue
				}
			}

			if result <= examRange {
				array[i] = strconv.Itoa(num1) + " " + getStringSymbol(symbols[index]) + " " + strconv.Itoa(num2) + " ="
				break
			}
		}

	}
	ctx.JSON(array)

}

func getStringSymbol(numSymbol string) string {
	switch numSymbol {
	case "1":
		return "+"
	case "2":
		return "-"
	case "3":
		return "x"
	case "4":
		return "÷"
	default:
		return "+"
	}
}

package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var apiKey string
var debug bool

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.check-close-windows")
	viper.AddConfigPath(".")

	// flag.Bool("debug", false, "resty debug mode")
	// pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Bool("debug", false, "resty debug mode")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	err := viper.ReadInConfig()
	if err != nil {
		// echo.New().StdLogger.Panicln(err)
	}

	apiKey = viper.GetString("apiKey")
	debug = viper.GetBool("debug")
}

func api(c echo.Context, lat string, long string, unitSystem string, fields string) error {
	client := resty.New()

	client.SetHostURL("https://api.climacell.co")
	client.SetHeader("apikey", apiKey)
	client.QueryParam.Add("lat", lat)
	client.QueryParam.Add("lon", long)
	client.QueryParam.Add("unit_system", unitSystem)
	if fields != "" {
		client.QueryParam.Add("fields", fields)
	}
	client.Debug = debug

	response, err := client.R().Get("v3/weather/realtime")
	if err != nil {
		panic(err)
	}

	switch {
	case 200 <= response.StatusCode() && response.StatusCode() < 300:
		c.Response().Writer.Header().Add("Content-Type", response.Header().Get("Content-Type"))
		return c.String(response.StatusCode(), string(response.Body()))
	default:
		c.Response().Writer.Header().Add("Content-Type", response.Header().Get("Content-Type"))
		return c.String(response.StatusCode(), string(response.Body()))
	}
}

type appConfig struct {
	apiKey string
}

func configureApp() {
}

func main() {
	configureApp()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.Static("./dist"))

	e.GET("api", func(c echo.Context) error {
		// 	return c.JSON(http.StatusOK, struct {
		// 		foo string
		// 		bar float64
		// 	}{"blah", 123.45})
		// })

		err := api(c, c.QueryParam("lat"), c.QueryParam("lon"), c.QueryParam("unit_system"), c.QueryParam("fields"))
		return err
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

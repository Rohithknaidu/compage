package main

import (
	"context"
	"os"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"github.com/Rohithknaidu/Instagram/user-service/config"
	restcontrollers "github.com/Rohithknaidu/Instagram/user-service/pkg/rest/server/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sinhashubham95/go-actuator"
	log "github.com/sirupsen/logrus"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	serviceName  = os.Getenv("SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure     = os.Getenv("INSECURE_MODE")
)

func main() {

	// rest server configuration
	router := echo.New()
	var restTraceProvider *sdktrace.TracerProvider
	if len(serviceName) > 0 && len(collectorURL) > 0 {
		// add opentel
		restTraceProvider = config.InitRestTracer(serviceName, collectorURL, insecure)
		router.Use(OpenTelemetryMiddleware(serviceName))
	}
	defer func() {
		if restTraceProvider != nil {
			if err := restTraceProvider.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}
	}()
	router.Use(middleware.Logger())
	// add actuator
	addActuator(router)
	// add prometheus
	addPrometheus(router)

	userController, err := restcontrollers.NewUserController()
	if err != nil {
		log.Errorf("error occurred: %v", err)
		os.Exit(1)
	}

	v1 := router.Group("/v1")
	{

		v1.POST("/users", userController.CreateUser)

		v1.GET("/users", userController.ListUsers)

		v1.GET("/users/:id", userController.FetchUser)

		v1.PUT("/users/:id", userController.UpdateUser)

		v1.DELETE("/users/:id", userController.DeleteUser)

	}

	Port := ":1337"
	log.Println("Server started")
	if err = router.Start(Port); err != nil {
		log.Errorf("error occurred: %v", err)
		os.Exit(1)
	}

}
func OpenTelemetryMiddleware(serviceName string) echo.MiddlewareFunc {
	return otelecho.Middleware(serviceName)
}

func prometheusHandler() echo.HandlerFunc {
	h := promhttp.Handler()

	return func(c echo.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

func addPrometheus(router *echo.Echo) {
	router.GET("/metrics", prometheusHandler())
}

func addActuator(router *echo.Echo) {
	actuatorHandler := actuator.GetActuatorHandler(&actuator.Config{Endpoints: []int{
		actuator.Env,
		actuator.Info,
		actuator.Metrics,
		actuator.Ping,
		// actuator.Shutdown,
		actuator.ThreadDump,
	},
		Env:     "dev",
		Name:    "user-service",
		Port:    1337,
		Version: "0.0.1",
	})
	echoActuatorHandler := func(ctx echo.Context) error {
		actuatorHandler(ctx.Response(), ctx.Request())
		return nil
	}
	router.GET("/actuator/*endpoint", echoActuatorHandler)
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

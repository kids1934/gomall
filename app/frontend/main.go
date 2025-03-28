package main

import (
	"context"
	"github.com/kids1934/gomall/app/frontend/biz/router"
	"github.com/kids1934/gomall/app/frontend/conf"
	"github.com/kids1934/gomall/app/frontend/infra/mtl"
	"github.com/kids1934/gomall/app/frontend/infra/rpc"
	"github.com/kids1934/gomall/app/frontend/middleware"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/gzip"
	"github.com/hertz-contrib/logger/accesslog"
	hertzprom "github.com/hertz-contrib/monitor-prometheus"
	hertzotelprovider "github.com/hertz-contrib/obs-opentelemetry/provider"
	hertzoteltracing "github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/hertz-contrib/pprof"
	"github.com/hertz-contrib/sessions"
	"github.com/hertz-contrib/sessions/redis"
	"github.com/joho/godotenv"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func main() {
	_ = godotenv.Load()

	mtl.InitMtl()
	rpc.InitClient()
	address := conf.GetConf().Hertz.Address

	p := hertzotelprovider.NewOpenTelemetryProvider(
		hertzotelprovider.WithSdkTracerProvider(mtl.TracerProvider),
		hertzotelprovider.WithEnableMetrics(false),
	)
	defer p.Shutdown(context.Background())
	tracer, cfg := hertzoteltracing.NewServerTracer(hertzoteltracing.WithCustomResponseHandler(func(ctx context.Context, c *app.RequestContext) {
		c.Header("shop-trace-id", oteltrace.SpanFromContext(ctx).SpanContext().TraceID().String())
	}))

	h := server.New(server.WithHostPorts(address), server.WithTracer(
		hertzprom.NewServerTracer(
			"",
			"",
			hertzprom.WithRegistry(mtl.Registry),
			hertzprom.WithDisableServer(true),
		),
	),
		tracer,
	)
	h.LoadHTMLGlob("template/*")
	h.Delims("{{", "}}")

	h.Use(hertzoteltracing.ServerMiddleware(cfg))
	registerMiddleware(h)

	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{"ping": "pong"})
	})

	router.GeneratedRegister(h)

	h.GET("sign-in", func(ctx context.Context, c *app.RequestContext) {
		c.HTML(consts.StatusOK, "sign-in", utils.H{
			"title": "Sign in",
			"next":  c.Query("next"),
		})
	})
	h.GET("sign-up", func(ctx context.Context, c *app.RequestContext) {
		c.HTML(consts.StatusOK, "sign-up", utils.H{
			"title": "Sign up",
		})
	})
	h.GET("/redirect", func(ctx context.Context, c *app.RequestContext) {
		c.HTML(consts.StatusOK, "about", utils.H{
			"title": "Error",
		})
	})
	if os.Getenv("GO_ENV") != "online" {
		h.GET("/robots.txt", func(ctx context.Context, c *app.RequestContext) {
			c.Data(consts.StatusOK, "text/plain", []byte(`User-agent: *
Disallow: /`))
		})
	}

	h.Static("/static", "./")

	h.Spin()
}

func registerMiddleware(h *server.Hertz) {
	// pprof
	if conf.GetConf().Hertz.EnablePprof {
		pprof.Register(h)
	}

	store, err := redis.NewStore(100, "tcp", conf.GetConf().Redis.Address, "", []byte(os.Getenv("SESSION_SECRET")))
	if err != nil {
		panic(err)
	}
	store.Options(sessions.Options{MaxAge: 86400, Path: "/"})
	rs, err := redis.GetRedisStore(store)
	if err == nil {
		rs.SetSerializer(sessions.JSONSerializer{})
	}
	h.Use(sessions.New("cloudwego-shop", store))

	// gzip
	if conf.GetConf().Hertz.EnableGzip {
		h.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	// access log
	if conf.GetConf().Hertz.EnableAccessLog {
		h.Use(accesslog.New())
	}

	// recovery
	h.Use(recovery.Recovery())

	h.OnShutdown = append(h.OnShutdown, mtl.Hooks...)

	// cores
	h.Use(cors.Default())
	middleware.RegisterMiddleware(h)
}

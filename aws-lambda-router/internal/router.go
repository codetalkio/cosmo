package internal

import (
	"fmt"

	"github.com/wundergraph/cosmo/router/core"
	"github.com/wundergraph/cosmo/router/pkg/config"
	"github.com/wundergraph/cosmo/router/pkg/metric"
	"github.com/wundergraph/cosmo/router/pkg/trace"
	"go.uber.org/zap"
)

type Option func(*RouterConfig)

type RouterConfig struct {
	RouterConfigPath     string
	ConfigPath           string
	TelemetryServiceName string
	RouterOpts           []core.Option
	GraphApiToken        string
	HttpPort             string
	EnableTelemetry      bool
	Stage                string
	TraceSampleRate      float64
	Logger               *zap.Logger
}

func NewRouter(opts ...Option) *core.Router {

	rc := &RouterConfig{}

	for _, opt := range opts {
		opt(rc)
	}

	if rc.Logger == nil {
		rc.Logger = zap.NewNop()
	}

	logger := rc.Logger

	routerOpts := []core.Option{
		core.WithLogger(logger),
		core.WithAwsLambdaRuntime(),
	}

	routerConfigPath := rc.RouterConfigPath

	// Loading the config file is optional, but if a path is specified we assume it's required.
	if rc.ConfigPath != "" {
		cfg, err := config.LoadConfig(rc.ConfigPath, "")
		if err != nil {
			logger.Fatal("Could not load config", zap.Error(err), zap.String("path", rc.ConfigPath))
		} else {
			// Selectively apply the config options to the router, based on what is supported.
			routerOpts = append(routerOpts,
				core.WithGraphApiToken(cfg.Config.Graph.Token),
				core.WithGraphQLPath(cfg.Config.GraphQLPath),
				core.WithIntrospection(cfg.Config.IntrospectionEnabled),
				core.WithPlayground(cfg.Config.PlaygroundEnabled),
				core.WithPlaygroundPath(cfg.Config.PlaygroundPath),
				core.WithOverrideRoutingURL(cfg.Config.OverrideRoutingURL),
			)

			if cfg.Config.RouterConfigPath != "" {
				routerConfigPath = cfg.Config.RouterConfigPath
			}
		}
	} else {
		// If no config file is specified, set up the defaults for the Lambda Cosmo Router.
		routerOpts = append(routerOpts,
			core.WithPlayground(true),
			core.WithIntrospection(true),
			core.WithGraphApiToken(rc.GraphApiToken),
		)
	}

	if routerConfigPath == "" {
		routerConfigPath = "router.json"
	}
	routerConfig, err := core.SerializeConfigFromFile(routerConfigPath)
	if err != nil {
		logger.Fatal("Could not read router config", zap.Error(err), zap.String("path", routerConfigPath))
	}
	routerOpts = append(routerOpts, core.WithStaticRouterConfig(routerConfig))

	if rc.HttpPort != "" {
		routerOpts = append(routerOpts, core.WithListenerAddr(":"+rc.HttpPort))
	}

	if rc.EnableTelemetry {
		routerOpts = append(routerOpts,
			core.WithGraphQLMetrics(&core.GraphQLMetricsConfig{
				Enabled:           true,
				CollectorEndpoint: "https://cosmo-metrics.wundergraph.com",
			}),
			core.WithMetrics(&metric.Config{
				Name:    rc.TelemetryServiceName,
				Version: Version,
				OpenTelemetry: metric.OpenTelemetry{
					Enabled: true,
				},
			}),
			core.WithTracing(&trace.Config{
				Enabled: true,
				Name:    rc.TelemetryServiceName,
				Version: Version,
				Sampler: rc.TraceSampleRate,
				Propagators: []trace.Propagator{
					trace.PropagatorTraceContext,
				},
			}),
		)
	}

	if rc.Stage != "" {
		routerOpts = append(routerOpts,
			core.WithGraphQLWebURL(fmt.Sprintf("/%s%s", rc.Stage, "/graphql")),
		)
	}

	r, err := core.NewRouter(append(rc.RouterOpts, routerOpts...)...)
	if err != nil {
		logger.Fatal("Could not create router", zap.Error(err))
	}

	return r
}

func WithRouterConfigPath(path string) Option {
	return func(r *RouterConfig) {
		r.RouterConfigPath = path
	}
}

func WithConfigPath(path string) Option {
	return func(r *RouterConfig) {
		r.ConfigPath = path
	}
}

func WithTelemetryServiceName(name string) Option {
	return func(r *RouterConfig) {
		r.TelemetryServiceName = name
	}
}

func WithRouterOpts(opts ...core.Option) Option {
	return func(r *RouterConfig) {
		r.RouterOpts = append(r.RouterOpts, opts...)
	}
}

func WithGraphApiToken(token string) Option {
	return func(r *RouterConfig) {
		r.GraphApiToken = token
	}
}

func WithHttpPort(port string) Option {
	return func(r *RouterConfig) {
		r.HttpPort = port
	}
}

func WithEnableTelemetry(enable bool) Option {
	return func(r *RouterConfig) {
		r.EnableTelemetry = enable
	}
}

func WithStage(stage string) Option {
	return func(r *RouterConfig) {
		r.Stage = stage
	}
}

func WithTraceSampleRate(rate float64) Option {
	return func(r *RouterConfig) {
		r.TraceSampleRate = rate
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(r *RouterConfig) {
		r.Logger = logger
	}
}

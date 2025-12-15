package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/ICan-TC/lib/db"
	"github.com/ICan-TC/lib/logging"
	"github.com/ICan-TC/lib/tokens"
	"github.com/ICan-TC/users/cmd"
	"github.com/ICan-TC/users/internal/config"
	"github.com/ICan-TC/users/internal/handlers"
	"github.com/ICan-TC/users/internal/service"
)

// Options for the CLI. Pass `--port` or set the `SERVICE_PORT` env var.
type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8888"`
}

// GreetingOutput represents the greeting operation response.
type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

func main() {

	flag.Int("port", 8888, "Port to listen on")
	flag.String("config", "", "Path to config file")
	flag.Parse()

	logging.InitLogger("text")

	config.Load()
	logging.InitLogger(config.Get().Server.LogFormat)

	// --- OpenTelemetry Tracing Setup ---
	// ctx := context.Background()
	// shutdown, err := observability.SetupOTelTracing(ctx, "questara-backend")
	// if err != nil {
	// 	panic(fmt.Sprintf("failed to set up OpenTelemetry: %v", err))
	// }
	// defer shutdown(context.Background())
	// -----------------------------------

	l := logging.L()

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		action := "help"
		if len(os.Args) > 2 {
			switch os.Args[2] {
			case "up", "-up", "--up", "-u":
				action = "up"
			case "down", "-down", "--down", "-d":
				action = "down"
			case "status", "-status", "--status", "-s":
				action = "status"
			case "help", "-help", "--help", "-h":
				action = "help"
			default:
				action = "help"
			}
		}
		cmd.Migrate(action)
		os.Exit(0)
	}

	cfg := config.Get()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)
	if cfg.DB.SSL {
		dsn += "?sslmode=require"
	} else {
		dsn += "?sslmode=disable"
	}

	dbconn, err := db.New(dsn)
	if err != nil {
		panic(err)
	}

	l.Info().Msg("Starting API server")
	zerolog.Ctx(context.Background()).Info().Msg("Test log from zerolog.Ctx(context.Background())")
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		cfg := config.Get()
		router := chi.NewMux()

		h := cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"https://*", "http://*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		})
		router.Use(h)

		api := humachi.New(router, huma.DefaultConfig("API Server", "1.0.0"))

		// Wire up the handlers
		usersSvc, err := service.NewUsersService(dbconn)
		if err != nil {
			l.Err(err).Msg("Skipping Users Service")
		} else {
			handlers.RegisterUsersRoutes(api, usersSvc)
		}

		studentsSvc, err := service.NewStudentsService(dbconn)
		if err != nil {
			l.Err(err).Msg("Skipping Students Service")
		} else {
			handlers.RegisterStudentsRoutes(api, studentsSvc)
		}

		teachersSvc, err := service.NewTeachersService(dbconn)
		if err != nil {
			l.Err(err).Msg("Skipping Teachers Service")
		} else {
			handlers.RegisterTeachersRoutes(api, teachersSvc)
		}

		employeesSvc, err := service.NewEmployeesService(dbconn)
		if err != nil {
			l.Err(err).Msg("Skipping Employees Service")
		} else {
			handlers.RegisterEmployeesRoutes(api, employeesSvc)
		}

		parentsSvc, err := service.NewParentsService(dbconn)
		if err != nil {
			l.Err(err).Msg("Skipping Parents Service")
		} else {
			handlers.RegisterParentsRoutes(api, parentsSvc)
		}

		studentParentsSvc, err := service.NewStudentParentsService(dbconn)
		if err != nil {
			l.Err(err).Msg("Skipping StudentParents Service")
		} else {
			handlers.RegisterStudentParentsRoutes(api, studentParentsSvc)
		}

		groupsSvc, err := service.NewGroupsService(dbconn)
		if err != nil {
			l.Err(err).Msg("Skipping Groups Service")
		} else {
			handlers.RegisterGroupsRoutes(api, groupsSvc)
		}

		enrollmentsSvc, err := service.NewEnrollmentsService(dbconn)
		if err != nil {
			l.Err(err).Msg("Skipping Enrollments Service")
		} else {
			handlers.RegisterEnrollmentsRoutes(api, enrollmentsSvc)
		}

		tokenProvider, err := tokens.NewTokenProvider(tokens.TokenProviderArgs{
			Secret:          cfg.Auth.Secret,
			AccessTokenTTL:  cfg.Auth.AccessTokenTTL,
			RefreshTokenTTL: cfg.Auth.RefreshTokenTTL,
		})
		if err != nil {
			l.Err(err).Msg("Failed to create token provider, this is a critical module, exiting")
			os.Exit(1)
		}

		tokensSvc := service.NewTokensService(tokenProvider, dbconn)
		authSvc, err := service.NewAuthService(usersSvc, tokensSvc)
		if err != nil {
			l.Err(err).Msg("Skipping Auth Service")
		} else {
			handlers.RegisterAuthRoutes(api, authSvc)
		}

		huma.Get(api, "/greeting/{name}", func(ctx context.Context, input *struct {
			Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
		},
		) (*GreetingOutput, error) {
			resp := &GreetingOutput{}
			resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)
			return resp, nil
		})

		hooks.OnStart(func() {
			l.Info().Int("port", cfg.Server.Port).Msg("API server listening")
			importOtelHTTP := func() {} // dummy to ensure import
			_ = importOtelHTTP
			wrapped := otelhttp.NewHandler(router, "http.server")
			wrappedWithLogger := logging.RequestLoggingHandler(wrapped)
			if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), wrappedWithLogger); err != nil {
				panic(fmt.Sprintf("failed to start server: %v", err))
			}
		})
	})
	cli.Run()
}

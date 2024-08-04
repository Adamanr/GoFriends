package main

import (
	"accessCloude/internal/config"
	"accessCloude/internal/storage"
	"log"
	"net"
	"net/http"

	"log/slog"

	api "accessCloude/internal/handler"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpmiddleware "github.com/oapi-codegen/nethttp-middleware"
)

func main() {
	cfg, err := config.GetConfigs()
	if err != nil {
		slog.Info("Failed to get config")
		panic(err)
	}

	swagger, err := api.GetSwagger()
	if err != nil {
		panic(err)
	}

	swagger.Servers = nil
	r := chi.NewRouter()

	db := storage.NewDatabase(cfg)

	accessCloude := api.NewAccessCloude(db)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.Logger)
	api.HandlerFromMux(accessCloude, r)

	h := httpmiddleware.OapiRequestValidator(swagger)(r)
	openapi3filter.RegisterBodyDecoder("image/jpeg", openapi3filter.FileBodyDecoder)
	openapi3filter.RegisterBodyDecoder("image/png", openapi3filter.FileBodyDecoder)

	s := &http.Server{
		Handler: h,
		Addr:    net.JoinHostPort(cfg.CS.Host, cfg.CS.Port),
	}

	slog.Info("Starting server on", cfg.CS.Host, cfg.CS.Port)
	log.Fatal(s.ListenAndServe())
}

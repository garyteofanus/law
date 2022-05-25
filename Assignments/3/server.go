package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

type Server struct {
	router *chi.Mux
	logger *logrus.Logger
}

func NewServer(logger *logrus.Logger) Server {
	return Server{
		router: chi.NewRouter(),
		logger: logger,
	}
}

func (s Server) SetupRoutes() {
	handlerLogger := s.logger.WithFields(logrus.Fields{
		"method": "handler",
	})

	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		userQueryParam := r.URL.Query().Get("param")
		handlerLogger.Infof("query param: %s", userQueryParam)

		t, err := template.New("response template").Parse(`{{.}}`)
		if err != nil {
			handlerLogger.Errorf("failed to parse template: %v", err)
		}

		typeParam := r.URL.Query().Get("type")
		handlerLogger.Infof("type param: %s", typeParam)
		switch typeParam {
		case "json":
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(map[string]string{
				"query": userQueryParam,
			}); err != nil {
				handlerLogger.Errorf("failed to encode json: %v", err)
			}
		case "html", "":
			w.Header().Set("Content-Type", "text/html")
			if err := t.Execute(w, userQueryParam); err != nil {
				handlerLogger.Errorf("failed to execute template: %v", err)
			}
		}
	})

}

func (s Server) Start(port string) error {
	if err := http.ListenAndServe(":"+port, s.router); err != nil {
		return err
	}
	return nil
}

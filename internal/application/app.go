package application

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/fernandoocampo/basic-micro/internal/adapter/stores"
	"github.com/fernandoocampo/basic-micro/internal/adapter/web"
	"github.com/fernandoocampo/basic-micro/internal/pets"
	"github.com/fernandoocampo/basic-micro/internal/setups"
	"go.uber.org/zap"
)

// Event contains an application event.
type Event struct {
	Message string
	Error   error
}

// Setup contains application metadata
type Setup struct {
	Version    string
	BuildDate  string
	CommitHash string
}

// Server is the server of our application.
type Server struct {
	logger     *zap.Logger
	store      *stores.Store
	setup      setups.Application
	version    string
	buildDate  string
	commitHash string
}

var (
	errStartingApplication = errors.New("unable to start application")
)

func NewServer(settings Setup) *Server {
	newServer := Server{
		version:    settings.Version,
		buildDate:  settings.BuildDate,
		commitHash: settings.CommitHash,
	}

	return &newServer
}

func (s *Server) Run() error {
	s.notifyStart()

	confError := s.loadConfiguration()
	if confError != nil {
		return errStartingApplication
	}

	loggerError := s.initializeLogger()
	if loggerError != nil {
		return errStartingApplication
	}

	s.logger.Debug("application configuration", zap.String("parameters", fmt.Sprintf("%+v", s.setup)))

	s.logger.Info("starting database connection")

	err := s.createStorer()
	if err != nil {
		return errStartingApplication
	}

	petServiceSetup := pets.ServiceSetup{
		Storer: s.store,
		Logger: s.logger,
	}
	petService := pets.NewService(petServiceSetup)
	petEndpoints := pets.NewEndpoints(petService, s.logger)

	eventStream := make(chan Event)
	s.listenToOSSignal(eventStream)
	s.startWebServer(petEndpoints, eventStream)

	eventMessage := <-eventStream
	s.logger.Info("ending server", zap.String("event", eventMessage.Message))

	if eventMessage.Error != nil {
		s.logger.Error("ending server with error", zap.Error(eventMessage.Error))

		return errStartingApplication
	}

	return nil
}

func (s *Server) initializeLogger() error {
	var logger *zap.Logger
	var err error

	if s.setup.LogLevel == setups.ProductionLog {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		fmt.Println(
			"level", "ERROR",
			"msg", "unable to initialize logger",
			"error", err,
		)

		return fmt.Errorf("unable to initialize logger: %w", err)
	}

	s.logger = logger

	return nil
}

func (s *Server) notifyStart() {
	log.Println(
		"starting service",
		"version:", s.version,
		"commit:", s.commitHash,
		"build date:", s.buildDate,
	)
}

// Stop stop application, take advantage of this to clean resources
func (s *Server) Stop() {
	s.logger.Info("stopping the application")
}

func (s *Server) listenToOSSignal(eventStream chan<- Event) {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		osSignal := (<-c).String()
		event := Event{
			Message: osSignal,
		}
		eventStream <- event
	}()
}

// startWebServer starts the web server.
func (s *Server) startWebServer(petEndpoints pets.Endpoints, eventStream chan<- Event) {
	go func() {
		s.logger.Info("starting http server", zap.String("port", s.setup.ApplicationPort))
		router := petsRouter{
			router:    web.NewRouter(),
			endpoints: petEndpoints,
			decoders:  web.NewPetDecoders(s.logger),
			encoders:  web.NewPetEncoders(s.logger),
		}
		handler := newPetsRouter(router)
		err := http.ListenAndServe(s.setup.ApplicationPort, handler)
		if err != nil {
			eventStream <- Event{
				Message: "web server was ended with error",
				Error:   err,
			}
			return
		}
		eventStream <- Event{
			Message: "web server was ended",
		}
	}()
}

func (s *Server) loadConfiguration() error {
	applicationSetUp, err := setups.Load()
	if err != nil {
		log.Println("level", "ERROR", "msg", "application setup could not be loaded", "error", err)

		return errors.New("application setup could not be loaded")
	}
	s.setup = applicationSetUp
	return nil
}

func (s *Server) createStorer() error {
	storeSetup := stores.Setup{
		Logger: s.logger,
	}
	storer := stores.NewStore(storeSetup)
	s.store = storer

	return nil
}

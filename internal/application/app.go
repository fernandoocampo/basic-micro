package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/fernandoocampo/basic-micro/internal/adapter/stores"
	"github.com/fernandoocampo/basic-micro/internal/adapter/telemetry"
	"github.com/fernandoocampo/basic-micro/internal/adapter/web"
	"github.com/fernandoocampo/basic-micro/internal/pets"
	"github.com/fernandoocampo/basic-micro/internal/setups"
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
	logger     *slog.Logger
	store      *stores.Store
	setup      setups.Application
	version    string
	buildDate  string
	commitHash string
}

const ServiceName = "basic-micro"

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
	slog.Info("starting application")
	ctx, stop := s.initializeApplication()
	defer stop()

	slog.Info("loading configuration")
	confError := s.loadConfiguration()
	if confError != nil {
		return errStartingApplication
	}

	slog.Info("initializing logger")
	loggerError := s.initializeLogger()
	if loggerError != nil {
		return errStartingApplication
	}

	slog.Info("initialize telemetry")
	shutdownTelemetry, telemetryErr := s.initializeTelemetry(ctx)
	if telemetryErr != nil {
		s.logger.Error("initializing telemetry", "error", telemetryErr.Error)

		return errStartingApplication
	}
	defer shutdownTelemetry(ctx)

	s.logger.Debug("application configuration", slog.String("parameters", fmt.Sprintf("%+v", s.setup)))

	s.logger.Info("starting database connection")
	err := s.createStorer()
	if err != nil {
		return errStartingApplication
	}

	s.logger.Info("initializing service")
	petServiceSetup := pets.ServiceSetup{
		Storer: s.store,
		Logger: s.logger,
	}
	petService := pets.NewService(petServiceSetup)

	s.logger.Info("initializing endpoints")
	petEndpoints := pets.NewEndpoints(petService, s.logger)

	eventStream := Or(
		s.stopApplication(ctx),
		s.startWebServer(petEndpoints),
	)

	eventMessage := <-eventStream
	s.logger.Info("ending server", slog.String("event", eventMessage.Message))

	if eventMessage.Error != nil {
		s.logger.Error("ending server with error", "error", eventMessage.Error)

		return errStartingApplication
	}

	return nil
}

func (s *Server) initializeApplication() (context.Context, context.CancelFunc) {
	s.notifyStart()

	ctx, stopFunc := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		os.Kill,
		syscall.SIGTERM,
	)

	return ctx, stopFunc
}

func (s *Server) initializeLogger() error {
	logLevel := slog.LevelDebug

	if s.setup.LogLevel == setups.ProductionLog {
		slog.Info("setting log level", slog.String("level", "INFO"))
		logLevel = slog.LevelInfo
	}

	handlerOptions := &slog.HandlerOptions{
		Level: logLevel,
	}

	loggerHandler := slog.NewJSONHandler(os.Stdout, handlerOptions)
	logger := slog.New(loggerHandler)

	logger.Info(
		"logger has been initialized",
		slog.String("level", handlerOptions.Level.Level().String()),
	)

	slog.SetDefault(logger)

	s.logger = logger

	return nil
}

func (s *Server) initializeTelemetry(ctx context.Context) (func(context.Context), error) {
	telemetrySetup := telemetry.OtelSDKSetup{
		ServiceName:    ServiceName,
		ServiceVersion: s.version,
		Logger:         s.logger,
	}

	telemetryShutdown, err := telemetry.NewOtelSDK(ctx, telemetrySetup)
	if err != nil {
		s.logger.Error("initializing telemetry", "error", err)

		return nil, errors.New("unable to initialize telemetry")
	}

	return telemetryShutdown, nil
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

func (s *Server) listenToOSSignal() <-chan Event {
	osSignalStream := make(chan Event)
	go func() {
		defer close(osSignalStream)
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		osSignal := (<-c).String()
		event := Event{
			Message: osSignal,
		}
		osSignalStream <- event
	}()
	return osSignalStream
}

// startWebServer starts the web server.
func (s *Server) startWebServer(petEndpoints pets.Endpoints) <-chan Event {
	serverSignalStream := make(chan Event)
	go func() {
		defer close(serverSignalStream)
		s.logger.Info("starting http server", slog.String("port", s.setup.ApplicationPort))
		router := petsRouter{
			router:    web.NewRouter(),
			endpoints: petEndpoints,
			decoders:  web.NewPetDecoders(s.logger),
			encoders:  web.NewPetEncoders(s.logger),
		}
		handler := newPetsRouter(router)
		err := http.ListenAndServe(s.setup.ApplicationPort, handler)
		if err != nil {
			serverSignalStream <- Event{
				Message: "web server was ended with error",
				Error:   err,
			}
			return
		}
		serverSignalStream <- Event{
			Message: "web server was ended",
		}
	}()
	return serverSignalStream
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

func (s *Server) stopApplication(ctx context.Context) <-chan Event {
	stopAppSignal := make(chan Event)

	go func(ctx context.Context) {
		defer close(stopAppSignal)
		<-ctx.Done()
		message := ctx.Err().Error()
		s.logger.Info("context was cancelled", "message", message)
		stopAppSignal <- Event{
			Message: fmt.Sprintf("context was cancelled: %s", message),
		}
	}(ctx)

	return stopAppSignal
}

func Or(channels ...<-chan Event) <-chan Event {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan Event)
	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-Or(append(channels[3:], orDone)...):
			}
		}
	}()

	return orDone
}

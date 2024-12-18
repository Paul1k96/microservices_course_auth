package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Paul1k96/microservices_course_auth/internal/interceptor"
	"github.com/Paul1k96/microservices_course_auth/pkg/proto/gen/user_v1"
	_ "github.com/Paul1k96/microservices_course_auth/statik" // required for statik
	"github.com/Paul1k96/microservices_course_platform_common/pkg/closer"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// App is the main application structure.
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
	swaggerServer   *http.Server
	logger          *slog.Logger
}

// NewApp creates a new instance of the App.
func NewApp(ctx context.Context, logger *slog.Logger) (*App, error) {
	a := &App{
		logger: logger,
	}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to init dependencies: %w", err)
	}

	return a, nil
}

// Run runs the application.
// nolint: funlen // this is the main function of the app
func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	ctx, cancel := context.WithCancel(ctx)

	errGroup, ctx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		select {
		case <-ctx.Done():
			a.logger.Info("terminating: context cancelled")
		case <-waitForSignal():
			a.logger.Info("terminating: via signal")
		}

		cancel()
		return nil
	})

	errGroup.Go(func() error {
		userCreateConsumer, err := a.serviceProvider.UserCreateConsumer(ctx)
		if err != nil {
			a.logger.Error("failed to get user create consumer", slog.String("error", err.Error()))
			return fmt.Errorf("failed to get user create consumer: %w", err)
		}

		err = userCreateConsumer.RunConsumer(ctx)
		if err != nil {
			a.logger.Error("failed to run user create consumer", slog.String("error", err.Error()))
			return fmt.Errorf("failed to run user create consumer: %w", err)
		}

		return nil
	})

	errGroup.Go(func() error {
		err := a.runGRPCServer(ctx)
		if err != nil {
			log.Printf("failed to run grpc server: %s", err.Error())
			return fmt.Errorf("failed to run grpc server: %w", err)
		}

		return nil
	})

	errGroup.Go(func() error {
		err := a.runHTTPServer(ctx)
		if err != nil {
			log.Printf("failed to run http server: %s", err.Error())
			return fmt.Errorf("failed to run http server: %w", err)
		}

		return nil
	})

	errGroup.Go(func() error {
		err := a.runSwaggerServer(ctx)
		if err != nil {
			log.Printf("failed to run swagger server: %s", err.Error())
			return fmt.Errorf("failed to run swagger server: %w", err)
		}

		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("failed to run app: %w", err)
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider(a.logger)
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(interceptor.ValidateInterceptor),
	)

	userV1Impl, err := a.serviceProvider.UserV1Impl(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user v1 implementation: %w", err)
	}

	reflection.Register(a.grpcServer)
	user_v1.RegisterUserServer(a.grpcServer, userV1Impl)

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	grpcOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := user_v1.RegisterUserHandlerFromEndpoint(ctx, mux, a.serviceProvider.GRPCConfig().GetAddress(), grpcOpts)
	if err != nil {
		return fmt.Errorf("failed to register user handler from endpoint: %w", err)
	}

	httpConfig, err := a.serviceProvider.HTTPConfig()
	if err != nil {
		return fmt.Errorf("failed to get http config: %w", err)
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:              httpConfig.GetAddress(),
		ReadHeaderTimeout: httpConfig.GetReadHeaderTimeout(),
		Handler:           corsMiddleware.Handler(mux),
	}

	return nil
}

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFS, err := fs.New()
	if err != nil {
		return fmt.Errorf("failed to create statik fs: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFS)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json", a.logger))

	swaggerConf, err := a.serviceProvider.SwaggerConfig()
	if err != nil {
		return fmt.Errorf("failed to get swagger config: %w", err)
	}

	a.swaggerServer = &http.Server{
		Addr:              swaggerConf.GetAddress(),
		ReadHeaderTimeout: swaggerConf.GetReadHeaderTimeout(),
		Handler:           mux,
	}
	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	a.logger.Info("GRPC server is running on", slog.Any("addr", a.serviceProvider.GRPCConfig().GetAddress()))

	errChan := make(chan error, 1)
	go func() {
		lis, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().GetAddress())
		if err != nil {
			errChan <- fmt.Errorf("failed to listen: %w", err)
			return
		}

		if err = a.grpcServer.Serve(lis); err != nil {
			errChan <- fmt.Errorf("failed to serve: %w", err)
			return
		}
	}()

	select {
	case <-ctx.Done():
		a.logger.Info("GRPC server is stopping")
	case err := <-errChan:
		return fmt.Errorf("failed to run grpc server: %w", err)
	}

	a.grpcServer.GracefulStop()

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	config, err := a.serviceProvider.HTTPConfig()
	if err != nil {
		return fmt.Errorf("failed to get http config: %w", err)
	}

	a.logger.Info("HTTP server is running on", slog.Any("addr", config.GetAddress()))

	errChan := make(chan error, 1)
	go func() {
		if err = a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("failed to listen and serve: %w", err)
			return
		}
	}()

	select {
	case <-ctx.Done():
		a.logger.Info("HTTP server is stopping")
	case err = <-errChan:
		return fmt.Errorf("failed to run http server: %w", err)
	}

	ctxShutdown, cancel := context.WithTimeout(context.Background(), config.GetGracefulShutdownTimeout())
	defer cancel()

	err = a.httpServer.Shutdown(ctxShutdown)
	if err != nil {
		return fmt.Errorf("failed to shutdown http server: %w", err)
	}

	return nil
}

func (a *App) runSwaggerServer(ctx context.Context) error {
	config, err := a.serviceProvider.SwaggerConfig()
	if err != nil {
		return fmt.Errorf("failed to get swagger config: %w", err)
	}

	a.logger.Info("Swagger server is running on", slog.Any("addr", config.GetAddress()))

	errChan := make(chan error, 1)
	go func() {
		if err = a.swaggerServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("failed to listen and serve: %w", err)
			return
		}
	}()

	select {
	case <-ctx.Done():
		a.logger.Info("Swagger server is stopping")
	case err = <-errChan:
		return fmt.Errorf("failed to run swagger server: %w", err)
	}

	ctxShutdown, cancel := context.WithTimeout(context.Background(), config.GetGracefulShutdownTimeout())
	defer cancel()

	err = a.swaggerServer.Shutdown(ctxShutdown)
	if err != nil {
		return fmt.Errorf("failed to shutdown swagger server: %w", err)
	}

	return nil
}

func waitForSignal() chan os.Signal {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	return sigCh
}

func serveSwaggerFile(file string, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		logger.Info("serving swagger file", slog.String("path", file))

		statikFS, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("open swagger file", slog.String("path", file))

		f, err := statikFS.Open(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() {
			closeErr := f.Close()
			if closeErr != nil {
				logger.Error(
					"failed to close file",
					slog.String("path", file),
					slog.String("error", closeErr.Error()))
			}
		}()

		logger.Info("read swagger file", slog.String("path", file))

		content, err := io.ReadAll(f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("Write swagger file", slog.String("path", file))

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("served swagger file", slog.String("path", file))
	}
}

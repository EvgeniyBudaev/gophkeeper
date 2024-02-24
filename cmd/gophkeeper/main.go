package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/EvgeniyBudaev/gophkeeper/internal/adapters/store"
	"github.com/EvgeniyBudaev/gophkeeper/internal/app"
	"github.com/EvgeniyBudaev/gophkeeper/internal/config"
	"github.com/EvgeniyBudaev/gophkeeper/internal/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

const (
	timeoutServerShutdown = time.Second * 5
	timeoutShutdown       = time.Second * 10
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelCtx()

	l, err := logger.NewLogger()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	c, err := config.Load(l)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	conn, err := store.NewPostgresConnection(c)
	if err != nil {
		return fmt.Errorf("failed new postgres connection: %w", err)
	}

	s := store.NewStore(conn)
	err = conn.Ping()
	if err != nil {
		return fmt.Errorf("failed ping database: %w", err)
	}

	wg := &sync.WaitGroup{}
	defer func() {
		wg.Wait()
	}()

	wg.Add(1)
	go func() {
		defer l.Info("closed DB")
		defer wg.Done()
		<-ctx.Done()
		conn.Close()
	}()

	componentsErrs := make(chan error, 1)

	a := app.NewApp(c, s, l.Named("app"))
	srv, err := a.NewServer()
	if err != nil {
		l.Fatalf("error creating server: %w", err)
	}

	//go func(errs chan<- error) {
	//	if c.EnableHTTPS {
	//		_, errCert := os.ReadFile(c.TLSCertPath)
	//		_, errKey := os.ReadFile(c.TLSKeyPath)
	//		if errors.Is(errCert, os.ErrNotExist) || errors.Is(errKey, os.ErrNotExist) {
	//			privateKey, certBytes, err := app.CreateCertificates(l.Named("certs-builder"))
	//			if err != nil {
	//				errs <- fmt.Errorf("error creating tls certs: %w", err)
	//				return
	//			}
	//			if err := app.WriteCertificates(certBytes, c.TLSCertPath, privateKey, c.TLSKeyPath, l); err != nil {
	//				errs <- fmt.Errorf("error writing tls certs: %w", err)
	//				return
	//			}
	//		}
	//		srv.TLSConfig = &tls.Config{
	//			MinVersion:         tls.VersionTLS12,
	//			ClientAuth:         tls.RequestClientCert,
	//			KeyLogWriter:       bufio.NewWriter(os.Stdout),
	//			InsecureSkipVerify: true,
	//		}
	//		if err := srv.ListenAndServeTLS(c.TLSCertPath, c.TLSKeyPath); err != nil {
	//			if errors.Is(err, http.ErrServerClosed) {
	//				return
	//			}
	//			errs <- fmt.Errorf("run tls server has failed: %w", err)
	//			return
	//		}
	//	}
	//	l.Warnf("serving http server %s without TLS: Use only for development", srv.Addr)
	//	if err := srv.ListenAndServe(); err != nil {
	//		if errors.Is(err, http.ErrServerClosed) {
	//			return
	//		}
	//		errs <- fmt.Errorf("run server has failed: %w", err)
	//	}
	//}(componentsErrs)

	go func(errs chan<- error) {
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			errs <- fmt.Errorf("run server has failed: %w", err)
		}
	}(componentsErrs)

	wg.Add(1)
	go func() {
		defer l.Info("server has been shutdown")
		defer wg.Done()
		<-ctx.Done()

		shutdownTimeoutCtx, cancelShutdownTimeoutCtx := context.WithTimeout(context.Background(), timeoutServerShutdown)
		defer cancelShutdownTimeoutCtx()
		if err := srv.Shutdown(shutdownTimeoutCtx); err != nil {
			l.Errorf("an error occurred during server shutdown: %v", err)
		}
	}()

	select {
	case <-ctx.Done():
	case err := <-componentsErrs:
		l.Error(err)
		cancelCtx()
	}

	go func() {
		ctx, cancelCtx := context.WithTimeout(context.Background(), timeoutShutdown)
		defer cancelCtx()
		<-ctx.Done()
		l.Fatal("failed to gracefully shutdown the service")
	}()

	return nil
}

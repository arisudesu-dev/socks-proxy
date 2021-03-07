package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/armon/go-socks5"
	"github.com/caarlos0/env"
)

var config struct {
	User     string `env:"PROXY_USER"     envDefault:""`
	Password string `env:"PROXY_PASSWORD" envDefault:""`
	Port     string `env:"PROXY_PORT"     envDefault:"1080"`
}

func main() {
	err := env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}

	socsk5conf := &socks5.Config{
		Rules: &socks5.PermitCommand{
			EnableConnect: true,
		},
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}

	if config.User+config.Password != "" {
		creds := socks5.StaticCredentials{
			config.User: config.Password,
		}
		cator := socks5.UserPassAuthenticator{Credentials: creds}
		socsk5conf.AuthMethods = []socks5.Authenticator{cator}
	}

	server, err := socks5.New(socsk5conf)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Start listening on port %s", config.Port)

	listener, err := net.Listen("tcp", ":"+config.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	errCh := make(chan error)
	sigCh := make(chan os.Signal)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Start accepting proxy connections")
	go func() {
		errCh <- server.Serve(listener)
	}()

	select {
	case err := <-errCh:
		log.Fatalf("Server error, terminating: %+v", err)
	case sig := <-sigCh:
		log.Printf("Normal shutdown by signal: %+v", sig)
	}

	os.Exit(0)
}

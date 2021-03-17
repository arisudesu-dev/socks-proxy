package main

import (
	"errors"
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

	brkCh := make(chan struct{})
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Start accepting proxy connections")
	go func() {
		log.Printf("Signal received: %+v", <-sigCh)
		_ = listener.Close()
		close(brkCh)
	}()

	if err := server.Serve(listener); !IsErrNetClosing(err) {
		log.Fatalf("Server error, terminating: %+v", err)
	}

	log.Printf("Waiting for server shutdown...")
	<-brkCh
	os.Exit(0)
}

func IsErrNetClosing(err error) bool {
	var ErrNetClosing = errors.New("use of closed network connection")

	if e, ok := err.(*net.OpError); ok {
		return e.Err.Error() == ErrNetClosing.Error()
	}
	return false
}

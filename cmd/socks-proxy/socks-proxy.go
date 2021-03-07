package main

import (
	"log"
	"os"

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

	log.Printf("Start listening proxy service on port %s\n", config.Port)
	if err := server.ListenAndServe("tcp", ":"+config.Port); err != nil {
		log.Fatal(err)
	}
}

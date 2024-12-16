package serve

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/z411392/golab/config"
	"github.com/z411392/golab/container"
)

const use = "serve"

const readTimeout = 5 * time.Second
const writeTimeout = 10 * time.Second

func runE(command *cobra.Command, args []string) (err error) {
	container.Init()
	defer container.Release()
	var server *http.Server
	address := fmt.Sprintf("%s:%d", host, port)
	handler := config.NewHttpHandler()
	server = &http.Server{
		Addr:         address,
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	failed := make(chan error)
	go (func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			failed <- err
		}
		close(failed)
	})()
	select {
	case <-stop:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = server.Shutdown(ctx)
		if err != nil {
			log.Printf("%s", err.Error())
		}
	case err = <-failed:
	}
	return
}

var (
	host string
	port int
)

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:  use,
		RunE: runE,
	}
	flag := command.Flags()
	defaultHost := os.Getenv("HOST")
	defaultPort, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		defaultPort = 8080
	}
	flag.StringVar(&host, "host", defaultHost, "") // -h 已經被 helper 用走
	flag.IntVarP(&port, "port", "p", defaultPort, "")
	return command
}

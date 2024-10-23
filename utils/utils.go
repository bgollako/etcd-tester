package utils

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Fatal(err error) {
	if err != nil {
		log.Println(err)
	}
}

func FatalAnyErr(a any, err error) any {
	if err != nil {
		log.Fatalln(err)
	}
	return a
}

func GetShutdownHandle() <-chan os.Signal {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)
	return shutdownCh
}

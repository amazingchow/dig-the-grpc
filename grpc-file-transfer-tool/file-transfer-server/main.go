package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
)

var (
	portFlag     = flag.Int("port", 8999, "server port")
	certFileFlag = flag.String("cert", "grpc-file-transfer-tool/cert/cert.pem", "cert file")
	keyFileFlag  = flag.String("key", "grpc-file-transfer-tool/cert/key.pem", "private key file")
)

func main() {
	flag.Parse()

	cfg := &GrpcStreamServerCfg{
		Port: *portFlag,
		Cert: *certFileFlag,
		Key:  *keyFileFlag,
	}

	srv, err := NewGrpcStreamServer(cfg)
	if err != nil {
		panic(err)
	}
	if err = srv.Init(); err != nil {
		panic(err)
	}

	go srv.Run()
	defer func() {
		srv.Close()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
MAIN_LOOP:
	for { // nolint
		select {
		case <-sigCh:
			break MAIN_LOOP
		}
	}
}

package main

import (
	//"flag"
	//pb "github.com/flamurey/prototrade"
	//"golang.org/x/net/context"
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/grpclog"
	//"github.com/flamurey/tranzaq-connector/connector/wrapper"
	//"log"
	//"os"

	"flag"
	"github.com/flamurey/tranzaq-connector/connector/console"
	log "github.com/flamurey/tranzaq-connector/connector/logger"
	"github.com/flamurey/tranzaq-connector/connector/wrapper"
)

var (
	wrapperDir      = flag.String("wrapper_path", "/home/lesha/ws/ats/tranzaq-docker/bin", "path to exe of wrapper program")
	winePath        = flag.String("wine_path", "/usr/bin/wine", "The path to wine for runnig tranzaq wrapper")
	interectiveMode = flag.Bool("i", false, "Interactive mode. Commands can be created by console")
	serviceMode     = flag.Bool("service", false, "Service mode. Appication trying connect to server from which will receive commands")
	serverAddr      = flag.String("server_addr", "127.0.0.1:10000", "The grpc server address")
)

func main() {
	flag.Parse()
	log.Init()
	//start tranzaq wrapper
	tranzaq := wrapper.New(*wrapperDir, *winePath)
	err := tranzaq.Start()
	if err != nil {
		log.WithError(err).Fatalf("Wrapper program couldn't start")
	}
	defer tranzaq.Close()

	if *interectiveMode {
		console.StartReadInput(tranzaq)
	}
	if *serviceMode {
		//TODO connect to grpc server
	}
}

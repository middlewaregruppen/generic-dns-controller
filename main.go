package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/middlewaregruppen/generic-dns-controller/dns"
)

func main() {
	opts := slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}
	handler := slog.NewJSONHandler(os.Stdout, &opts)
	slog.SetDefault(slog.New(handler))
	//logger := slog.New(handler)

	clientset, err := dns.RestK8sClient()
	if err != nil {
		slog.Error("Error creating a kubernetes client")
		return
	}
	// use switch statement for different DNS providers
	switch os.Getenv("DNS_PROVIDER") {
	case "ExternalDNS":
		//Use ExternalDNS
	case "Infoblox":
		//Use Infoblox
	case "Route53":
		//Use Route53
	default:
		log.Fatalf("Unknown DNS Provider")
	}
	//initialising DNS Provider
	dnsProvider, err := dns.NewInfobloxProvider()
	if err != nil {
		log.Fatalf("Error initialising Infoblox provider: %v", err)
	}

	// initialise DNS Controller
	controller := dns.NewDNSController(clientset, dnsProvider)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go controller.Run(ctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Println("Shutting down DNS Controller")
}

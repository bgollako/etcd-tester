package main

import (
	"context"
	"etcd-tester/client"
	"etcd-tester/clog"
	topics "etcd-tester/topics/generate"
	"etcd-tester/utils"
	"flag"
	"os"
	"strings"
)

func main() {
	// Reading cmdline arguments
	numClients := flag.Int("numClients", 1, "number of etcd clients")
	clientName := flag.String("clientName", "client-1", "name of client")
	flag.Parse()

	ctx := clog.NewContextWithDefaultLogger(context.Background())
	logger := clog.MustFromContext(ctx)
	// Reading topics from topics/topics.txt
	data := string(utils.FatalAnyErr(os.ReadFile(topics.Topics)).([]byte))
	allTopics := strings.Split(data, topics.Delimiter)
	// The last topic is empty since it ends with ,
	numTopicsPerClient := (len(allTopics) - 1) / *numClients
	start := 0
	logger.Sugar().Infof("Total no of topics - %d", len(allTopics)-1)
	logger.Sugar().Infof("Generating %d clients, with %d topic per client", *numClients, numTopicsPerClient)

	var clients []client.Client
	for i := 0; i < *numClients; i++ {
		client := client.NewClient(ctx, *clientName, allTopics[start:start+numTopicsPerClient], []string{"localhost:2379"})
		utils.Fatal(client.Start())
		clients = append(clients, client)
		start += numTopicsPerClient
	}
	logger.Info("Completed launching all the clients!!")

	shutdownCh := utils.GetShutdownHandle()
	<-shutdownCh

	logger.Info("Received shutdown signal, shutting down all the clients")
	for _, client := range clients {
		client.Close()
	}

}

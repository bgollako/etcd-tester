package main

import (
	"context"
	"etcd-tester/client"
	"etcd-tester/clog"
	topics "etcd-tester/topics/generate"
	"etcd-tester/utils"
	"flag"
	"fmt"
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
	// Remove the last empty topic
	allTopics = allTopics[:len(allTopics)-1]

	totalTopics := len(allTopics)
	numTopicsPerClient := totalTopics / *numClients
	remainingTopics := totalTopics % *numClients

	logger.Sugar().Infof("Total no of topics - %d", totalTopics)
	logger.Sugar().Infof("Generating %d clients, with %d-%d topics per client", *numClients, numTopicsPerClient, numTopicsPerClient+1)

	var clients []client.Client
	start := 0
	for i := 0; i < *numClients; i++ {
		end := start + numTopicsPerClient
		if remainingTopics > 0 {
			end++
			remainingTopics--
		}
		client := client.NewClient(ctx, fmt.Sprintf("%s-%d", *clientName, i+1), allTopics[start:end], []string{"localhost:2379"})
		utils.Fatal(client.Start())
		clients = append(clients, client)
		start = end
	}
	logger.Info("Completed launching all the clients!!")

	shutdownCh := utils.GetShutdownHandle()
	<-shutdownCh

	logger.Info("Received shutdown signal, shutting down all the clients")
	for _, client := range clients {
		client.Close()
	}

}

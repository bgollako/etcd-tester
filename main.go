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
	"sort"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

func main2() {
	// Reading cmdline arguments
	numClients := flag.Int("numClients", 30, "number of etcd clients")
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
		client := client.NewClient(ctx, fmt.Sprintf("%s-%d", *clientName, i+1), allTopics[start:end], []string{"etcd:2379"})
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

func main() {
	generateLatencyStats(10000)
}

func generateLatencyStats(numKeys int) {
	ctx := clog.NewContextWithDefaultLogger(context.Background())
	logger := clog.MustFromContext(ctx)

	// Create keys
	keys := make([]string, numKeys)
	for i := range keys {
		keys[i] = fmt.Sprintf("key-%d", i)
	}

	// Assuming we have a client instance
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logger.Error("error while creating etcd client", zap.Error(err))
		return
	}
	defer etcdClient.Close()

	var writeLatencies []int64
	for _, key := range keys {
		startWrite := time.Now()
		_, err := etcdClient.Put(ctx, key, "value")
		if err != nil {
			logger.Error("error while writing to etcd", zap.Error(err))
			return
		}
		writeLatencies = append(writeLatencies, time.Since(startWrite).Nanoseconds())
	}
	totalWriteLatency := int64(0)
	for _, latency := range writeLatencies {
		totalWriteLatency += latency
	}
	avgWriteLatency := float64(totalWriteLatency) / float64(len(writeLatencies))
	writeLatencyPercentiles := getPercentiles(writeLatencies, []float64{0.5, 0.9, 0.9999})

	fmt.Printf("Total write latency: %d ns\n", totalWriteLatency)
	fmt.Printf("Average write latency: %.3f ns\n", avgWriteLatency)
	fmt.Printf("50th percentile write latency: %d ns\n", writeLatencyPercentiles[0])
	fmt.Printf("90th percentile write latency: %d ns\n", writeLatencyPercentiles[1])
	fmt.Printf("99.99th percentile write latency: %d ns\n", writeLatencyPercentiles[2])

	var readLatencies []int64
	for _, key := range keys {
		startRead := time.Now()
		_, err := etcdClient.Get(ctx, key)
		if err != nil {
			logger.Error("error while reading from etcd", zap.Error(err))
			return
		}
		readLatencies = append(readLatencies, time.Since(startRead).Nanoseconds())
	}
	totalReadLatency := int64(0)
	for _, latency := range readLatencies {
		totalReadLatency += latency
	}
	avgReadLatency := float64(totalReadLatency) / float64(len(readLatencies))
	readLatencyPercentiles := getPercentiles(readLatencies, []float64{0.5, 0.9, 0.9999})

	fmt.Printf("Total read latency: %d ns\n", totalReadLatency)
	fmt.Printf("Average read latency: %.3f ns\n", avgReadLatency)
	fmt.Printf("50th percentile read latency: %d ns\n", readLatencyPercentiles[0])
	fmt.Printf("90th percentile read latency: %d ns\n", readLatencyPercentiles[1])
	fmt.Printf("99.99th percentile read latency: %d ns\n", readLatencyPercentiles[2])

	var deleteLatencies []int64
	for _, key := range keys {
		startDelete := time.Now()
		_, err := etcdClient.Delete(ctx, key)
		if err != nil {
			logger.Error("error while deleting from etcd", zap.Error(err))
			return
		}
		deleteLatencies = append(deleteLatencies, time.Since(startDelete).Nanoseconds())
	}
	totalDeleteLatency := int64(0)
	for _, latency := range deleteLatencies {
		totalDeleteLatency += latency
	}
	avgDeleteLatency := float64(totalDeleteLatency) / float64(len(deleteLatencies))
	deleteLatencyPercentiles := getPercentiles(deleteLatencies, []float64{0.5, 0.9, 0.9999})

	fmt.Printf("Total delete latency: %d ns\n", totalDeleteLatency)
	fmt.Printf("Average delete latency: %.3f ns\n", avgDeleteLatency)
	fmt.Printf("50th percentile delete latency: %d ns\n", deleteLatencyPercentiles[0])
	fmt.Printf("90th percentile delete latency: %d ns\n", deleteLatencyPercentiles[1])
	fmt.Printf("99.99th percentile delete latency: %d ns\n", deleteLatencyPercentiles[2])
}

func getPercentiles(latencies []int64, percentiles []float64) []int64 {
	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})
	var percentilesDurations []int64
	for _, percentile := range percentiles {
		index := int(float64(len(latencies)) * percentile)
		percentilesDurations = append(percentilesDurations, latencies[index])
	}
	return percentilesDurations
}

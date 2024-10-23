package main

import (
	topics "etcd-tester/topics/generate"
	"flag"
)

func main() {
	num := flag.Int("count", 1, "no of topics to generate")
	flag.Parse()

	topics.GenerateTopics(*num, topics.Prefix, topics.TopicsFile)
}

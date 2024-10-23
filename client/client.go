package client

import (
	"context"
	"etcd-tester/clog"
	"time"

	"math/rand/v2"

	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/zap"
)

type Client interface {
	Start() error
	Close()
}

func NewClient(ctx context.Context, id string, topics, endpoints []string) Client {
	return &client{
		id:        id,
		ctx:       ctx,
		topics:    topics,
		endpoints: endpoints,
		rng:       rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(uuid.New().ID()))),
	}
}

type client struct {
	id        string
	ctx       context.Context
	cancel    func()
	topics    []string
	endpoints []string
	c         *clientv3.Client
	rng       *rand.Rand
}

func (c *client) Start() error {
	var err error

	// Adding client id to the logger
	log := clog.MustFromContext(c.ctx)
	log = log.With(zap.String(clog.KeyClientId, c.id))
	c.ctx = clog.NewContextWithLogger(c.ctx, log)

	// Creating a etcd client.
	// This client will be used for leader election on all the topics
	c.c, err = clientv3.New(clientv3.Config{
		Endpoints:   c.endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Error("error while creating etcd client", zap.Error(err))
		return err
	}

	// Launch routines to contest leader elections for given topics
	c.ctx, c.cancel = context.WithCancel(c.ctx)
	for _, topic := range c.topics {
		go c.keepContesting(topic)
	}
	return nil
}

// Shutsdowns the client
// Calls the cancel function to propagate the cancellation to all the topics
func (c *client) Close() {
	if c.cancel != nil {
		c.cancel()
	}
	if c.c != nil {
		if err := c.c.Close(); err != nil {
			clog.MustFromContext(c.ctx).Error("error shutting down client", zap.Error(err))
		}
	}
}

// keepContesting - contests the elections on the given
// topic endlessly, till it recieves a shutdown signal on the
// context.
func (c *client) keepContesting(topic string) {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-c.ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			c.contest(topic)
		}
	}
}

// contest - contests the elections on the given topic
func (c *client) contest(topic string) {
	var err error
	log := clog.MustFromContext(c.ctx)
	// Creating a new session
	var session *concurrency.Session
	session, err = concurrency.NewSession(c.c, concurrency.WithTTL(1))
	if err != nil {
		log.Error("error while creating etcd session", zap.Error(err))
		return
	}

	// Campaigning for a new election
	election := concurrency.NewElection(session, topic)
	log.Sugar().Infof("contesting election for topic %s", topic)
	// Blocks till it wins the election
	err = election.Campaign(c.ctx, uuid.NewString())
	if err != nil {
		log.Error("error while creating campaigning for election", zap.Error(err))
		return
	}
	// Yay!! we won the election
	log.Sugar().Infof("won election for topic %s", topic)
	// Sleep for a random duration between 30 to 60 seconds
	sleepDuration := time.Duration(30+c.rng.IntN(31)) * time.Second
	log.Sugar().Infof("sleeping for %v seconds as leader for topic %s", sleepDuration.Seconds(), topic)
	time.Sleep(sleepDuration)

	// Resigning from the election
	log.Sugar().Infof("resigning from leadership for topic %s", topic)
	err = election.Resign(c.ctx)
	if err != nil {
		log.Error("error while creating campaigning for election", zap.Error(err))
		return
	}

	// Close the session
	session.Close()
}

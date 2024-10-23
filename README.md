Each tester instance will have 100 etcd clients
Each client will participate in leader election of 500 predefined topics, one topic per session.
So total of 50k topics will be distributed among 100 clients.
Each leader will relinquish its leadership of its topic at a randomly chosen interval between 30-60 seconds, forcing a leadership election once again.

We will have 30 such tester instances, each with a unique name

All of the above numbers can be taken as a command line arguments.

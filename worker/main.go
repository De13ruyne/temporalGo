package main

import (
	"log"

	"my-org/greeting"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {

	c, err := client.Dial(client.Options{
		// HostPort: "172.18.0.4:7233",
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "my-task-queue", worker.Options{})

	// w2 := worker.New(c, "my-query-queue", worker.Options{})

	w.RegisterWorkflow(greeting.SayHelloWorkflow)
	w.RegisterActivity(greeting.Greet)

	w.RegisterWorkflow(greeting.QueryWorkflow)

	err = w.Run(worker.InterruptCh())
	// err = w2.Run(worker.InterruptCh())

	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}

}

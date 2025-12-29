package query

import (
	"context"
	"encoding/json"
	"log"

	greeting "my-org/greeting"

	"go.temporal.io/sdk/client"
)

func Query(jobID string) (interface{}, error) {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "q-wf",
		TaskQueue: "my-task-queue",
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, greeting.QueryWorkflow, jobID)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	var result interface{}
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	res, err := json.Marshal(result)
	if err != nil {
		log.Fatalln("json marshal error", err)
	}
	log.Println("Workflow result:", string(res))
	return result, nil
}

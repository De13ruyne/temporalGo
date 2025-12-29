package start

import (
	"context"
	"fmt"
	"log"

	// "os"

	greeting "my-org/greeting"

	"go.temporal.io/sdk/client"
)

func Start(para1 string, para2 string) {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "greeting-workflow",
		TaskQueue: "my-task-queue",
	}

	// para1 := os.Args[1]
	// para2 := os.Args[2]

	// retry 3 times
	for i := 0; i < 2; i++ {
		// para1 随便一个字符串
		// para2 设置是否第一次失败
		we, err := c.ExecuteWorkflow(context.Background(), options, greeting.SayHelloWorkflow, para1, para2, i+1)

		if para2 == "1" {
			para2 = "0"
		}
		if i > 0 {
			fmt.Println("Retry: ", i)
		}

		if err != nil {
			// log.Fatalln("Unable to execute workflow", err)
			log.Println("Unable to execute workflow", err)
			continue
		}
		log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

		var result string
		err = we.Get(context.Background(), &result)
		if err != nil {
			// log.Fatalln("Unable get workflow result", err)
			log.Println("Unable get workflow result", err)
			continue
		}
		log.Println("Workflow result:", result)
		break

	}

}

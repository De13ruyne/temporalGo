package greeting

import (
	"errors"
	"log/slog"
	"sync"
	"time"

	"go.temporal.io/sdk/workflow"
)

type jobmsg struct {
	JobID   string      `json:"job_id"`
	Status  string      `json:"status"`
	Attempt int         `json:"attempt"`
	Result  interface{} `json:"result"`
	Error   error       `json:"error"`
}

type safejobHash struct {
	mu   sync.RWMutex
	data map[string]jobmsg
}

var jobHash = safejobHash{
	data: make(map[string]jobmsg),
}

type safectxMap struct {
	mu   sync.RWMutex
	data map[string]workflow.Context
}

var ctxMap = safectxMap{
	data: make(map[string]workflow.Context),
}

// JobHash的写操作（Set）：加写锁
func (jh *safejobHash) Set(key string, val jobmsg) {
	jh.mu.Lock()
	defer jh.mu.Unlock()
	jh.data[key] = val
}

// JobHash的读操作（Get）：加读锁
func (jh *safejobHash) Get(key string) (jobmsg, bool) {
	jh.mu.RLock() // 读操作加读锁（可多个读锁并发）
	defer jh.mu.RUnlock()
	val, ok := jh.data[key]
	return val, ok
}

func (cm *safectxMap) Set(key string, val workflow.Context) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.data[key] = val
}

func (cm *safectxMap) Get(key string) (workflow.Context, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	val, ok := cm.data[key]
	return val, ok
}

// var jobHash = map[string]jobmsg{}

// var ctxMap = map[string]workflow.Context{}

func SayHelloWorkflow(ctx workflow.Context, name string, failedAttempts string, times int) (string, error) {
	println(workflow.GetInfo(ctx).WorkflowStartTime.Format("2006-01-02 15:04:05"))
	if failedAttempts == "1" {
		err := errors.New("simulated failure")
		// workflow.GetInfo(ctx).Attempt++
		println("retry!!!!!!!!!!!!!!")
		return "", err
	}
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 10,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	info := workflow.GetInfo(ctx)

	// println("ctx1: ", ctx)
	// workflow.GetLogger(ctx).Debug("")
	jobtmp := jobmsg{
		JobID:   info.WorkflowExecution.ID,
		Status:  "running",
		Attempt: times,
		Result:  nil,
		Error:   nil,
	}

	// jobHash.data[info.WorkflowExecution.ID] = jobtmp

	jobHash.Set(info.WorkflowExecution.ID, jobtmp)

	// ctxMap[info.WorkflowExecution.ID] = ctx
	ctxMap.Set(info.WorkflowExecution.ID, ctx)

	var result string
	err := workflow.ExecuteActivity(ctx, Greet, name).Get(ctx, &result)
	if err != nil {
		return "", err
	}
	println("ctx2: ", ctx)

	// // Getting the logger from the context.
	workflow.GetLogger(ctx).Debug("")

	info = workflow.GetInfo(ctx)
	println("info: ", info.WorkflowExecution.ID)

	jobtmp = jobmsg{
		JobID:   info.WorkflowExecution.ID,
		Status:  "finished",
		Attempt: times,
		Result:  result,
		Error:   nil,
	}

	// jobHash[info.WorkflowExecution.ID] = jobtmp
	jobHash.Set(info.WorkflowExecution.ID, jobtmp)

	workflow.GetLogger(ctx).Debug("")
	slog.Debug("hhhhhc")
	return result, nil
}

func QueryWorkflow(ctx workflow.Context, jobID string) (jobmsg, error) {

	// info := workflow.GetInfo(ctx)
	println("---------------------------")
	// println(jobHash[jobID].Attempt)
	// println(ctxMap[jobID])
	// println((workflow.GetInfo(ctxMap[jobID]).WorkflowStartTime.Format("2006-01-02 15:04:05")))
	// println((workflow.GetInfo(ctxMap[jobID]).Attempt))
	println("---------------------------")

	// println("Query info: ", info.WorkflowExecution.ID)

	// 返回预定义的jobMap
	// return jobHash[jobID], nil
	if jobret, ok := jobHash.Get(jobID); ok {
		return jobret, nil
	}
	return jobmsg{}, errors.New("job ID not found")
}

// func MyWorkflow(ctx workflow.Context, input string) error {
// 	ao := workflow.ActivityOptions{
// 		StartToCloseTimeout: time.Second * 10,
// 	}
// 	currentState := "started" // this could be any serializable struct
// 	err := workflow.SetQueryHandler(ctx, "current_state", func() (string, error) {
// 		return currentState, nil
// 	})
// 	if err != nil {
// 		currentState = "failed to register query handler"
// 		return err
// 	}
// 	// your normal workflow code begins here, and you update the currentState as the code makes progress.
// 	currentState = "waiting timer"
// 	err = workflow.NewTimer(ctx, time.Hour).Get(ctx, nil)
// 	if err != nil {
// 		currentState = "timer failed"
// 		return err
// 	}

// 	currentState = "waiting activity"
// 	ctx = workflow.WithActivityOptions(ctx, ao)
// 	err = workflow.ExecuteActivity(ctx, Greet, "my_input").Get(ctx, nil)
// 	if err != nil {
// 		currentState = "activity failed"
// 		return err
// 	}
// 	currentState = "done"
// 	return nil
// }

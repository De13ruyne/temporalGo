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

// JobHash的写操作：加写锁
func (jh *safejobHash) Set(key string, val jobmsg) {
	jh.mu.Lock()
	defer jh.mu.Unlock()
	jh.data[key] = val
}

// JobHash的读操作：加读锁
func (jh *safejobHash) Get(key string) (jobmsg, bool) {
	jh.mu.RLock() // 读操作加读锁 可多个读锁并发
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

	jobtmp := jobmsg{
		JobID:   info.WorkflowExecution.ID,
		Status:  "running",
		Attempt: times,
		Result:  nil,
		Error:   nil,
	}

	jobHash.Set(info.WorkflowExecution.ID, jobtmp)

	ctxMap.Set(info.WorkflowExecution.ID, ctx)

	var result string
	err := workflow.ExecuteActivity(ctx, Greet, name).Get(ctx, &result)
	if err != nil {
		return "", err
	}
	println("ctx2: ", ctx)

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

	jobHash.Set(info.WorkflowExecution.ID, jobtmp)

	workflow.GetLogger(ctx).Debug("")
	slog.Debug("hhhhhc")
	return result, nil
}

func QueryWorkflow(ctx workflow.Context, jobID string) (jobmsg, error) {

	// 返回预定义的jobMap
	// return jobHash[jobID], nil
	if jobret, ok := jobHash.Get(jobID); ok {
		return jobret, nil
	}
	return jobmsg{}, errors.New("job ID not found")
}

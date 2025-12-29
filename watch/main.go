package main

import (
	query "my-org/greeting/query"
	start "my-org/greeting/start"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 对应POST /jobs的请求体结构
type Input struct {
	Numbers []int `json:"numbers"` // 绑定请求体的input.numbers
}

type Options struct {
	FailFirstAttempt bool `json:"fail_first_attempt"` // 绑定请求体的options.fail_first_attempt
}

type JobRequest struct {
	Input   string `json:"input"`
	Options string `json:"options"`
}

func startJobHandler(c *gin.Context) {
	var req JobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体格式错误: " + err.Error()})
		return
	}
	start.Start(req.Input, req.Options)

	c.JSON(http.StatusOK, gin.H{"job_id": "greeting-workflow"})
}

func queryJobHandler(c *gin.Context) {
	jobID := c.Param("job_id")

	res, err := query.Query(jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询任务状态失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)

}

func main() {
	// 初始化Gin引擎
	r := gin.Default()

	// 注册接口路由
	r.POST("/jobs", startJobHandler)        // 启动任务
	r.GET("/jobs/:job_id", queryJobHandler) // 查询任务状态

	// 启动服务（端口8080）
	if err := r.Run(":8080"); err != nil {
		panic("服务启动失败: " + err.Error())
	}
}

package greeting

import (
	"context"
	"fmt"
	"time"
)

func Greet(ctx context.Context, name string) (string, error) {
	// 模拟耗时操作
	time.Sleep(5 * time.Second)
	return fmt.Sprintf("Hello %s", name), nil
}

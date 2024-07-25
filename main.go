package main

import (
	"context"

	"data-backup/cmd"
	"github.com/ihezebin/oneness/logger"
)

func main() {
	ctx := context.Background()
	if err := cmd.Run(ctx); err != nil {
		logger.Fatalf(ctx, "cmd run error: %v", err)
	}
}

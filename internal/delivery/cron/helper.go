package cron

import (
	"log"
	"runtime/debug"
)

func safeExecute(logger *log.Logger, jobName string, fn func()) {
	defer func() {
		if err := recover(); err != nil {
			logger.Printf("[CRON PANIC RECOVERED] Job: %s | Error: %v\nStack trace:\n%s", jobName, err, debug.Stack())
		}
	}()
	fn()
}

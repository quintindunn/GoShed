package telemetry

import (
	"fmt"
	"time"
)

func Log(msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s | %s\n", timestamp, msg)
}

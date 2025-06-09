package cron

import (
	"com.quintindev/APIShed/util"
	"time"
)

const nullificationCooldown = 10 * time.Second

func Manager() {
	ticker := time.NewTicker(nullificationCooldown)
	defer ticker.Stop()

	for {
		go codeNullificationCron()
		<-ticker.C
	}
}

func codeNullificationCron() {
	util.NullifyAllocatedCodes()
	util.UpdateExpiredRollingCodes()
}

package cron

import (
	"com.quintindev/APIShed/util"
	"time"
)

func Manager() {
	var nullificationCooldown = util.QueryConfigValue[int64]("code_expiration_check_interval")
	ticker := time.NewTicker(time.Duration(nullificationCooldown) * time.Millisecond)
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

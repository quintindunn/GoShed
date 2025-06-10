package hardware

import (
	"com.quintindev/APIShed/audit"
	"com.quintindev/APIShed/util"
	"fmt"
	"time"
)

func SetLockState(newState bool) {
	state := "UNARMED"
	if newState {
		state = "ARMED"
	}
	util.SetConfigValue[bool]("lock_state", newState)
	audit.LogInitiator("SYSTEM", fmt.Sprintf("Setting lock state to %s", state))
}

func HandleCodedUnlock() {
	util.NullifyAllocatedCodes()
	util.UpdateExpiredRollingCodes()

	SetLockState(false)
	time.Sleep(time.Duration(util.QueryConfigValue[int64]("unlock_time")) * time.Millisecond)
	SetLockState(true)
}

func GetArmedState() bool {
	return util.QueryConfigValue[bool]("lock_state")
}

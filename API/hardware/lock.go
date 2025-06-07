package hardware

import (
	"com.quintindev/APIShed/audit"
	"com.quintindev/APIShed/util"
	"fmt"
	"time"
)

type LockHardwareState struct {
	Locked bool
}

var LockState = LockHardwareState{}

func SetLockState(newState bool) {
	LockState.Locked = newState
	state := "UNLOCKED"
	if newState {
		state = "LOCKED"
	}
	audit.LogInitiator("SYSTEM", fmt.Sprintf("Setting lock state to %s", state))
}

func HandleCodedUnlock() {
	util.NullifyAllocatedCodes()
	util.UpdateExpiredRollingCodes()

	SetLockState(false)
	time.Sleep(time.Duration(util.QueryConfigValue[int64]("unlock_time")) * time.Millisecond)
	SetLockState(true)
}

func GetLockedState() bool {
	return LockState.Locked
}

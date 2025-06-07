package hardware

import (
	"com.quintindev/APIShed/audit"
	"fmt"
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

func GetLockedState() bool {
	return LockState.Locked
}

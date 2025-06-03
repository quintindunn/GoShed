package hardware

import (
	"com.quintindev/APIShed/telemetry"
	"fmt"
)

type LockHardwareState struct {
	Locked bool
}

var LockState = LockHardwareState{}

func SetLockState(newState bool) {
	LockState.Locked = newState
	telemetry.Log(fmt.Sprintf("Setting lock state to %v", newState))
}

func GetLockedState() bool {
	return LockState.Locked
}

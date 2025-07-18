package audit

import (
	"com.quintindev/APIShed/database"
	"com.quintindev/APIShed/models"
	"fmt"
	"time"
)

func Log(msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s | %s\n", timestamp, msg)

	log := models.Log{
		Msg: msg,
	}

	if err := database.DB.Create(&log).Error; err != nil {
		fmt.Println("Error creating log:", err)
	}
}

func LogInitiator(initiator string, msg string) {
	Log(fmt.Sprintf("%s | %s", initiator, msg))
}

func NullifyRollingCodes(codes []string) {
	msg := fmt.Sprintf("Nullified %d rolling codes: %+v", len(codes), codes)
	LogInitiator("SYSTEM", msg)
}

func CreateNewRollingCodes(codes []string) {
	msg := fmt.Sprintf("Created %d rolling codes: %+v", len(codes), codes)
	LogInitiator("SYSTEM", msg)
}

func NullifyAllocatedCodes(codes [][]string) {
	msg := fmt.Sprintf("Nullified %d allocated codes: %+v", len(codes), codes)
	LogInitiator("SYSTEM", msg)
}

func UnlockByAllocatedCode(code models.AllocatedCode, initiator string) {
	msg := fmt.Sprintf("Unlocking from user \"%s\", with code: \"%s\"", code.Name, code.Code)
	LogInitiator(initiator, msg)
}

func UnlockByRollingCode(code models.RollingCode, initiator string) {
	msg := fmt.Sprintf("Unlocking from rolling code: \"%s\"", code.Code)
	LogInitiator(initiator, msg)
}

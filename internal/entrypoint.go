package internal

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"time"
)

const (
	// TODO move these to env variables

	EntriesToGrab = 10
	lowThreshold  = 5.0
	highThreshold = 10.0
	urgentLow     = 3.0
	urgentHigh    = 15.0
)

func Entrypoint() {
	logger := GetLogger()
	zap.ReplaceGlobals(logger)

	logger.Info("Starting KateAlerts")

	xdripDataGrabber := NewXDropDataGrabber(os.Getenv("xdripURL"), os.Getenv("xdripToken"))
	ifttt := NewIFTTT(os.Getenv("iftttURL"))

	entries := xdripDataGrabber.Entries(EntriesToGrab)
	currentLevel := MgdlToMmol(entries[0].Sgv)

	// Check if data is recent
	stale := entries[0].SysTime.Before(time.Now().Add(time.Duration(-10) * time.Minute))

	if stale {
		difference := time.Now().Sub(entries[0].SysTime)
		ifttt.TriggerEndpoint("NOT READ", fmt.Sprintf(
			"CGM data stale, last recorded value %2.fmmol/L. %2.f minutes ago", currentLevel, difference.Minutes()), "")
		return
	}

	delta := MgdlToMmol(entries[0].Sgv - entries[EntriesToGrab-1].Sgv)
	deltaText := fmt.Sprintf("%.2f", delta)
	sendWarning := true
	value1 := ""
	if currentLevel <= urgentLow {
		value1 = "URGENT LOW"
	} else if currentLevel > urgentLow && currentLevel <= lowThreshold {
		value1 = "Low"
	} else if currentLevel >= urgentHigh {
		value1 = "URGENT HIGH"
	} else if currentLevel < urgentHigh && currentLevel >= highThreshold {
		value1 = "High"
	} else {
		sendWarning = false
	}

	if sendWarning {
		if delta > 0 {
			deltaText = "+" + deltaText
		}
		ifttt.TriggerEndpoint("NOT READ", fmt.Sprintf("CGM Reports Kate is %s: %.2f Δ%smmol/L", value1, currentLevel, deltaText), "")
	}
}

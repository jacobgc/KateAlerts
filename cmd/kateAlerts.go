package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"jacobgc.me/katealert/internal"
	"os"
)

const (
	// TODO move these to env variables

	EntriesToGrab = 5
	lowThreshold  = 5.0
	highThreshold = 10.0
	urgentLow     = 3.0
	urgentHigh    = 15.0
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// HelloPubSub consumes a Pub/Sub message.
func HelloPubSub(ctx context.Context, m PubSubMessage) error {
	main()
	return nil
}

func main() {

	logger := internal.GetLogger()
	zap.ReplaceGlobals(logger)

	logger.Info("Starting KateAlerts")

	xdripDataGrabber := internal.NewXDropDataGrabber(os.Getenv("xdripURL"))
	ifttt := internal.NewIFTTT(os.Getenv("iftttURL"))

	entries := xdripDataGrabber.Entries(EntriesToGrab)

	currentLevel := internal.MgdlToMmol(entries[0].Sgv)
	delta := internal.MgdlToMmol(entries[0].Sgv - entries[EntriesToGrab-1].Sgv)
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
		ifttt.TriggerEndpoint("Kate Not In Range", fmt.Sprintf("CGM Reports Kate is %s: %.2f Δ%smmol/L", value1, currentLevel, deltaText), "")
	}

	println(entries[0].Sgv)
	fmt.Printf("%.2f", internal.MgdlToMmol(entries[0].Sgv))
}

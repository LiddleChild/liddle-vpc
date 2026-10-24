package timeutil

import (
	"fmt"
	"time"
)

var (
	LocationBangkok *time.Location
)

func init() {
	var err error
	LocationBangkok, err = time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(fmt.Errorf("failed to load location: %w", err))
	}
}

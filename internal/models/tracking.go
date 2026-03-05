package models

import "time"

type Tracking struct {
	Track     string
	URI       string
	Method    string
	Request   interface{}
	Response  interface{}
	StartDate time.Time
}

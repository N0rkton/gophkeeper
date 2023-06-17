package datamodels

import "time"

type Auth struct {
	ID       uint32
	Password string
}
type Data struct {
	UserID    uint32
	DataID    uint32
	Data      string
	Metadata  string
	ChangedAt time.Time
}

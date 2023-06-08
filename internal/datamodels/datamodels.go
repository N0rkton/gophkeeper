// Package datamodels represents structs used in program
package datamodels

import "time"

// Auth - struct used to save info about new user
type Auth struct {
	ID       uint32 `json:"ID"`
	Login    string `json:"Login"`
	Password string `json:"Password"`
}

// Login - struct for login
type Login struct {
	ID       uint32 `json:"ID"`
	Password string `json:"Password"`
}

// Data - struct for all information about 1 note
type Data struct {
	UserID    uint32    `json:"UserID"`
	DataID    string    `json:"DataID"`
	Data      string    `json:"Data"`
	Metadata  string    `json:"Metadata"`
	ChangedAt time.Time `json:"ChangedAt"`
	Deleted   bool      `json:"Deleted"`
}

// UniqueData - unique constraint from database for in memory storage
type UniqueData struct {
	DataID string
	UserID uint32
}

package models

import (
	"gorm.io/gorm"
)

// User structure
type Revision struct {
	gorm.Model
	MemID   uint `json:"mem_id" example:"1"`
	Mem     Mem
	Result  uint `json:"result_int" example:"0"` // 0 means false, 1 means true ! This should fix Metabase issue
	Quality uint `json:"quality" example:"0"`    // [0: Blackout - 1: Error with choices - 2: Error with hints - 3: Error - 4: Good with hints - 5: Perfect]
}

package api

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/ChaunceyXCX/OpenTools/internal/core/database"
)

const winStateKey = "ZTOOLS/window_state"

type WindowState struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

func SaveWindowState(state *WindowState) {
	if globalDB == nil {
		return
	}
	data, err := json.Marshal(state)
	if err != nil {
		log.Printf("[Window] marshal state error: %v", err)
		return
	}
	dataStr := string(data)
	globalDB.Put(&database.DbDoc{
		ID:   winStateKey,
		Data: dataStr,
	})
}

func LoadWindowState() *WindowState {
	if globalDB == nil {
		return nil
	}
	doc, err := globalDB.Get(winStateKey)
	if err != nil || doc == nil {
		return nil
	}
	var state WindowState
	if err := json.Unmarshal([]byte(doc.Data), &state); err != nil {
		log.Printf("[Window] unmarshal state error: %v", err)
		return nil
	}
	return &state
}

func (s *WindowState) String() string {
	return strconv.Itoa(s.X) + "," + strconv.Itoa(s.Y) + "," + strconv.Itoa(s.Width) + "," + strconv.Itoa(s.Height)
}

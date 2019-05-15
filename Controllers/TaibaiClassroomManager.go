package Controllers

import (
	"TaiBaiSupport/Models"
	"sync"
)


type TaibaiClassroomManager struct {
	OperationRWMux sync.RWMutex
	WSConns        map[int]map[int]*Models.TaibaiWSConn
}

func NewTaibaiClassroomManager() *TaibaiClassroomManager {
	M := &TaibaiClassroomManager{
		WSConns:       make(map[int]map[int]*Models.TaibaiWSConn),
	}
	return M
}

var TaibaiClassroomManagerInstance *TaibaiClassroomManager

func init() {
	TaibaiClassroomManagerInstance = NewTaibaiClassroomManager()
}

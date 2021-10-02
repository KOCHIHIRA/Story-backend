package websocket

type Manager struct {
	RoomList map[string]*Room
}

func NewManager() Manager {
	return Manager{
		RoomList: make(map[string]*Room),
	}
}

//新たにルームを作成する処理
func (manager Manager) CreateRoom(roomName, roomTitle string) (bool, *Room) {
	if _, ok := manager.RoomList[roomName]; ok {
		return false, nil
	}
	room := NewRoom(roomName, roomTitle)
	manager.RoomList[roomName] = room
	return true, room
}

//ルームを読み込む処理
func (manager Manager) ReadRoom(roomName string) (*Room, bool) {
	if room, err := manager.RoomList[roomName]; err {
		return room, true
	}
	return nil, false
}

//ルームを削除する処理
func (manager Manager) DeleteRoom(roomName string) {
	if _, ok := manager.RoomList[roomName]; ok {
		delete(manager.RoomList, roomName)
	}
}

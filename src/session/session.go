package session

import (
	"container/list"
	"time"
)

var pder = &Provider{list: list.New()}

func init() {
	pder.sessions = make(map[string]*list.Element, 0)
}

type Session struct {
	sessionID  string
	isLogin    bool
	lastAccess time.Time
	value      map[interface{}]interface{}
}

//セッション領域から値を取得して、アクセス時間を更新している。
func (session *Session) Get(key interface{}) interface{} {
	//pder.SessionUpdate(session.sessionID)
	if v, ok := session.value[key]; ok {
		return v
	} else {
		return nil
	}
}

//セッション領域に値を入れて、アクセス時間を更新している。
func (session *Session) Set(key, value interface{}) error {
	session.value[key] = value
	//pder.SessionUpdate(session.sessionID)
	return nil
}

//valueの中の特定の要素を削除して、アクセスした時間を更新している。
func (session *Session) Delete(key string) error {
	delete(session.value, key)
	//pder.SessionUpdate(session.sessionID)
	return nil
}

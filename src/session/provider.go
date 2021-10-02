package session

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type Provider struct {
	lock     sync.Mutex
	sessions map[string]*list.Element
	list     *list.List
}

func (prv *Provider) NewProvider() *Provider {
	return &Provider{list: list.New()}
}

func (prv *Provider) SessionInit(sid string) (*Session, error) {
	prv.lock.Lock()
	defer prv.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	newsess := &Session{sessionID: sid, lastAccess: time.Now(), value: v, isLogin: true}
	element := prv.list.PushBack(newsess)
	prv.sessions[sid] = element
	return newsess, nil
}

/*
func (pder *Provider) SessionRead(sid string) (*Session, error) {
	if element, ok := pder.sessions[sid]; ok {
		return element.Value.(*Session), nil
	} else {
		sess, err := pder.SessionInit(sid)
		return sess, err
	}
	//return nil, nil
}
*/

func (prv *Provider) SessionRead(sid string) (*Session, error) {
	prv.lock.Lock()
	defer prv.lock.Unlock()
	if element, ok := pder.sessions[sid]; ok {
		element.Value.(*Session).lastAccess = time.Now()
		//時間を更新した後、要素をリストの先頭に持ってきている。
		pder.list.MoveToFront(element)
		return element.Value.(*Session), nil
	}
	return nil, fmt.Errorf("session_out")
}

func (pder *Provider) SessionDestroy(sid string) {
	if element, ok := pder.sessions[sid]; ok {
		delete(pder.sessions, sid)
		pder.list.Remove(element)
		return
	}
	return
}

//一定時間アクセスしていない期限切れのセッションを破棄してる。
func (pder *Provider) SessionGC(maxlifetime int64) {
	pder.lock.Lock()
	defer pder.lock.Unlock()

	for {
		//リストの後方の要素を一つづつ取り出している。
		element := pder.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*Session).lastAccess.Unix() + maxlifetime) < time.Now().Unix() {
			pder.list.Remove(element)
			delete(pder.sessions, element.Value.(*Session).sessionID)
		} else {
			break
		}
	}
}

//アクセスしてきた時間を現在の時間に更新している。
func (pder *Provider) SessionUpdate(sid string) error {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	if element, ok := pder.sessions[sid]; ok {
		element.Value.(*Session).lastAccess = time.Now()
		//時間を更新した後、要素をリストの先頭に持ってきている。
		pder.list.MoveToFront(element)
		return nil
	} else {
		return fmt.Errorf("session_out")
	}
}

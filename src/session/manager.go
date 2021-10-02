package session

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"sync"
	"time"
)

type Manager struct {
	CookieName  string
	Provider    *Provider
	Lock        sync.Mutex
	Maxlifetime int64
}

func NewManager() (*Manager, error) {

	manager := &Manager{CookieName: "SessionID", Provider: pder, Maxlifetime: 3600}
	return manager, nil
}

func (mng *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session *Session) {
	mng.Lock.Lock()
	defer mng.Lock.Unlock()
	//セッション切れはクッキーにセッションIDはあるがsessions[]にセッションIDで検索かけてもインスタンスが取れない
	sid := mng.sessionId()
	session, _ = mng.Provider.SessionInit(sid)
	cookie := http.Cookie{Name: mng.CookieName, Value: sid, Path: "/", HttpOnly: true}
	http.SetCookie(w, &cookie)
	return session
}

func (manager *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.CookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		manager.Lock.Lock()
		defer manager.Lock.Unlock()
		manager.Provider.SessionDestroy(cookie.Value)
		expire := time.Now()
		cookie := http.Cookie{Name: manager.CookieName, Path: "/", HttpOnly: true, Expires: expire, MaxAge: -1}
		http.SetCookie(w, &cookie)

	}
}

func (manager *Manager) GC() {
	manager.Lock.Lock()
	defer manager.Lock.Unlock()
	manager.Provider.SessionGC(manager.Maxlifetime)
	time.AfterFunc(time.Duration(manager.Maxlifetime), func() { manager.GC() })
}

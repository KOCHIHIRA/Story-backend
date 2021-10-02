package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/src/session"
	"server/src/sqls"
	"server/src/websocket"
)

var sessionManager *session.Manager
var roomManager websocket.Manager

const DBname = "STORY"

func init() {
	sessionManager, _ = session.NewManager()
	go sessionManager.GC()
	roomManager = websocket.NewManager()
}

type Response struct {
	Status       bool        `json:"status"`
	Data         interface{} `json:"data"`
	ErrorMessage string      `json:"error"`
}

type RequestRoom struct {
	Offset int `json:"offset"`
}

//cors対策の処理
func check_cors(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")
	return w
}

//ルームリストのリクエストを受けたときの処理
func RoomList(w http.ResponseWriter, r *http.Request) {
	var request RequestRoom
	var res Response
	w = check_cors(w)
	if r.Method != "OPTIONS" {
		cookie, err := r.Cookie(sessionManager.CookieName)
		if err != nil || cookie.Value == "" {
			res = Response{Status: false, ErrorMessage: "session_out"}
		} else {
			err := sessionManager.Provider.SessionUpdate(cookie.Value)
			if err != nil {
				res = Response{Status: false, ErrorMessage: "session_out"}
			} else {
				json.NewDecoder(r.Body).Decode(&request)
				errMess, result := sqls.GetAllRoomDetail("STORY", (request.Offset * 10))
				isGetData := true
				if result == nil {
					isGetData = false
				}
				res = Response{Status: isGetData, Data: result, ErrorMessage: errMess}
			}
		}

		response, err := json.Marshal(res)
		if err != nil {
			fmt.Println("エラーが出力されました。")
			log.Fatal(err)
			return
		}
		w.Write(response)
		return
	}
}

//ルーム作成のリクエストを受けたときの処理
func CreateRoom(w http.ResponseWriter, r *http.Request) {
	var res Response
	w = check_cors(w)
	if r.Method != "OPTIONS" {
		var room sqls.Room
		cookie, err := r.Cookie(sessionManager.CookieName)
		if err != nil || cookie.Value == "" {
			res = Response{Status: false, ErrorMessage: "session_out"}
		} else {
			session, err := sessionManager.Provider.SessionRead(cookie.Value)
			if err != nil {
				res = Response{Status: false, ErrorMessage: "session_out"}
			} else {
				json.NewDecoder(r.Body).Decode(&room)
				userName := session.Get("userName")
				result, errMess := sqls.CreateRoom("STORY", room.Name, room.Title, fmt.Sprint(userName))
				res = Response{Status: result, ErrorMessage: errMess}
				//チャットルームを作る

				if ok, chatRoom := roomManager.CreateRoom(room.Name, room.Title); ok {
					go chatRoom.Start()
				}
			}
		}
		response, err := json.Marshal(res)
		if err != nil {
			log.Fatal(err)
			fmt.Println("エラーが出力されました。")
			return
		}
		w.Write(response)
		return
	}
}

//ルームに参加したときの処理
func JoinRoom(w http.ResponseWriter, r *http.Request) {
	var res Response
	w = check_cors(w)
	if r.Method != "OPTIONS" {
		var room sqls.Room
		cookie, err := r.Cookie(sessionManager.CookieName)
		if err != nil || cookie.Value == "" {
			res = Response{Status: false, ErrorMessage: "session_out"}
		} else {
			session, err := sessionManager.Provider.SessionRead(cookie.Value)
			if err != nil {
				res = Response{Status: false, ErrorMessage: "session_out"}
			} else {
				json.NewDecoder(r.Body).Decode(&room)
				//チャットルームを探す
				_, err := roomManager.ReadRoom(room.Name)
				if !err {
					res = Response{Status: false, ErrorMessage: "部屋が存在しません"}
				} else {
					session.Set("roomName", room.Name)
					session.Set("roomTitle", room.Title)
					session.Set("roomOwner", room.Owner)
					session.Set("roomVote", room.Vote)

					res = Response{Status: true, ErrorMessage: ""}
				}
			}
		}
		response, err := json.Marshal(res)
		if err != nil {
			log.Fatal(err)
			fmt.Println("エラーが出力されました。")
			return
		}
		w.Write(response)
		return
	}
}

//ルーム参加
func ChatRoom(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionManager.CookieName)
	if err == nil {
		session, err := sessionManager.Provider.SessionRead(cookie.Value)
		if err == nil {

			userName := session.Get("userName")
			roomName := session.Get("roomName")
			roomTitle := session.Get("roomTitle")
			roomOwner := session.Get("roomOwner")
			roomVote := session.Get("roomVote")

			//チャットルームを探す
			chatRoom, err := roomManager.ReadRoom(roomName.(string))
			if err {

				//以前に書き込みしたかどうかを調べる
				conn, err := websocket.Upgrade(w, r)
				if err != nil {
					return
				}
				permission, _ := sqls.CheckPermission("STORY", roomName.(string), fmt.Sprint(userName))
				//storyの中身を並べ替えてクライアントに送信する
				status, messages := sqls.GetRoomSentence(DBname, roomName.(string))
				//roomに存在するユーザーの一覧を送信 する
				client := &websocket.Client{Name: userName.(string), Conn: conn, WritePermission: permission}
				if status {
					fmt.Println(len(messages))
					client.Conn.WriteJSON(websocket.Message{Type: "INIT_DATA",
						Users: chatRoom.GetUserList(), RoomName: roomName.(string), RoomTitle: roomTitle.(string),
						Owner: roomOwner.(string), Vote: roomVote.(int), Storys: messages})
				} else {
					client.Conn.WriteJSON(websocket.Message{Type: "INIT_DATA",
						RoomName: roomName.(string), RoomTitle: roomTitle.(string), Owner: roomOwner.(string),
						Users: chatRoom.GetUserList(), Vote: roomVote.(int),
						Storys: []sqls.Story{{UserName: "", Sentence: ""}}})
				}
				chatRoom.Register <- client
				client.Read(chatRoom)
			}
		}
	}
}

//投票の処理
func Voting(w http.ResponseWriter, r *http.Request) {
	var res Response
	fmt.Println("I get Voting request")
	w = check_cors(w)
	if r.Method != "OPTIONS" {
		var room sqls.Room
		cookie, err := r.Cookie(sessionManager.CookieName)
		if err != nil || cookie.Value == "" {
			res = Response{Status: false, ErrorMessage: "session_out"}
		} else {
			session, err := sessionManager.Provider.SessionRead(cookie.Value)
			if err != nil {
				res = Response{Status: false, ErrorMessage: "session_out"}
			} else {
				json.NewDecoder(r.Body).Decode(&room)
				userName := session.Get("userName")
				result, errMsg := sqls.Voting("STORY", userName.(string), room.Name)
				res = Response{Status: result, ErrorMessage: errMsg}
				//チャットルームを作る
			}
		}
		response, err := json.Marshal(res)
		if err != nil {
			log.Fatal(err)
			fmt.Println("エラーが出力されました。")
			return
		}
		w.Write(response)
		return
	}
}

//ランキングの処理
func Ranking(w http.ResponseWriter, r *http.Request) {
	var res Response
	w = check_cors(w)
	if r.Method != "OPTIONS" {
		cookie, err := r.Cookie(sessionManager.CookieName)
		if err != nil || cookie.Value == "" {
			res = Response{Status: false, ErrorMessage: "session_out"}
		} else {
			_, err := sessionManager.Provider.SessionRead(cookie.Value)
			if err != nil {
				res = Response{Status: false, ErrorMessage: "session_out"}
			} else {
				errMsg, result := sqls.GetRanking(DBname)
				fmt.Println(result)
				isGetData := true
				if result == nil {
					isGetData = false
				}
				res = Response{Status: isGetData, Data: result, ErrorMessage: errMsg}
				//チャットルームを作る
			}
		}
		response, err := json.Marshal(res)
		if err != nil {
			log.Fatal(err)
			fmt.Println("エラーが出力されました。")
			return
		}
		w.Write(response)
		return
	}
}

//ルーム検索
func Search(w http.ResponseWriter, r *http.Request) {
	var res Response
	w = check_cors(w)
	if r.Method != "OPTIONS" {
		var room sqls.Room
		cookie, err := r.Cookie(sessionManager.CookieName)
		if err != nil || cookie.Value == "" {
			res = Response{Status: false, ErrorMessage: "session_out"}
		} else {
			_, err := sessionManager.Provider.SessionRead(cookie.Value)
			if err != nil {
				res = Response{Status: false, ErrorMessage: "session_out"}
			} else {
				json.NewDecoder(r.Body).Decode(&room)
				errMsg, result := sqls.GetSearch(DBname, room.Name)
				isGetData := true
				if result == nil {
					isGetData = false
				}
				res = Response{Status: isGetData, Data: result, ErrorMessage: errMsg}
				//チャットルームを作る
			}
		}
		response, err := json.Marshal(res)
		if err != nil {
			log.Fatal(err)
			fmt.Println("エラーが出力されました。")
			return
		}
		w.Write(response)
		return
	}
}

//ユーザー登録リクエストを受け付けた時の処理
func SignUp(w http.ResponseWriter, r *http.Request) {
	w = check_cors(w)
	if r.Method != "OPTIONS" {
		var user sqls.User
		json.NewDecoder(r.Body).Decode(&user)
		result, errMsg := sqls.SignUp("STORY", user.Name, user.Password, user.Mail)
		res := Response{Status: result, ErrorMessage: errMsg}
		response, err := json.Marshal(res)
		if err != nil {
			log.Fatal(err)
			fmt.Println("エラーが出力されました。")
			return
		}
		w.Write(response)
		return
	}
}

func LogIn(w http.ResponseWriter, r *http.Request) {
	w = check_cors(w)
	if r.Method != "OPTIONS" {
		var user sqls.User
		var res Response
		json.NewDecoder(r.Body).Decode(&user)
		//ユーザーIDとパスワードでユーザー認証をする。
		result, errMsg := sqls.LogIn("STORY", user.Name, user.Password)
		if result {
			session := sessionManager.SessionStart(w, r)
			session.Set("userName", user.Name)
		}
		res = Response{Status: result, ErrorMessage: errMsg}

		response, err := json.Marshal(res)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(response)
		return
	}
}

//ルーム作成のリクエストを受けたときの処理
func LogOut(w http.ResponseWriter, r *http.Request) {
	w = check_cors(w)
	if r.Method != "OPTIONS" {
		sessionManager.SessionDestroy(w, r)
		res := Response{Status: true}
		response, err := json.Marshal(res)
		if err != nil {
			log.Fatal(err)
			fmt.Println("エラーが出力されました。")
			return
		}
		w.Write(response)
		return
	}
}

func main() {
	http.HandleFunc("/roomlist", RoomList)
	http.HandleFunc("/regist", SignUp)
	http.HandleFunc("/login", LogIn)
	http.HandleFunc("/logout", LogOut)
	http.HandleFunc("/create_room", CreateRoom)
	http.HandleFunc("/voting", Voting)
	http.HandleFunc("/ranking", Ranking)
	http.HandleFunc("/search", Search)
	http.HandleFunc("/roomjoin", JoinRoom)
	http.HandleFunc("/ws", ChatRoom)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":3001", nil); err != nil {
		fmt.Println(err)
	}
}

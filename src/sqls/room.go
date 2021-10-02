package sqls

/*
import (
	"time"
)
*/

type Room struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Owner       string `json:"owner"`
	CreatedDate string `json:"created_date"`
	Vote        int    `json:"vote"`
}

type Rooms []Room

//全てのルームの詳細をROOMデータベースから取得する
func GetAllRoomDetail(DBname string, offset int) (string, []Room) {
	var rooms Rooms
	var room Room
	db := NewDB(DBname)
	defer db.Close()
	state := "SELECT r.room_name, r.title, r.create_user, create_day, (SELECT COUNT(*) FROM VOTING v WHERE r.room_name = v.room_name) as 'vote' FROM ROOM r limit 10 offset ?"
	ftmt, err := db.Prepare(state)
	if err != nil {
		//log.Fatal(err)
		return "サーバー側でエラーが発生しました。", nil
	}
	result, err := ftmt.Query(offset)
	if err != nil {
		//log.Fatal(err)
		return "データの取得に失敗しました。", nil
	}
	defer result.Close()
	for result.Next() {
		err := result.Scan(&room.Name, &room.Title, &room.Owner, &room.CreatedDate, &room.Vote)
		if err != nil {
			return "", rooms
			//log.Fatal(err)
		}
		rooms = append(rooms, room)
	}
	return "", rooms
}

func GetRanking(DBname string) (string, []Room) {
	var rooms Rooms
	var room Room
	db := NewDB(DBname)
	defer db.Close()
	SQL := "SELECT r.room_name, r.title, r.create_user, create_day, (SELECT COUNT(*) FROM VOTING v WHERE r.room_name = v.room_name) as 'vote' FROM ROOM r ORDER BY `vote` DESC LIMIT 10"
	ftmt, err := db.Prepare(SQL)
	if err != nil {
		//log.Fatal(err)
		return "サーバー側でエラーが発生しました。", nil
	}
	result, err := ftmt.Query()
	if err != nil {
		//log.Fatal(err)
		return "データの取得に失敗しました。", nil
	}
	defer result.Close()
	for result.Next() {
		err := result.Scan(&room.Name, &room.Title, &room.Owner, &room.CreatedDate, &room.Vote)
		if err != nil {
			return "", rooms
			//log.Fatal(err)
		}
		rooms = append(rooms, room)
	}
	return "", rooms
}

//ルームを作る
//func CreateRoom()
func CreateRoom(DBname, Name, Title, Owner string) (bool, string) {
	db := NewDB(DBname)
	defer db.Close()

	insert_keyword := "INSERT INTO `ROOM` (room_name, title, create_user, create_day) VALUES (?, ?, ?, NOW())"
	insert, err := db.Prepare(insert_keyword)
	if err != nil {
		return false, "接続エラーが発生しました"
		//log.Fatal(err)
	}
	_, errs := insert.Exec(Name, Title, Owner)
	if errs != nil {
		return false, "ルーム名を変更してください"
		//log.Fatal(err)
	}
	insert.Close()
	return true, ""
}

func GetSearch(DBname, RoomName string) (string, []Room) {
	var rooms Rooms
	var room Room
	db := NewDB(DBname)
	defer db.Close()
	SQL := "SELECT r.room_name, r.title, r.create_user, create_day, (SELECT COUNT(*) FROM VOTING v WHERE r.room_name = v.room_name) as 'vote' FROM ROOM r where room_name like ?"
	ftmt, err := db.Prepare(SQL)
	if err != nil {
		//log.Fatal(err)
		return "サーバー側でエラーが発生しました。", nil
	}
	result, err := ftmt.Query("%" + RoomName + "%")
	if err != nil {
		//log.Fatal(err)
		return "データの取得に失敗しました。", nil
	}
	defer result.Close()
	for result.Next() {
		err := result.Scan(&room.Name, &room.Title, &room.Owner, &room.CreatedDate, &room.Vote)
		if err != nil {
			return "", rooms
			//log.Fatal(err)
		}
		rooms = append(rooms, room)
	}
	return "", rooms
}

//ルームを削除する
//func DeleteRoom()

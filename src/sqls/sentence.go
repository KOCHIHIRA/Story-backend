package sqls

import (
	"log"
)

type Story struct {
	UserName string `json:"user"`
	Sentence string `json:"sentence"`
}

type Storys []Story

//STATEMENTテーブルからルームの入力文字を取得する

func CheckPermission(DBname, RoomName, UserName string) (bool, string) {
	db := NewDB(DBname)
	defer db.Close()
	SQL := "SELECT talk_word FROM TALKED WHERE room_name = ? AND user_name = ?"
	ftmt, err := db.Prepare(SQL)
	if err != nil {
		log.Fatal(err)
		return false, ""
	}
	result, err := ftmt.Query(RoomName, UserName)
	if err != nil {
		//log.Fatal(err)
		return false, ""
	}
	defer result.Close()
	if result.Next() {
		return false, ""
	}
	return true, ""
}

//STATEMENTテーブルにルームに入力した文字を入れる。
func SetSentence(DBname, RoomName, UserName, sentence string) (bool, string) {
	db := NewDB(DBname)
	defer db.Close()

	SQL := "INSERT INTO `TALKED` (user_name, room_name, talk_word, talked_time) VALUES (?, ?, ?, Now())"
	insert, err := db.Prepare(SQL)
	if err != nil {
		return false, "接続エラーが発生しました"
		//log.Fatal(err)
	}
	_, errs := insert.Exec(UserName, RoomName, sentence)
	if errs != nil {
		return false, "ルーム名を変更してください"
		//log.Fatal(err)
	}
	insert.Close()
	return true, ""
}

func GetRoomSentence(DBname, RoomName string) (bool, []Story) {
	var storys Storys
	var story Story
	SQL := "SELECT user_name, talk_word FROM TALKED WHERE room_name = ? order by talked_time"
	db := NewDB(DBname)
	defer db.Close()
	ftmt, err := db.Prepare(SQL)
	if err != nil {
		//log.Fatal(err)
		return false, nil
	}
	result, err := ftmt.Query(RoomName)
	if err != nil {
		log.Fatal(err)
		return false, nil
	}
	defer result.Close()
	for result.Next() {
		err := result.Scan(&story.UserName, &story.Sentence)
		if err != nil {
			//return false, storys
			log.Fatal(err)
		}
		storys = append(storys, story)
	}
	return true, storys
}

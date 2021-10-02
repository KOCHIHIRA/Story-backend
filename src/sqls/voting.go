package sqls

import (
	"log"
)

func Voting(DBname, userName, roomName string) (bool, string) {
	db := NewDB(DBname)
	defer db.Close()

	SQL := "INSERT INTO `VOTING` (user_name, room_name) SELECT * FROM (SELECT ? as `user`, ? as `room`) AS `VOTE` WHERE NOT EXISTS (SELECT * FROM `VOTING` WHERE user_name = ?)"
	insert, err := db.Prepare(SQL)
	if err != nil {
		//return false, "接続エラーが発生しました"
		log.Fatal(err)
	}
	result, errs := insert.Exec(userName, roomName, userName)
	if errs != nil {
		return false, "ルーム名を変更してください"
		//log.Fatal(err)
	}
	rows, err := result.RowsAffected()
	if rows == 0 {
		return false, "既に投票しています。"
	}
	insert.Close()
	return true, ""
}

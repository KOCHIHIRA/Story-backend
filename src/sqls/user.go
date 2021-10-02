package sqls

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Mail     string `json:"mail"`
}

//ユーザーをDBに登録する
func SignUp(DBname, User, Password, Mail string) (bool, string) {
	db := NewDB(DBname)
	defer db.Close()
	SQL := "INSERT INTO `USER` (name, password, mail) VALUES (?, ?, ?)"
	insert, err := db.Prepare(SQL)
	if err != nil {
		return false, "接続エラーが発生しました"
		//log.Fatal(err)
	}

	_, errs := insert.Exec(User, Password, Mail)

	if errs != nil {
		return false, "ユーザー名を変更してください"
		//log.Fatal(err)
	}

	insert.Close()
	return true, "ユーザーの登録が成功しました。"
}

//ユーザーログイン認証する時用
func LogIn(DBname, UserName, Password string) (bool, string) {
	db := NewDB(DBname)
	defer db.Close()
	SQL := "SELECT name FROM USER where name = ? AND password = ?"
	ftmt, err := db.Prepare(SQL)
	if err != nil {
		//log.Fatal(err)
		return false, ""
	}
	result, err := ftmt.Query(UserName, Password)
	if err != nil {
		//log.Fatal(err)
		return false, "接続エラー"
	}
	defer result.Close()
	if result.Next() {
		return true, ""
	}
	return false, "パスワードが違いますよ。"
}

package model

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var MyDB *sql.DB

func DataBaseInit() {
	dataSource := os.Getenv("HEROKU_DATABASE_URL")
	var err error
	MyDB, err = sql.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}
}

func InsertNewUserByTwitter(twitterId int64) error {
	_, err := MyDB.Exec("INSERT INTO users (userName, userImage, homeImage, moodMessage, twitterId, myPoint)"+
		" VALUES(?, ?, ?, ?, ?, ?)", DefaultUser.UserName, DefaultUser.UserImage, DefaultUser.HomeImage,
		DefaultUser.MoodMessage, twitterId, DefaultUser.MyPoint)
	if err != nil {
		return err
	}
	return nil
}

func IsUserExistByTwitter(twitterId int64) bool {
	var count int64
	if err := MyDB.QueryRow("select count(id) from users where twitterId = ?;", twitterId).Scan(&count); err != nil {
		return false
	}
	if count == 1 {
		return true
	} else {
		return false
	}
}

func IsUserExistById(Id string) bool {
	var count int64
	if err := MyDB.QueryRow("select count(id) from users where id = ?;", Id).Scan(&count); err != nil {
		return false
	}
	if count == 1 {
		return true
	} else {
		return false
	}
}

func SelectUserDataById(id string) (User, error) {
	var userData User
	if err := MyDB.QueryRow("select id, userName, userImage, homeImage, moodMessage, "+
		"myPoint from users where id = ?;", id).Scan(&userData.Id, &userData.UserName,
		&userData.UserImage, &userData.HomeImage, &userData.MoodMessage, &userData.MyPoint); err != nil {
		return User{}, err
	}
	return userData, nil
}

func SelectUserPointByTwitter(id int64) (int64, error) {
	var myPoint int64
	if err := MyDB.QueryRow("select myPoint from users where twitterId = ?;", id).Scan(&myPoint); err != nil {
		return myPoint, err
	}
	return myPoint, nil
}

func UpdatePointAddByTwitter(id int64, cp int64) error {
	_, err := MyDB.Exec("update users set myPoint = myPoint + ? where id = ?;", cp, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdatePointSubByTwitter(twitterId int64, cp int64) error {
	_, err := MyDB.Exec("update users set myPoint = myPoint - ? where twitterId = ?;", cp, twitterId)
	if err != nil {
		return err
	}
	return nil
}

func SelectAllUserLIMIT100() []User {
	rows, err := MyDB.Query("select id, userName, userImage, homeImage, moodMessage, twitterId, myPoint from users LIMIT 100")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	users := make([]User, 0)

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.UserName, &user.UserImage, &user.HomeImage, &user.MoodMessage, &user.TwitterId, &user.MyPoint); err != nil {
			log.Println(err)
			continue
		}
		users = append(users, user)
	}
	return users
}

package main

import (
	"database/sql"
	_ "encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type User struct {
	id               int    `json:"id"`
	icon             int    `json:"icon"`
	first_name       string `json:"first_name"`
	last_name        string `json:"last_name"`
	middle_name      string `json:"middle_name"`
	login            string `json:"login"`
	password         string `json:"password"`
	status           int    `json:"status"`
	keyInTime        string `json:"keyintime"`
	score            int    `json:"score"`
	listIdOfTourment string `json:"listIdOfTourment"`
	invites          string `json:"invites"`
}

var infoDb = InfoAboutdb()
var db, errDb = sql.Open("mysql", infoDb)

// func ConnectDB()  {

// }
// func createJson() {

// 	js, _ := json.Marshal(u)
// 	fmt.Println(js)
// }

func InfoAboutdb() string { //получаем инфо о пользователе, пароле, и название db

	file, err := os.Open("pass.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	date := make([]byte, 64) //файл читается как послеедовательность байт, поэтому заводим для него байт массив

	n, err := file.Read(date) //в date мы кладем байты, а n указывает на то сколько байт лежит в переменной date
	if err == io.EOF {
		fmt.Println("EOF")
	}
	//fmt.Println("info: ", string(date[:n]))
	return string(date[:n])
}
func AddDb(firstName string, icon int, lastName string, middleName string, login string, password string, status string) { // Добавляем в базу данных запись о пользователе, используется только при регистрации нового пользователя

	//id - должен быть уникален
	//пароль может быть не уникален

	var (
		id = FindMaxId() + 1 //максимальный id в базе +1
		//id               = 1234
		standartScore    = 1200
		listIdOfTourment = ""
		invites          = ""
		//userStatus    = 0 //пока что по умолчанию участник
	)

	//user := User{id, firstName, lastName, userStatus}
	//db.Query("INSERT INTO UsersTable (id, first_name, last_name, user_status) VALUES (", id, ",", firstName, ",", lastName, ",", userStatus, ");")
	//result, err := db.Exec("INSERT INTO users.users (ID ,Icon ,FirstName ,LastName ,MiddleName ,Login , Password , Status , KeyInTime , Score, ListIdOfTourment JSON,Invites JSON) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)", id, nil, firstName, lastName, middleName, login, password, status, 0, standartScore, nil, nil)
	result, err := db.Exec("INSERT INTO users.users_second (ID ,Icon ,FirstName ,LastName ,MiddleName ,Login , Password , Status , KeyInTime , Score, ListIdOfTourment,Invites) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)", id, 0, firstName, lastName, middleName, login, password, status, 0, standartScore, listIdOfTourment, invites)
	if err != nil {
		panic(err)
	}
	fmt.Println("download record: ", result.LastInsertId)
	//db.Close()
	//_,err = db.Exec("UPDATE User_db.User_Table WHERE id=?")
}
func ShowDb() {
	rows, err := db.Query("select * from Users.users_second")
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()

	users := []User{}
	for rows.Next() {
		u := User{}
		err_scan := rows.Scan(&u.id, &u.icon, &u.first_name, &u.last_name, &u.middle_name, &u.login, &u.password, &u.status, &u.score, &u.keyInTime, &u.listIdOfTourment, &u.invites)
		if err_scan != nil {
			fmt.Println(err_scan)
			continue
		}
		users = append(users, u)
	}
	for i := range users {
		fmt.Println("id = ", users[i].id, "FirstName:", users[i].first_name, "LastName:", users[i].last_name)
	}

}
func Remove(id int) { // удаление пользователя по id
	_, err := db.Exec("DELETE FROM users.users_second WHERE id=?", id)
	if err != nil {
		log.Panic(err)
	}
	// _, err = db.Exec("update User_db.User_Table WHERE id=?", id)
	// if err != nil {
	// 	panic(err)
	// }
}
func FindMaxId() int { // Поиск максимального id в бд, нужна тобы id был всегда уникальным
	maxID := db.QueryRow("SELECT MAX(id) FROM Users_second")
	var u int
	err := maxID.Scan(&u)
	if err != nil {
		panic(err)
	}
	//fmt.Println(u)
	return u
}
func FindUser(id int) User { // пооиск пользователя по id
	result := db.QueryRow("SELECT *FROM users.users_second WHERE id=?", id)
	u := User{}
	err := result.Scan(&u.id, &u.icon, &u.first_name, &u.last_name, &u.middle_name, &u.login, &u.password, &u.status, &u.score, &u.keyInTime, &u.listIdOfTourment, &u.invites)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(u.first_name, u.last_name)
	return u
}
func FindUserFNLN(first_name string, last_name string) int {
	result := db.QueryRow("SELECT *FROM Users.Users_second WHERE first_name=? and last_name=?", first_name, last_name)
	u := User{}
	err := result.Scan(&u.id, &u.first_name, &u.last_name, &u.status)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(u.id)
	return u.id
}

package main

import (
	"database/sql"
	_ "encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	_ "strings"
)

type User struct {
	Id               int    `json:"id"`
	Icon             int    `json:"icon"`
	First_name       string `json:"first_name"`
	Last_name        string `json:"last_name"`
	Middle_name      string `json:"middle_name"`
	Login            string `json:"login"`
	Password         string `json:"password"`
	Status           int    `json:"status"`
	KeyInTime        string `json:"keyintime"`
	Score            int    `json:"score"`
	ListIdOfTourment string `json:"listIdOfTourment"`
	Invites          string `json:"invites"`
}

type Tourments struct {
	ID            int    `json:"id"`
	ListOfMembers string `json:"listOfMembers"`
	Name          string `json:"name"`
	Date          string `json:"date"`
	MinRating     int    `json:"minRating"`
	FullName      string `json:"fullName"`
	Result        string
}
type ListOfMembersList struct {
	UsersId []int `json:"ListOfUsersId"`
}
type ListOfTourments struct {
	Tourments []int `json"tourmentsId"`
}

var infoDb = InfoAboutdb()
var db, errDb = sql.Open("mysql", infoDb)

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
		err_scan := rows.Scan(&u.Id, &u.Icon, &u.First_name, &u.Last_name, &u.Middle_name, &u.Login, &u.Password, &u.Status, &u.Score, &u.KeyInTime, &u.ListIdOfTourment, &u.Invites)
		if err_scan != nil {
			fmt.Println(err_scan)
			continue
		}
		users = append(users, u)
	}
	for i := range users {
		fmt.Println("id = ", users[i].Id, "FirstName:", users[i].First_name, "LastName:", users[i].Last_name)
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
		return 1234
	}
	//fmt.Println(u)
	return u
}
func FindUser(id int) User { // пооиск пользователя по id
	result := db.QueryRow("SELECT *FROM users.users_second WHERE id=?", id)
	u := User{}
	err := result.Scan(&u.Id, &u.Icon, &u.First_name, &u.Last_name, &u.Middle_name, &u.Login, &u.Password, &u.Status, &u.KeyInTime, &u.Score, &u.ListIdOfTourment, &u.Invites)
	if err != nil {
		log.Panic(err, "FindUser")
	}
	//fmt.Println(u.First_name, u.Last_name)
	return u
}

func FindUserFNLN(login string, password string) int {
	result := db.QueryRow("SELECT *FROM Users.Users_second WHERE Login=? and Password=?", login, password)
	u := User{}
	err := result.Scan(&u.Id, &u.Icon, &u.First_name, &u.Last_name, &u.Middle_name, &u.Login, &u.Password, &u.Status, &u.KeyInTime, &u.Score, &u.ListIdOfTourment, &u.Invites)
	if err != nil {
		return 0 //Если возвращаетс 0 то пользователь с таким логином и паролем не найден
	}
	//fmt.Println(u.Id)
	return u.Id // если вернулось число больше 1236 то пользовательский id найден
}

func AddDbTourment(name string, data string, minRating int, fullName string) { // Добавляем в базу данных запись о пользователе, используется только при регистрации нового пользователя

	//id - должен быть уникален
	//пароль может быть не уникален

	var (
		id            = FindMaxIdInTourment() + 1 //максимальный id в базе +1
		listOfMembers = ""
	)

	result, err := db.Exec("INSERT INTO users.TourmentsTable (Id,ListOfMembers,Name ,Date ,MinRating ,FullName, Result) VALUES (?,?,?,?,?,?,?)", id, listOfMembers, name, data, minRating, fullName, "")
	if err != nil {
		panic(err)
	}
	fmt.Println("download record: ", result.LastInsertId)
}
func FindMaxIdInTourment() int { // Поиск максимального id в бд, нужна тобы id был всегда уникальным
	maxID := db.QueryRow("SELECT MAX(ID) FROM TourmentsTable")
	var u int
	err := maxID.Scan(&u)
	if err != nil {
		return 1234
	}
	//fmt.Println(u)
	return u
}
func AddUserInTourment(userId int, tourmentId int) string {
	result := db.QueryRow("SELECT *FROM users.TourmentsTable WHERE ID=?", tourmentId)
	resultUser := db.QueryRow("SELECT *FROM users.Users_second WHERE ID=?", userId)
	tourment := Tourments{}
	err := result.Scan(&tourment.ID, &tourment.ListOfMembers, &tourment.Name, &tourment.Date, &tourment.MinRating, &tourment.FullName, &tourment.Result)
	if err != nil {
		log.Panic(err, "This id is not found...")
		return "This id is not found..."
	}
	u := User{}
	err1 := resultUser.Scan(&u.Id, &u.Icon, &u.First_name, &u.Last_name, &u.Middle_name, &u.Login, &u.Password, &u.Status, &u.KeyInTime, &u.Score, &u.ListIdOfTourment, &u.Invites)
	if err1 != nil {
		log.Panic(err1, "This id is not found...")
		return "This id is not found..."
	}
	listUser := u.ListIdOfTourment
	list := tourment.ListOfMembers
	if listUser == "" {
		listUser = strconv.Itoa(tourmentId)
		_, err = db.Exec("UPDATE users.Users_second set ListIdOfTourment=? WHERE ID=?", listUser, userId)
	} else {
		str := strings.Split(listUser, ",")
		for _, value := range str {
			if value == strconv.Itoa(tourmentId) {
				return "This user is participate !"
			}
		}
		listUser += "," + strconv.Itoa(tourmentId)
		_, err = db.Exec("UPDATE users.Users_second set ListIdOfTourment=? WHERE ID=?", listUser, userId)

	}
	if list == "" {
		list = strconv.Itoa(userId)
		_, err = db.Exec("UPDATE users.TourmentsTable set ListOfMembers=? WHERE ID=?", list, tourmentId)
		return "list is empty"
	}
	list += "," + strconv.Itoa(userId)
	_, err = db.Exec("UPDATE users.TourmentsTable set ListOfMembers=? WHERE ID=?", list, tourmentId)
	return list
}
func getAllMembersTourment(tourmentId int) []string {
	result := db.QueryRow("SELECT *FROM users.TourmentsTable WHERE ID=?", tourmentId)
	tourment := Tourments{}
	err := result.Scan(&tourment.ID, &tourment.ListOfMembers, &tourment.Name, &tourment.Date, &tourment.MinRating, &tourment.FullName, &tourment.Result)
	if err != nil {
		log.Panic(err)
	}
	str := strings.Split(tourment.ListOfMembers, ",")
	return str
}

// func ChangeStr(str string) []string{
// 	str = strings.Split()
// }
func DeleteUserFromTourment(userId int, tourmentId int) string {
	result := db.QueryRow("SELECT *FROM users.TourmentsTable WHERE ID=?", tourmentId)
	resultUser := db.QueryRow("SELECT *FROM users.Users_second WHERE ID=?", userId)
	tourment := Tourments{}
	err := result.Scan(&tourment.ID, &tourment.ListOfMembers, &tourment.Name, &tourment.Date, &tourment.MinRating, &tourment.FullName, &tourment.Result)
	if err != nil {
		log.Panic(err, "This id is not found...")
		return "This id is not found..."
	}
	u := User{}
	err1 := resultUser.Scan(&u.Id, &u.Icon, &u.First_name, &u.Last_name, &u.Middle_name, &u.Login, &u.Password, &u.Status, &u.KeyInTime, &u.Score, &u.ListIdOfTourment, &u.Invites)
	if err1 != nil {
		log.Panic(err1, "This id is not found...")
		return "This id is not found..."
	}
	listMembers := tourment.ListOfMembers
	if listMembers == "" {
		return "list of members is empty!"
	} else {
		str := strings.Split(listMembers, ",")
		for i, value := range str {
			if value == strconv.Itoa(userId) {
				str = append(str[:i], str[i+1:]...)
			}
		}
		var updateListMembers string
		for _, value := range str {
			if updateListMembers != "" {
				updateListMembers += "," + value
			} else {
				updateListMembers += value
			}
		}
		_, err = db.Exec("UPDATE users.TourmentsTable set ListOfMembers=? WHERE ID=?", updateListMembers, tourmentId)
	}
	listTourments := u.ListIdOfTourment
	str1 := strings.Split(listTourments, ",")
	for i, value := range str1 {
		if value == strconv.Itoa(tourmentId) {
			str1 = append(str1[:i], str1[i+1:]...)
		}
	}
	var updateListTourments string
	for _, value := range str1 {
		if updateListTourments != "" {
			updateListTourments += "," + value
		} else {
			updateListTourments += value
		}
	}
	_, err = db.Exec("UPDATE users.users_second set ListIdOfTourment=? WHERE ID=?", updateListTourments, userId)
	return "Success !"
}
func ChangeInfoInUser(id int, firstName string, icon int, lastName string, middleName string, login string, password string, status string) {
	_, _ = db.Exec("UPDATE users.users_second set Icon=?, FirstName=?, LastName=?, MiddleName=?, Login=?, Password=?, Status=? WHERE ID=?", icon, firstName, lastName, middleName, login, password, status, id)
}
func ChangeInfoInTorument(id int, name string, data string, minRating int, fullName string) {
	_, _ = db.Exec("UPDATE users.TourmentsTable set Name=?, Date=?, MinRating=?, FullName=? WHERE ID=?", name, data, minRating, fullName, id)
}
func DeleteTourments(id int) {
	_, err := db.Exec("DELETE FROM users.TourmentsTable WHERE ID=?", id)
	if err != nil {
		log.Panic(err)
	}
}

func Help() map[string]interface{} {
	HelpList := make(map[string]interface{})
	HelpList["addDb"] = "firstName string, lastName string, middleName string, login string, password string, status int"
	HelpList["deleteDb"] = "id int"
	HelpList["getinfo"] = "id int"
	HelpList["auntification"] = "login string, password string"
	HelpList["addDbTourment"] = "name string, date string, fullname string"
	HelpList["addUsersInTourment"] = "idUser int, idTourment int"
	HelpList["getAllMembersTourment"] = "tourmentID int"
	HelpList["deleteUserFromTourment"] = "idUser int, idTourment int"
	HelpList["changeInfoInUser"] = "id int, firstName string, icon int, lastName string, middleName string, login string, password string, status string"
	return HelpList
}

//Func for db:
//SELECT * from TourmentsTable;
// SELECT * from users_second;
//describe users_second;
//describe users_second

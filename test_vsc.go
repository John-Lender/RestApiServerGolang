package main

import (
	_ "database/sql"
	"encoding/json"
	"fmt"
	"log"
	_ "log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	if errDb != nil {
		log.Panicln(errDb)
	}
}

func main() {
	//AddDb("Ivan", "Pokidov")
	//Remove(989)
	//ShowDb()
	//Remove(980)
	//ShowDb()
	//FindUser(984)

	http.HandleFunc("/", mainPage)
	port := ":9090"
	println("Server listen on port", port)
	err := http.ListenAndServe(port, nil) // Мы помещаем в переменную err ощибку, если ошибки не будет то она будет равна
	if err != nil {                       //если вернулся ощибка то мы реагируем на нее определенным образом
		log.Fatal("ListenAndServe", err) // Вызываем прерывание, из за ошибки
	}

}

// func connect_db() {
// 	db, err_ := sql.Open("mysql", "C:/Users/Andrey/Documents/DataBase.sql")
// 	if err_ != nil {
// 		log.Fatal(err_)
// 	}
// 	fmt.Println("status", db.Stats().OpenConnections)
// }

func mainPage(w http.ResponseWriter, r *http.Request) { // Ответ пишем в w(Writer). Стадартный набор для обработки http
	//decoder := json.NewDecoder(r.Body)
	// r.ParseForm()
	// //fmt.Println(r.URL.Path)
	// fmt.Println(1)
	// var result map[string]interface{}
	//user := User{12, "Shannon", "Masq", 0}
	// json.NewDecoder(r.Body).Decode(&result)
	// //println(&result)
	//js, _ := json.Marshal(user)
	//w.Write(js)
	r.ParseForm()
	var method = strings.Split(r.URL.Path, "/")

	switch method[1] {
	case "addDb":
		fmt.Fprintf(w, "Added")
		AddDb(r.Form.Get("firstName"), 0, r.Form.Get("lastName"), r.Form.Get("middleName"), r.Form.Get("login"), r.Form.Get("password"), r.Form.Get("status"))
		ShowDb()
	case "deleteDb":
		index, err := strconv.Atoi(r.Form.Get("id"))
		if err != nil {
			fmt.Fprintf(w, "Parameter id not found")
			log.Panic(err)
		} else {
			Remove(index)
			ShowDb()
		}
	case "getinfo":
		index, err := strconv.Atoi(r.Form.Get("id"))
		if err != nil {
			fmt.Fprintf(w, "error 1: Parameter id not found")
			fmt.Println(err)

		} else {
			u := FindUser(index)
			users := &User{u.id, u.icon, u.first_name, u.last_name, u.middle_name, u.login, u.password, u.status, u.keyInTime, u.score, u.listIdOfTourment, u.invites}
			//users := &User{124, 0, "dueudeu", "rfrfrfr", "rfnrfrfb", "frfrfr", "frfrfrfr", 0, "rfrfdw", 1234, "r vjr vrj", "4vvrvrr"}
			js, _ := json.Marshal(users)
			fmt.Fprintf(w, string(js))
			//fmt.Fprintf(w, string(u.id))
			//fmt.Fprintf(w, string(u.id), "   ", string(u.login))
			//fmt.Fprintf(w, []byte(js))
			//fmt.Println(users)
			//w.Write(js)
		}
	case "findFNLN":
		ans := FindUserFNLN(r.Form.Get("firstName"), r.Form.Get("lastName"))
		fmt.Println(ans)
	}

	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
		fmt.Println("---------------------")
	}

}

func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //анализ аргументов,
	//fmt.Println(r.Form) // ввод информации о форме на стороне сервера
	// var method = strings.Split(r.URL.Path, "/")
	// var meth string
	// for i := range method {
	// 	//println(i, method[i])
	// 	if i == 2 {
	// 		meth = method[2]
	// 	}
	// }
	// if method[1] == "addDb" {
	// 	fmt.Println("true")
	// } else {
	// 	fmt.Println("false")
	// }
	//fmt.Println(meth)
	//fmt.Println(method[0], method[1])
	// fmt.Println("path", strings.Split(r.URL.Path, "/"))
	// fmt.Println("scheme", r.URL.Scheme)
	// fmt.Println(r.Form["url_long"])

	fmt.Println("---------------------")
	// for k, v := range r.Form {
	// 	fmt.Println("key:", k)
	// 	fmt.Println("val:", strings.Join(v, ""))
	// 	fmt.Println("---------------------")
	// }
	//var method = r.URL.Path
	//fmt.Println(strings.Split(method, "/"))
	//fmt.Println("path", strings.Split(r.URL.Path, "/"))
	//user := User{12, r.Form.Get("first_name"), r.Form.Get("last_name"), 0}
	//js, _ := js.Marshal(user)
	//fmt.Fprintf(w, user.first_name) // отправляем данные на клиентскую сторону
	//fmt.Println(user.first_name, user.last_name)
}

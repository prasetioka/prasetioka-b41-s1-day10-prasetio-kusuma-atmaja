package main

import (
	"context"

	"bootcamp-day-10/connection"

	"fmt"

	"html/template"

	"log"

	"net/http"

	"time"

	"strconv"

	"github.com/gorilla/mux"

	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/sessions"
)

type MetaData struct {
	Title     string
	IsLogin   bool
	UserName  string
	FlashData string
}

var Data = MetaData{
	Title: "Personal Web",
}

type Project struct {
	Id           int
	Title        string
	Start_date   time.Time
	Format_start string
	End_date     time.Time
	Format_end   string
	Description  string
	Technologies []string
	Image        string
	Duration     string
	IsLogin      bool
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

var Projects = []Project{}

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnect()

	// route path folder untuk public
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// routing
	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/formProject", formProject).Methods("GET")
	route.HandleFunc("/projectPage", projectPage).Methods("GET")
	route.HandleFunc("/projectDetail/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/addProject", addProject).Methods("POST")
	route.HandleFunc("/deleteProject/{id}", deleteProject).Methods("GET")
	route.HandleFunc("/updateForm{id}", updateForm).Methods("GET")
	route.HandleFunc("/updateProject", updateProject).Methods("POST")

	route.HandleFunc("/formRegister", formRegister).Methods("GET")
	route.HandleFunc("/register", register).Methods("POST")

	route.HandleFunc("/formLogin", formLogin).Methods("GET")
	route.HandleFunc("/login", login).Methods("POST")

	route.HandleFunc("/logout", logout).Methods("GET")

	fmt.Println("Server Running on port 5000")
	http.ListenAndServe("localhost:5000", route)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	// fm := session.Flashes("message")

	// var flashes []string
	// if len(fm) > 0 {
	// 	session.Save(r, w)
	// 	for _, fl := range fm {
	// 		flashes = append(flashes, fl.(string))
	// 	}
	// }

	rows, _ := connection.Conn.Query(context.Background(), "SELECT id, title, start_date, end_date, description, technologies, image FROM tb_projects")

	var result []Project

	for rows.Next() {
		var each = Project{}

		err := rows.Scan(&each.Id, &each.Title, &each.Start_date, &each.End_date, &each.Description, &each.Technologies, &each.Image)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		each.Format_start = each.Start_date.Format("2 January 2006")
		each.Format_end = each.End_date.Format("2 January 2006")

		each = Project{
			Id:           each.Id,
			Title:        each.Title,
			Duration:     DurationCount(each.Start_date, each.End_date),
			Description:  each.Description,
			Technologies: each.Technologies,
			Image:        each.Image,
		}
		result = append(result, each)
	}
	fmt.Println(result)

	respData := map[string]interface{}{
		"Data":     Data,
		"Projects": result,
		// "Data.FlashData": strings.Join(flashes, "message"),
	}
	fmt.Println(Data.IsLogin)

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contact-Type", "text/html; charset=utf-8")

	var Data = MetaData{}
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	var tmpl, err = template.ParseFiles("views/contact.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	data := map[string]interface{}{
		"Data": Data}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

func formProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contact-Type", "text/html; charset=utf-8")

	var Data = MetaData{}
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	var tmpl, err = template.ParseFiles("views/form-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	data := map[string]interface{}{
		"Data": Data}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

func projectPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contact-Type", "text/html; charset=utf-8")

	var Data = MetaData{}
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	var tmpl, err = template.ParseFiles("views/project-page.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	data := map[string]interface{}{
		"Data": Data}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contact-Type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var tmpl, err = template.ParseFiles("views/project-page.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var Data = MetaData{}

	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	var ProjectDetail = Project{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, title, start_date, end_date, description FROM tb_projects WHERE id=$1", id).Scan(&ProjectDetail.Id, &ProjectDetail.Title, &ProjectDetail.Start_date, &ProjectDetail.End_date, &ProjectDetail.Description)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	ProjectDetail.Format_start = ProjectDetail.Start_date.Format("2 January 2006")
	ProjectDetail.Format_end = ProjectDetail.End_date.Format("2 January 2006")

	ProjectDetail = Project{
		Id:           ProjectDetail.Id,
		Title:        ProjectDetail.Title,
		Format_start: ProjectDetail.Format_start,
		Format_end:   ProjectDetail.Format_end,
		Duration:     DurationCount(ProjectDetail.Start_date, ProjectDetail.End_date),
		Description:  ProjectDetail.Description,
		Technologies: ProjectDetail.Technologies,
		Image:        ProjectDetail.Image,
	}

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
	}

	data := map[string]interface{}{
		"Project": ProjectDetail,
		"Data":    Data,
	}

	// fmt.Println(data)

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var title = r.PostForm.Get("input-title")
	var startDate = r.PostForm.Get("start-date")
	var endDate = r.PostForm.Get("end-date")
	var description = r.PostForm.Get("project-description")
	var checkbox1 = r.PostForm.Get("checkbox1")
	var checkbox2 = r.PostForm.Get("checkbox2")
	var checkbox3 = r.PostForm.Get("checkbox3")
	var checkbox4 = r.PostForm.Get("checkbox4")
	var technologies = []string{
		checkbox1, checkbox2, checkbox3, checkbox4,
	}
	// var image = r.PostForm.Get("upload-image")

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(title, start_date, end_date, description, technologies, image) VALUES ($1, $2, $3, $4, $5, 'https://ewr1.vultrobjects.com/lmsbzzbx/blog/ufvyaQoTY5c.jpg')", title, startDate, endDate, description, technologies)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func updateForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contact-Type", "text/html; charset=utf-8")

	var Data = MetaData{}
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	var tmpl, err = template.ParseFiles("views/form-update.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	data := map[string]interface{}{
		"Data": Data}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

func updateProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var title = r.PostForm.Get("input-title")
	var startDate = r.PostForm.Get("start-date")
	var endDate = r.PostForm.Get("end-date")
	var description = r.PostForm.Get("project-description")

	_, err = connection.Conn.Exec(context.Background(), "UPDATE tb_projects SET title=$1, start_date=$2, end_date=$3, description=$4", title, startDate, endDate, description)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func formRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contact-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/form-register.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var name = r.PostForm.Get("inputName")
	var email = r.PostForm.Get("inputEmail")
	var password = r.PostForm.Get("inputPassword")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	// fmt.Println(passwordHash)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/formLogin", http.StatusMovedPermanently)
}

func formLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contact-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/form-login.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := r.PostForm.Get("inputEmail")
	password := r.PostForm.Get("inputPassword")

	user := User{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_user WHERE email=$1", email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	session.Values["IsLogin"] = true
	session.Values["Name"] = user.Name
	session.Options.MaxAge = 10800 // lama login 3 jam

	// session.AddFlash("Successfully Login!", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func logout(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DurationCount(Start_date time.Time, End_date time.Time) string {
	days := End_date.Sub(Start_date).Hours() / 24 // Selisih tanggal, output berubah ke jam, jam berubah jadi hari, masih dalam float64
	var duration string

	if days >= 30 {
		if (days / 30) == 1 {
			duration = "duration 1 month"
		} else {
			duration = "duration " + strconv.Itoa(int(days/30)) + " months" //float64 diubah ke int, diubah lagi ke string
		}
	} else {
		if days <= 1 {
			duration = "1 day"
		} else {
			duration = "duration " + strconv.Itoa(int(days)) + " days"
		}
	}
	return duration
}

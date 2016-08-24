package main

import (
"encoding/json"
 "html/template"
   "fmt"
   "net/http"
   "log"
    "github.com/garyburd/redigo/redis"
    
)

func handler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("home.html")
        t.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
   
    // Redis connection
    conn, err := redis.Dial("tcp", ":9090")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // Set UserName
    if _, err = conn.Do("SET", "username", "Prashant"); err != nil {
        log.Fatal(err)
    }

    //Set Password
    if _, err = conn.Do("SET", "password", "Prashant@91"); err != nil {
        log.Fatal(err)
    }

     strs, err := redis.Strings(conn.Do("MGET", "username", "password"))
    if err != nil {
        log.Fatal(err)
    }

    // prints [a b ]
    fmt.Println(strs[0])
    fmt.Println(strs[1])
    

     fmt.Println("method:", r.Method) //get request method
    if r.Method == "GET" {
        t, _ := template.ParseFiles("login.html")
        t.Execute(w, nil)
    } else {
        r.ParseForm()
        name := r.PostFormValue("username")
        pass := r.PostFormValue("password")

        fmt.Println(name)
        fmt.Println(pass)
      
       if(strs[0] == name && strs[1] == pass){
       	s, _ := template.ParseFiles("home.html")
        s.Execute(w, nil)
       }else{
        s1, _ := template.ParseFiles("login.html")
        s1.Execute(w, nil)
       }
    }
}

func register(w http.ResponseWriter, r *http.Request) {

	fmt.Println("method:", r.Method) //get request method
    if r.Method == "GET" {
        t, _ := template.ParseFiles("register.html")
        t.Execute(w, nil)
    } else {
       /* r.ParseForm()
        firstname := r.PostFormValue("firstname")
        lastname := r.PostFormValue("lastname")
        name := r.PostFormValue("username")
        pass := r.PostFormValue("password")*/
     
    decoder := json.NewDecoder(r.Body)
    var t test_struct   
    err := decoder.Decode(&t)
    if err != nil {
        panic(err)
    }

       

        // Redis connection
    conn, err := redis.Dial("tcp", ":9090")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

     // Set FirstName
    if _, err = conn.Do("SET", "firstname", t.FirstName); err != nil {
        log.Fatal(err)
    }


     // Set LastName
    if _, err = conn.Do("SET", "lastname", t.LastName); err != nil {
        log.Fatal(err)
    }


    // Set UserName
    if _, err = conn.Do("SET", "username", t.UserName); err != nil {
        log.Fatal(err)
    }

    //Set Password
    if _, err = conn.Do("SET", "password", t.Password); err != nil {
        log.Fatal(err)
    }

    strs, err := redis.Strings(conn.Do("MGET", "firstname","lastname", "username", "password"))
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(t.FirstName)

     js, err := json.Marshal(strs)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
    }


	
}





type test_struct struct {
    FirstName string   `json:"firstname"`
    LastName string     `json:"lastname"`
    UserName string     `json:"username"`
    Password string     `json:"password"`
}

func jsonTest(w http.ResponseWriter, r *http.Request) {

    decoder := json.NewDecoder(r.Body)
    var t test_struct   
    err := decoder.Decode(&t)
    if err != nil {
        panic(err)
    }
    log.Println(t.FirstName)
    log.Println(t.LastName)
    log.Println(t.UserName)
    log.Println(t.Password)
}

type Profile struct {
  Name    string
  Hobbies []string
}

func jsonResponse(w http.ResponseWriter, r *http.Request) {
    
     profile := Profile{"Alex", []string{"snowboarding", "programming"}}

  js, err := json.Marshal(profile)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  w.Write(js)
}



func main() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/login",login)
    http.HandleFunc("/register",register)
    http.HandleFunc("/jsonTest",jsonTest)
    http.HandleFunc("/jsonResponse",jsonResponse)
    http.ListenAndServe(":8080", nil)
}

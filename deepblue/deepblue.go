package deepblue

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	Id    		 string `json:"id"`
	Connection   bool   `json:"connection"`
	Point 		 int    `json:"point"`
}

type Node struct {
	Id     int
	Port   string
	Router *mux.Router
}

type Cluster []Node

var userCSV []int
var users	[]User
var node, port int

func randomIdGenerator(n int) []int {
	var rands []int
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	for i := 0; i < n; i++ {
		for {
		up:
			num := random.Intn(999999)
			if num > 99999 {
				for j := 0; j < len(rands); j++ {
					if rands[j] == num {
						goto up
					}
				}
				fmt.Println(num)
				rands = append(rands, num)
				break
			}
		}
	}
	return rands
}

func Init(n, scale int) {
	var lc Cluster
	master := GetMasterNode(n, 8080)
	lc = append(lc, *master)

	for i := 0; i < scale - 1; i++ {
		slave := GetSimNode()
		lc = append(lc, *slave)
	}

	for i := 0; i < scale; i++ {
		RaiseNode(&lc[i])
	}

	fmt.Println(lc)
	select {}
}

func RaiseNode(p *Node) {
	fmt.Println("Node ",p.Id," is running on PORT ",p.Port)
	go func() {
		log.Println(http.ListenAndServe(p.Port, p.Router))
	}()
}

func GetMasterNode(n, PORT int) *Node {
	port = PORT
	node = 0
	userCSV = randomIdGenerator(n)
	var root Node
	root.Id = node
	root.Port = "0.0.0.0:"+strconv.Itoa(port)
	root.Router = mux.NewRouter().StrictSlash(true)
	root.Router.HandleFunc("/", Index)
	root.Router.HandleFunc("/users", ConnectedUsers)
	root.Router.HandleFunc("/user/{userId}", GetConnected)
	root.Router.HandleFunc("/user/{userId}/{newPoint}", ChangePoint)
	return &root
}

func GetSimNode() *Node {
	node++
	port++
	var p Node
	p.Id = node
	p.Port = "0.0.0.0:"+strconv.Itoa(port)
	p.Router = mux.NewRouter().StrictSlash(true)
	p.Router.HandleFunc("/", Index)
	p.Router.HandleFunc("/users", ConnectedUsers)
	p.Router.HandleFunc("/user/{userId}", GetConnected)
	p.Router.HandleFunc("/user/{userId}/{newPoint}", ChangePoint)
	return &p
}

func ChangePoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	newPoint := vars["newPoint"]
	for i := 0; i < len(users); i++ {
		if userId == users[i].Id {
			var err error
			users[i].Point, err = strconv.Atoi(newPoint)
			if err != nil {
				fmt.Println("[Control Panel] Error detected on line 53!")
				_, _ = fmt.Fprint(w, "Not Updated")
				return
			}
			_, _ = fmt.Fprint(w, "Updated")
			return
		}
	}
	_, _ = fmt.Fprint(w, "Unauthorized")
}

func GetConnected(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	for i := 0; i < len(userCSV); i++ {
		str := strconv.Itoa(userCSV[i])
		if userId == str {
			conUser := User{
				Id:    userId,
				Connection:   true,
				Point: 100,
			}
			checkExist := false
			for j := 0; j < len(users); j++ {
				if users[j] == conUser {
					users[j].Point = 100
					users[j].Connection = true
					checkExist = true
				}
			}
			if !checkExist {
				users = append(users, conUser)
			}
			_, _ = fmt.Fprint(w, "Authorized")
			return
		}
	}
	_, _ = fmt.Fprint(w, "Unauthorized")
}

func ConnectedUsers(w http.ResponseWriter, _ *http.Request) {
	jm, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jm)
}

func Index(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "<h1><center>" +
		"<b>Control Panel Index</b></center></h1>")
}

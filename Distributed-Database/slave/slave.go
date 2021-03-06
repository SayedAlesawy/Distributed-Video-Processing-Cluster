package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // here
	"github.com/pebbe/zmq4"
)

//======================= Common Functions ==================

func initSubscriber(addr string) *zmq4.Socket {
	subscriber, _ := zmq4.NewSocket(zmq4.SUB)
	subscriber.SetLinger(0)

	subscriber.Connect(addr)
	subscriber.SetSubscribe("")
	return subscriber
}

func initPublisher(addr string) *zmq4.Socket {
	publisher, err := zmq4.NewSocket(zmq4.PUB)
	if err != nil {
		fmt.Print(err)
		return nil
	}
	publisher.SetLinger(0)
	publisher.Bind(addr)
	return publisher
}

func commandDataDeseralizer(s string) (string, string, string) {
	fields := strings.Split(s, ";")
	if len(fields) < 3 {
		if len(fields) < 2 {
			return fields[0], "", ""
		}

		return fields[0], fields[1], ""
	}
	return fields[0], fields[1], fields[2]
}
func registerUser(name string, email string, password string, db *sql.DB) bool {
	sqlStatement := `INSERT INTO clients (name, email, password) VALUES ($1,$2,$3);`
	fmt.Println("[RegisterUser] Saving user data ..")
	_, err := db.Exec(sqlStatement, name, email, password)
	if err != nil {
		log.Println(err)
		return false
	} else {
		fmt.Println("[RegisterUser] Success")
	}

	return true
}
func loginUser(email string, password string, db *sql.DB) int {
	sqlStatement := `SELECT * FROM clients WHERE email=$1 and password=$2;`

	var clientID int
	var clientName, clientEmail, clientPassword string

	row := db.QueryRow(sqlStatement, email, password)
	switch err := row.Scan(&clientID, &clientName, &clientEmail, &clientPassword); err {
	case sql.ErrNoRows:
		return -1
	case nil:
		return clientID
	default:
		fmt.Println(err)
		return -1

	}
}

func connectDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("[DB]Error loading .env file")
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("USER_NAME")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("[DB] Successfully connected!")
	return db
}

//======================= Common Functions ==================

//ReadQueryListner :
func ReadQueryListner(status *string, db *sql.DB, id int, clientIP string) {

	subscriber := initSubscriber(clientIP + "600" + strconv.Itoa(id))
	idPub := initPublisher(clientIP + "8093")
	defer subscriber.Close()

	for {
		s, err := subscriber.Recv(0)
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println("[ReadQueryListner] recieved", s)
		email, password, _ := commandDataDeseralizer(strings.Split(s, ":")[1])
		fmt.Println("[ReadQueryListner] " + email)
		fmt.Println("[ReadQueryListner] " + password)

		id := loginUser(email, password, db)
		if id > 0 {
			idPub.Send(strconv.Itoa(id), 0)
			fmt.Println("[ReadQueryListner] access granted ")
		} else {
			idPub.Send("-15", 0)
			fmt.Println("[ReadQueryListner] access denied")
		}

	}
}

//TrackerUpdateListner :
func TrackerUpdateListner(status *string, db *sql.DB, id int, trackerIP string) {
	subscriber := initSubscriber(trackerIP + "500" + strconv.Itoa(id))
	defer subscriber.Close()

	for {
		s, err := subscriber.Recv(0)
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println("[TrackerUpdateListner] rec", s)
		name, email, password := commandDataDeseralizer(s)
		registerUser(name, email, password, db)

	}
}

//HeartBeatPublisher :
func HeartBeatPublisher(id int, trackerIP string) {
	publisher := initPublisher(trackerIP + "300" + strconv.Itoa(id))

	defer publisher.Close()

	publisher.Bind(trackerIP + "300" + strconv.Itoa(id))

	for range time.Tick(time.Second * 2) {
		publisher.Send("Heartbeat", 0)
		log.Println("send", "Heartbeat:")
	}
}

func main() {
	clientIP := "tcp://127.0.0.1:"
	trackerIP := "tcp://127.0.0.1:"

	db := connectDB()
	defer db.Close()
	status := "Avaliable"
	id := 2
	go HeartBeatPublisher(id+1, trackerIP)
	go TrackerUpdateListner(&status, db, id+1, trackerIP)
	go ReadQueryListner(&status, db, id+1, clientIP)
	for {

	}
}

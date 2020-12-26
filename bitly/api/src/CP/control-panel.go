package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Redirect struct {
	Id 		int
	Slug 	string 	`db:"slug" form:"slug"`
	Url  	string	`db:"url" form:"url"`
	Hits    int     `db:"hits" form:"hits"`
}

//var db, err = sql.Open("mysql", "root:Kapil@123@tcp(localhost:3306)/bitly")
var db, err = sql.Open("mysql", "cmpe281:kapil123@tcp(10.0.1.43:3306)/bitly")
//var conn, errMQ = amqp.Dial("amqp://guest:guest@localhost:5672/")
var conn, errMQ = amqp.Dial("amqp://admin:rabbitbitly@10.0.1.123:5672/")
var ch, errChannel = conn.Channel()
var q, errQueue = ch.QueueDeclare(
				"hello", // name
				false,   // durable
				false,   // delete when unused
				false,   // exclusive
				false,   // no-wait
				nil,     // arguments
				)

var q1, errQueue1Declare = ch.QueueDeclare(
				"ForStats", // name
				false,   // durable
				false,   // delete when unused
				false,   // exclusive
				false,   // no-wait
				nil,     // arguments
			)

var msgs, errConsume = ch.Consume(
				q1.Name, // queue
				"",     // consumer
				true,   // auto-ack
				false,  // exclusive
				false,  // no-local
				false,  // no-wait
				nil,    // args
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func createTable(){

	stmt, err := db.Prepare("CREATE TABLE redirect (id int NOT NULL AUTO_INCREMENT, Slug varchar(40) UNIQUE, Url varchar(400), hits int, last_accessed DATETIME, PRIMARY KEY (id));")
	//stmt, err := db.Prepare("CREATE TABLE redirect (Slug varchar(40) UNIQUE, Url varchar(400), PRIMARY KEY (Slug));")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Redirect Table successfully migrated....")
	}
}

func generateSlug() string {
	var redirect Redirect
	var chars = []rune("0123456789abcdefghijklmnopqrstuvwxyz")
	s := make([]rune, 6)
	for i := range s {
		rand.Seed(time.Now().UnixNano())
		s[i] = chars[rand.Intn(len(chars))]
	}
	row := db.QueryRow("select slug from redirect  where Url=?", string(s))
	err = row.Scan(&redirect.Slug)
	return string(s)
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H {
		"message": fmt.Sprintf("200 OK"),
	})
}

type longURL struct {
	UrlLong string `json:"url"`
}

func getLongUrl(body []byte) (*longURL, error) {
	var s = new(longURL)
	err := json.Unmarshal(body, &s)
	if(err != nil){
		fmt.Println("whoops:", err)
	}
	return s, err
}

func add(c *gin.Context) {

	//UrlQuery := c.Query("url")
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	s, err := getLongUrl([]byte(jsonData))

	//fmt.Printf("jsonData: %s", s.UrlLong)
	Url := s.UrlLong
	//Url := strings.Replace(UrlQuery, "%2F", "/", -1)
	var redirect Redirect
	row := db.QueryRow("select slug from redirect  where Url=?", Url)
	err = row.Scan(&redirect.Slug)
	if err == nil {
		fmt.Print("Already exists")
		c.JSON(http.StatusOK, gin.H {
			"message": fmt.Sprintf("URL already exists"),
			"short_URL": fmt.Sprintf("http://cmpe.sjsu/%s",redirect.Slug),
		})
	} else {
		Slug := generateSlug()
		currentTime := time.Now()
		currentTime.Format("2006-01-02 15:04:05")
		stmt, err := db.Prepare("insert into redirect (`slug`, `url`, `hits`, `last_accessed`) values(?,?,0,?);")
		if err != nil {
			fmt.Print(err.Error())
		}

		_, err = stmt.Exec(Slug, Url, currentTime)
		if err != nil {
			fmt.Print(err.Error())
		}

		//Short_url := "http://cmpe.sjsu/"+Slug

		defer stmt.Close()

		//---------------- RABBIT MQ  --------------------
		//conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		failOnError(errMQ, "Failed to connect to RabbitMQ")
		//defer conn.Close()
		//ch, err := conn.Channel()
		failOnError(errChannel, "Failed to open a channel")
		////defer ch.Close()
		//
		//q, err := ch.QueueDeclare(
		//	"hello", // name
		//	false,   // durable
		//	false,   // delete when unused
		//	false,   // exclusive
		//	false,   // no-wait
		//	nil,     // arguments
		//)
		failOnError(errQueue, "Failed to declare a queue")

		body := Url+"^"+"http://cmpe.sjsu/"+Slug
		//body := Url, Short_url
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing {
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		failOnError(err, "Failed to publish a message")


		// ------------------------------------------------
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("201 Created"),
			"short_URL": fmt.Sprintf("http://cmpe.sjsu/%s",Slug),
			"long_URL": fmt.Sprintf("%s", s.UrlLong),
		})
	}

}

func updateStats(slug string) {
	currentTime := time.Now()
	currentTime.Format("2006-01-02 15:04:05")
	stmt, err := db.Prepare("update redirect set hits = hits + 1, last_accessed=? where slug = ?;" )
	if err != nil {
		fmt.Print(err.Error())
	}

	_, err = stmt.Exec(currentTime, slug)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer stmt.Close()
}

func main() {
	createTable();
	if err != nil {
		fmt.Print(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}
	
	//headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	//originsOk := handlers.AllowedOrigins([]string{"*"})
	//methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	router := gin.Default()
	router.Use(cors.Default())
	router.POST("/create", add)
	router.GET("/ping", ping)

	go func() {
		router.Run(":8000")
		//log.Fatal(http.ListenAndServe(":8000", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
	}()

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
		updateStats(string(d.Body))

	}
}

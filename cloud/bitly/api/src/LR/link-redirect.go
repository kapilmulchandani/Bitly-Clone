package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"strconv"
)

//var db, err = sql.Open("mysql", "root:Kapil@123@tcp(localhost:3306)/bitly")
var db, err = sql.Open("mysql", "cmpe281:kapil123@tcp(10.0.1.43:3306)/bitly")
//var conn, errMQ = amqp.Dial("amqp://guest:guest@localhost:5672/")
var conn, errMQ = amqp.Dial("amqp://admin:rabbitbitly@10.0.1.123:5672/")
var myClient = &http.Client{Timeout: 10 * time.Second}
var urlCache = make(map[string]Documents)
var ch, errChannel = conn.Channel()
var q, errQueueDeclare = ch.QueueDeclare(
	"hello", // name
	false,   // durable
	false,   // delete when unused
	false,   // exclusive
	false,   // no-wait
	nil,     // arguments
)
var msgs, errConsume = ch.Consume(
	q.Name, // queue
	"",     // consumer
	true,   // auto-ack
	false,  // exclusive
	false,  // no-local
	false,  // no-wait
	nil,    // args
)

var q1, errQueue1Declare = ch.QueueDeclare(
	"ForStats", // name
	false,   // durable
	false,   // delete when unused
	false,   // exclusive
	false,   // no-wait
	nil,     // arguments
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type Redirect struct {
	Id 		int
	Slug 	string 	`db:"slug" form:"slug"`
	Url  	string	`db:"url" form:"url"`
	Hits    int     `db:"hits" form:"hits"`
}

type LongURL struct {
	Url string
	Hits string
	Last_accessed string
}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	fmt.Println("getJSON : ", r.Body)
	return json.NewDecoder(r.Body).Decode(target)
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H {
		"message": fmt.Sprintf("200 OK"),
	})
}

func parseURL(urlToBeParsed string) string{
	u, err := url.QueryUnescape(urlToBeParsed)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Changed URL", u)
	return u
}

type  DocumentsArray []struct {
		Key string `json:"key"`
		Record string `json:"record"`
		Json string `json:"json"`
		Vector []string `json:"vector"`
		Message string `json:"message"`
	}


type Documents struct {
	Key string `json:"key"`
	Record string `json:"record"`
	Json string `json:"json"`
	Vector []string `json:"vector"`
	Message string `json:"message"`
}


func getAllDocuments(c *gin.Context){
	//var documents []LongURL
	//url := "http://localhost:9001/api"
	url := "http://nosql-db-cluster-nlb-87d30cfab1597f51.elb.us-east-1.amazonaws.com/api"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("----------------  In GET ALL DOCUMENTS   -------------------------")
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	all_documents := DocumentsArray{}

	body, _ := ioutil.ReadAll(resp.Body)
	errors := json.Unmarshal(body, &all_documents)
	if errors != nil {
		panic(err)
	}
	//fmt.Println("All Documents: ",all_documents)
	c.JSON(200, all_documents)
	//fmt.Println("response Body:", string(body))
}

func getBySlug(c *gin.Context){

	//short_url := c.Query("short_url")
	short_url := parseURL(c.Query("short_url"))
	fmt.Println("short_url in getBySlug:", short_url)
	t := strings.Split(short_url, ".sjsu/")[1]
	slug := strings.TrimSpace(t)
	fmt.Println("slug in getBySlug:", slug)
	long_url:= ""
	hitsInt:=0
	last_accessed := ""


	// ------------------ Check in Cache ---------------------

	//if urlCache[short_url] != {} {
		if longurl2, exist := urlCache[short_url]; exist {
		fmt.Println("FOUND in URL-CACHE")
		//longurl2 := urlCache[short_url]
		fmt.Println("long_url : ", longurl2)
		longurl3 := LongURL{}
		jsonLong := []byte(longurl2.Json)
		if err := json.Unmarshal(jsonLong, &longurl3); err != nil {
			panic(err)
		}
		fmt.Println("HITS IN CACHE: ", longurl3.Hits)
		hitsInt, err = strconv.Atoi(longurl3.Hits)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		long_url = longurl3.Url
		hitsInt = hitsInt+1
		//----------------------------------------------------------

	} else {

		fmt.Println("NOT FOUND in URL-CACHE")
		// ------------------ API : /api --------------------------

		longurl1 := LongURL{}
		//getJson("http://localhost:9001/api/"+slug, &longurl1)
		//getJson("http://10.0.1.69:9090/api/"+slug, &longurl1)
		getJson("http://nosql-db-cluster-nlb-87d30cfab1597f51.elb.us-east-1.amazonaws.com/api/"+slug, &longurl1)
		fmt.Println("longurl1 : ", longurl1)
		fmt.Println("url : ",longurl1.Url)
		fmt.Println("hits : ",longurl1.Hits)
		fmt.Println("last_accessed : ",longurl1.Last_accessed)
		long_url = longurl1.Url
		hitsInt, err = strconv.Atoi(longurl1.Hits)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		hitsInt = hitsInt+1
	}

	//-------------------------  API UPDATE  ------------------------
	fmt.Println("------------IN UPDATE-------------")
	hits := strconv.Itoa(hitsInt)
	//url := "http://localhost:9001/api/" + slug
	//url := "http://10.0.1.69:9090/api/" + slug
	url := "http://nosql-db-cluster-nlb-87d30cfab1597f51.elb.us-east-1.amazonaws.com/api/" + slug
	fmt.Println("URL:>", url)
	currentTime := time.Now()
	currentTime.Format("2006-01-02 15:04:05")
	last_accessed = currentTime.String()
	var jsonStr = []byte(`{"url":"` + long_url + `", "hits":"`+hits+`" , "last_accessed":"`+ last_accessed +`" }`)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	//---------------------------------------------------------------



		err = ch.Publish(
			"",     // exchange
			q1.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing {
				ContentType: "text/plain",
				Body:        []byte(slug),
			})
		failOnError(err, "Failed to publish a message")



		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("301 Found"),
			"url": fmt.Sprintf("%s",long_url),
			"hits": fmt.Sprintf("%s", hits),
			"last_accessed": fmt.Sprintf("%s", last_accessed),
		})



}

func main() {

	// ------------- RABBITMQ ------------------

	failOnError(errMQ, "Failed to connect to RabbitMQ")
	defer conn.Close()

	failOnError(errChannel, "Failed to open a channel")
	//defer ch.Close()


	failOnError(errQueueDeclare, "Failed to declare a queue")


	failOnError(errConsume, "Failed to register a consumer")

	//forever := make(chan bool)

	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/getUrl", getBySlug)
	router.GET("/ping", ping)
	router.GET("/getAllDocuments", getAllDocuments)
	//router.POST("/create", add)

	go func() {
		// spawn an HTTP server in `other` goroutine
		router.Run(":8009")
	}()
	//router.Run(":8001")

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
		message := string(d.Body)
		if strings.Contains(message, "^") {

		}
		s := strings.Split(message, "^")
		long_url, short_url := s[0], s[1]

		slug := strings.Split(short_url, ".sjsu/")[1]
		fmt.Println("SLUG: " + slug)

		enc := json.NewEncoder(os.Stdout)
		d := map[string]string{short_url: long_url}
		fmt.Println(enc.Encode(d))


		// Store in NoSQL DB

		// ------------------ API : /api --------------------------

		//url := "http://localhost:9001/api/" + slug
		//url := "http://10.0.1.69:9090/api/" + slug
		url := "http://nosql-db-cluster-nlb-87d30cfab1597f51.elb.us-east-1.amazonaws.com/api/" + slug
		fmt.Println("URL:>", url)
		currentTime := time.Now()
		currentTime.Format("2006-01-02 15:04:05")
		var jsonStr = []byte(`{"url":"` + long_url + `", "hits": "0", "last_accessed":"`+ currentTime.String() +`" }`)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()



		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		longurl_create := Documents{}
		errors := json.Unmarshal(body, &longurl_create)
		if errors != nil {
			panic(err)
		}

		//urlCache[short_url] = longurl_create
		fmt.Println("LONG URL CREATE :", longurl_create)
		fmt.Println("response Body:", string(body))

		// --------------------------------------------------------
	}

	//<-forever

	// -----------------------------------------
}

package main

import (
	"encoding/json"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	Glog                                                  *log.Logger
	idbUser, idbAddress, idbPass, idbDb, idbQueryInterval string
)

func logInit(fileName string, logtype string) {
	var writer io.Writer
	if fileName != "" && logtype == "both" {
		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("Failed to open log file", fileName, ":", err)
		}
		writer = io.MultiWriter(file, os.Stdout)
	} else if logtype == "file" {
		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("Failed to open log file", fileName, ":", err)
		}
		writer = io.Writer(file)
	} else {
		writer = io.Writer(os.Stdout)
	}
	Glog = log.New(writer, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// put an event into the db
func PutEvent(event *EfEvent) {
	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     idbAddress,
		Username: idbUser,
		Password: idbPass,
	})

	if err != nil {
		log.Fatal("Error creating InfluxDB Client: ", err.Error())
	}
	defer c.Close()

	// Create a new point batch
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  idbDb,
		Precision: "s",
	})

	/*
		// Create a point and add to batch
		tags := map[string]string{"cpu": "cpu-total"}
		fields := map[string]interface{}{
			"idle":   10.1,
			"system": 53.3,
			"user":   46.6,
		}*/

	pt, err := client.NewPoint(event.TagKey, event.Tags, event.Fields, time.Now())
	if err != nil {
		Glog.Println("Error: ", err.Error())
	}
	bp.AddPoint(pt)

	// Write the batch
	c.Write(bp)
}

/*
// retrieve a list of events from the db
func GetEvents(startTime time.Time, endTime time.Time) {

}
*/

//func PutEvent(key string, tags map[string]string, fields map[string]interface{}) {
type EfEvent struct {
	TagKey string
	Tags   map[string]string
	Fields map[string]interface{}
}

func SetEventHandler(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var efe EfEvent
	err := decoder.Decode(&efe)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	PutEvent(&efe)
}

func init() {

	// then contain the next set
	viper.SetConfigName("efconfig")
	viper.AddConfigPath("$HOME/.ef")
	viper.AddConfigPath("/etc/ef/")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("No configuration file loaded - aborting", err.Error())
	}

	viper.SetDefault("idbUser", "efuser")
	viper.SetDefault("idbPass", "pass")
	viper.SetDefault("idbDb", "efdata")
	idbUser = viper.GetString("idbUser")
	idbPass = viper.GetString("idbPass")
	idbDb = viper.GetString("idbDb")

	viper.SetDefault("logtype", "file")
	logFileName := viper.GetStringSlice("logfile")
	logInit(logFileName[0], "both")
}

func main() {
	http.HandleFunc("/event/add", SetEventHandler)
	// start the web server
	if err := http.ListenAndServe(": 8080", nil); err != nil {
		log.Fatal(" ListenAndServe:", err)
	}
}

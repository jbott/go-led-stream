package main

import (
	"fmt"
	"runtime"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
)

var mqtt_chan = make(chan []byte)

func main() {
	ConfigRuntime()
	go StartMQTT()
	StartGin()
}

func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}

var client MQTT.Client

func StartMQTT() {
	fmt.Print("Connecting to MQTT...\n")

	opts := MQTT.NewClientOptions().AddBroker("tcp://192.168.1.21:1883")
	opts.SetClientID("go-simple")

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		msg := <-mqtt_chan
		client.Publish("device/led0/cmd", 0, true, msg)
	}
}

func StartGin() {
	//gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	router.LoadHTMLGlob("resources/*.templ.html")
	router.Static("/static", "resources/static")

	router.GET("/", index)
	router.POST("/cmd", update)

	router.Run(":3000")
}

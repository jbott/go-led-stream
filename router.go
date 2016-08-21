package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/jbott/go-led-stream/led_stream"
)

func index(c *gin.Context) {
	c.HTML(200, "index.templ.html", gin.H{
		"hi": "index",
	})
}

func update(c *gin.Context) {
	var json struct {
		Entries []led_stream.Entry `json:"entries"`
	}

	err := c.BindJSON(&json)

	if err != nil {
		c.Error(err)
		c.Err()
		return
	}

	buf := new(bytes.Buffer)
	for _, e := range json.Entries {
		buf.Write(led_stream.EntryToBytes(e))
	}

	PublishMQTT(buf.Bytes())

	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func PublishMQTT(buf []byte) {
	packet := led_stream.WrapHeaderCRC(led_stream.PACKET_SET_CMDS, buf)
	mqtt_chan <- packet
}

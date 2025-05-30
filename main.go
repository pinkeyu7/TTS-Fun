package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

type Message struct {
	Text string `json:"text" binding:"required"`
}

type Input struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type Voice struct {
	Model  string `json:"model"`
	Speed  int    `json:"speed"`
	Pitch  int    `json:"pitch"`
	Energy int    `json:"energy"`
}

type AudioConfig struct {
	Encoding   string `json:"encoding"`
	SampleRate string `json:"sampleRate"`
}

type Payload struct {
	Input       Input       `json:"input"`
	Voice       Voice       `json:"voice"`
	AudioConfig AudioConfig `json:"audioConfig"`
}

type Response struct {
	AudioContent string      `json:"audioContent"`
	AudioConfig  AudioConfig `json:"audioConfig"`
}

func main() {
	router := gin.Default()

	// GET endpoint
	router.GET("/", func(c *gin.Context) {
		c.File("index.html")
	})

	// POST endpoint
	router.POST("/submit", func(c *gin.Context) {
		var msg Message
		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		payload := Payload{
			Input: Input{
				Text: msg.Text,
				Type: "text",
			},
			Voice: Voice{
				Model:  "zh_en_female_2",
				Speed:  1,
				Pitch:  1,
				Energy: 1,
			},
			AudioConfig: AudioConfig{
				Encoding:   "LINEAR16",
				SampleRate: "16K",
			},
		}

		// Convert map to JSON
		payloadData, err := json.Marshal(payload)
		if err != nil {
			panic(err)
		}

		hostURl := os.Getenv("HOST_URL")
		hostKey := os.Getenv("HOST_KEY")

		// Create a new POST request
		req, err := http.NewRequest(http.MethodPost, hostURl, bytes.NewBuffer(payloadData))
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("key", hostKey)

		// Send the request using http.Client
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		// Read and print response
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var response *Response
		err = json.Unmarshal(body, &response)
		if err != nil {
			fmt.Println("Error unmarshalling response:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process response"})
			return
		}

		c.JSON(http.StatusOK, response)
	})

	// Start server
	router.Run(":8080") // listens on port 8080
}

package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"os/exec"
	//   "fmt"
	"strings"
	//   "os"
)

type data struct {
	Ip  string `json:"ip"`
	Mac string `json:"mac"`
}

func getinfo(c *gin.Context) {

	cmd1 := exec.Command("arp", "-n", "-a")
	cmd2 := exec.Command("grep", "-v", "00:00:00:00:00:00")
	cmd3 := exec.Command("sed", "1d")
	cmd4 := exec.Command("awk", `{print $2 "\t" $4}`)
	// Pipe the cmd1 and cmd2
	cmd2.Stdin, _ = cmd1.StdoutPipe()
	cmd3.Stdin, _ = cmd2.StdoutPipe()
	cmd4.Stdin, _ = cmd3.StdoutPipe()
	stdout, err := cmd4.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	//sed and awk them to get info
	cmd4.Start()
	cmd3.Start()
	cmd2.Start()
	cmd1.Run()

	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(opBytes))
	var d = make([]data, 2, 2)
	var app data
	for _, name := range strings.Split(string(opBytes), "\n") {
		for index, info := range strings.Split(name, "\t") {
			if info != "" {
				if index%2 == 0 {
					app.Ip = info
				} else {
					app.Mac = info
				}
			d = append(d, app)
			}
		}
	}
	c.JSON(200, gin.H{
		"status":  200,
		"message": d,
	})
	return
}

func main() {
	r := gin.Default()
	r.GET("/getipinfo", getinfo)
	r.Run("192.168.1.135:8000")
}

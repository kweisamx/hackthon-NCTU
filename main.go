package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

type data struct {
	Host string `json:"host"`
	Ip   string `json:"ip"`
}

func getinfo(c *gin.Context) {

	cmd1 := exec.Command("nmap", "-sP", "192.168.1.0/24")
	cmd2 := exec.Command("grep", "Nmap")
	cmd3 := exec.Command("sed", "1d")
	cmd4 := exec.Command("sed", "$d")
	cmd5 := exec.Command("awk", `{print $5 "\t" $6}`)

	//pipe
	cmd2.Stdin, _ = cmd1.StdoutPipe()
	cmd3.Stdin, _ = cmd2.StdoutPipe()
	cmd4.Stdin, _ = cmd3.StdoutPipe()
	cmd5.Stdin, _ = cmd4.StdoutPipe()
	stdout, err := cmd5.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	//sed and awk them to get info
	cmd5.Start()
	cmd4.Start()
	cmd3.Start()
	cmd2.Start()
	cmd1.Run()

	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(opBytes))
	var d = make([]data, 0, 0)
	var app data
	for _, name := range strings.Split(string(opBytes), "\n") {
		if name != "" {
			fmt.Println("########", name)
			for index, info := range strings.Split(name, "\t") {
				if index%2 == 1 {
					app.Ip = strings.Trim(info, "(")
					app.Ip = strings.Trim(app.Ip, ")")

				} else {
					app.Host = info
				}
			}
			fmt.Println("!!!", app)
			d = append(d, app)
		}
	}
	c.JSON(200, gin.H{
		"status":  200,
		"message": d,
	})
	return
}

func gethash(c *gin.Context) {
	cmd1 := exec.Command("ipfs", "pin", "ls")
	cmd2 := exec.Command("grep", "recursive")

	cmd2.Stdin, _ = cmd1.StdoutPipe()
	stdout, _ := cmd2.StdoutPipe()

	cmd2.Start()
	cmd1.Run()

	opBytes, _ := ioutil.ReadAll(stdout)

	hashinfo := make([]string, 0, 0)
	for _, name := range strings.Split(string(opBytes), "\n") {
		if name != "" {
			hash := strings.Split(name, " ")
			hashinfo = append(hashinfo, hash[0])
		}
	}

	c.JSON(200, gin.H{
		"status": 200,
		"data":   hashinfo,
	})
	return
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/hash", gethash)
	r.GET("/ipinfo", getinfo)
	r.Run(":8000")
}

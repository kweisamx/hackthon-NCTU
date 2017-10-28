package main

import (
	"github.com/gin-gonic/gin"
	"os/exec"
	"io/ioutil"
	"strings"
)

func pooh(c *gin.Context) {
    cmd1 := exec.Command("ipfs", "pin", "ls")
	cmd2 := exec.Command("grep", "recursive")

    cmd2.Stdin ,_ = cmd1.StdoutPipe()
    stdout,_ := cmd2.StdoutPipe()
    
    cmd2.Start()
    cmd1.Run()
    
    opBytes,_ := ioutil.ReadAll(stdout)


	arrrr := make([]string, 10, 10)
	for _, name := range strings.Split(string(opBytes), "\n") {
		one := strings.Split(name, " ")
		arrrr = append(arrrr, one[0])
	}

    arrrr_next := make([]string, 0, 0)
	for _, son_zai_X := range arrrr {
	    if(son_zai_X != "") {
		    arrrr_next = append(arrrr_next, son_zai_X)
        }   
    }
    c.JSON(200, gin.H{
		"status":  200,
		"data": arrrr_next,						
    })
	return
}



func main() {
	r := gin.Default()
	r.GET("/hosthash", pooh)
	r.Run(":10000")
}

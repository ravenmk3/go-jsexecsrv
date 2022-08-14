package app

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/robertkrimen/otto"
)

func Run() {
	engine := gin.Default()
	engine.POST("/exec", handleExec)

	err := engine.Run("0.0.0.0:666")
	if err != nil {
		log.Fatal(err)
	}
}

func handleExec(c *gin.Context) {
	varName := c.Query("var")
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	script := string(body)
	value, err := execJs(script, varName)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"value": value})
}

func execJs(script string, varName string) (string, error) {
	vm := otto.New()

	value, err := vm.Run(script)
	if err != nil {
		return "", err
	}

	if varName == "" {
		return value.ToString()
	}

	value, err = vm.Get(varName)
	if err != nil {
		return "", nil
	}
	return value.ToString()
}

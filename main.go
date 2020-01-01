package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	// "golang-projects/aws_dynamodb_basic_test/dao"
)

var ctrl *DbController

func main() {
	// table := "Test"
	ctrl = InitDbConnection("http://localhost:8000")

	router := gin.New()
	router.Use(gin.Logger())
	router.GET("/user", GetUser)
	router.GET("/userlist", GetUserList)
	router.POST("/user", PutUser)

	fmt.Print("running server on localhost:5000")
	router.Run(":5000")

}
func GetUser(c *gin.Context) {
	u := c.Query("userid")
	// c.Param("userid")
	c.JSONP(200, gin.H{"data": u})
}

func GetUserList(c *gin.Context) {
	// http://github.com/gin-gonic/examples
	table := "Test"
	var resList []*UserObject // TODO cast to model struct
	err := ctrl.List(table, resList)
	if err != nil {
		c.AbortWithError(501, err)
	}
	// c.String(http.StatusOK, string(result))
	// c.HTML(http.StatusOK, "template.tmpl", gin.H{"title": "helloworld"})
	// c.Stream()
	c.JSONP(http.StatusOK, gin.H{"data": resList})
}

func PutUser(c *gin.Context) {
	// c.GetPostForm()
	u := c.PostForm("username")
	// e := c.PostForm("email")
	i := &UserObject{}
	i.Name = u
	ctrl.PutItem("Test", i)
	c.JSONP(200, gin.H{"data": u})
}

func UpdateUser(c *gin.Context) {
	u := c.PostForm("username")
	// e := c.PostForm("email")
	i := &Todo{}
	i.Action = u
	ctrl.PutItem("Test", i)
	c.JSONP(200, gin.H{"data": u})
}

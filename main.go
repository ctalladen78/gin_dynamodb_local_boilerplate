package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	// "golang-projects/aws_dynamodb_basic_test/dao"
)

var ctrl *DbController

func main() {
	ctrl = InitDbConnection("http://localhost:8000")

	router := gin.New()
	router.Use(gin.Logger())
	router.GET("/user", GetUser)
	router.GET("/userlist", GetUserList)
	router.POST("/user", PutUser)
	router.POST("/user/edit", UpdateUser)

	fmt.Print("running server on localhost:5000")
	router.Run(":5000")

}

// GET /user?userid=one
func GetUser(c *gin.Context) {
	userId := c.Query("userid")
	todo := c.Query("todo")
	// c.Param("userid")
	res, err := ctrl.GetItem(userId, todo, "Test2")
	if err != nil {
		c.AbortWithError(501, err)
	}
	c.JSONP(200, gin.H{"data": res})
}

func GetUserList(c *gin.Context) {
	// http://github.com/gin-gonic/examples
	table := "Test2"
	resList, err := ctrl.List(table)
	if err != nil {
		c.AbortWithError(501, err)
	}
	fmt.Printf("RESULTS %s", resList)
	// c.String(http.StatusOK, string(result))
	// c.HTML(http.StatusOK, "template.tmpl", gin.H{"title": "helloworld"})
	// c.Stream()
	c.JSONP(http.StatusOK, gin.H{"data": resList})
}

func PutUser(c *gin.Context) {
	// c.GetPostForm()
	u := c.PostForm("todo")
	// e := c.PostForm("email")
	t := &TodoObject{}
	t.Id = time.Now().Format(time.RFC3339) // uuid.New()
	t.Todo = u
	ctrl.PutItem("Test2", t)
	c.JSONP(200, gin.H{"data": u})
}

func UpdateUser(c *gin.Context) {
	nt := c.PostForm("newtodo")
	oid := c.PostForm("objectid")
	ot := c.PostForm("oldtodo")
	// e := c.PostForm("email")
	t := &TodoObject{
		Id:   oid,
		Todo: ot,
	}
	u, err := ctrl.Update("Test2", t, nt)
	if err != nil {
		c.AbortWithError(501, err)
	}
	c.JSONP(200, gin.H{"data": u})
}

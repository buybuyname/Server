package taskController

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swsad-dalaotelephone/Server/models/task"
	"github.com/swsad-dalaotelephone/Server/models/user"
	"github.com/swsad-dalaotelephone/Server/modules/log"
)

/*
VerifyTask : verify task
require: task_id, accepter_id, result, feedback, cookie
return: msg
*/
func VerifyTask(c *gin.Context) {
	// taskId := c.PostForm("task_id")
	// publisherId := c.PostForm("publisher_id")
	taskId := c.Param("task_id")
	user := c.MustGet("user").(userModel.User)
	publisherId := user.Id
	accepterId := c.PostForm("accepter_id")
	result := c.PostForm("result")
	feedback := c.PostForm("feedback")

	log.ErrorLog.Println(taskId)
	log.ErrorLog.Println(publisherId)
	log.ErrorLog.Println(accepterId)
	log.ErrorLog.Println(result)

	if taskId == "" || publisherId == "" || accepterId == "" || result == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "missing argument",
		})
		log.ErrorLog.Println("missing arugment")
		c.Error(errors.New("missing argument"))
		return
	}

	// check task_id exist or not
	tasks, err := taskModel.GetTasksByStrKey("id", taskId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		log.ErrorLog.Println(err)
		c.Error(err)
		return
	}

	// check accepter_id exist or not
	accepters, err := userModel.GetUsersByStrKey("id", accepterId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		log.ErrorLog.Println(err)
		c.Error(err)
		return
	}

	exist := len(tasks) == 1 && len(accepters) == 1
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid argument",
		})
		log.ErrorLog.Println("invalid argument")
		c.Error(errors.New("invalid argument"))
		return
	}
	//publisher of task is not this user
	if tasks[0].PublisherId != publisherId {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "permission denied",
		})
		log.ErrorLog.Println("permission denied")
		c.Error(errors.New("permission denied"))
		return
	}

	acceptance, err := taskModel.GetAcceptanceByTaskAccepterId(taskId, accepterId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		log.ErrorLog.Println(err)
		c.Error(err)
		return
	}
	// todo check acceptance invalid or not
	if acceptance.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "can not find acceptance record",
		})
		log.ErrorLog.Println("can not find acceptance record")
		c.Error(errors.New("can not find acceptance record"))
		return
	}

	if result == "true" {
		acceptance.Status = taskModel.StatusAcceptFinished
	} else {
		acceptance.Status = taskModel.StatusAcceptUnqualified
	}
	acceptance.Feedback = feedback

	if err := taskModel.UpdateAcceptance(acceptance); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
		log.InfoLog.Println(err)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		log.ErrorLog.Println(err)
		c.Error(err)
		return
	}
}

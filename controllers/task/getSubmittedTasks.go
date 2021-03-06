package taskController

import (
	"errors"
	"net/http"

	"github.com/swsad-dalaotelephone/Server/models/task"
	"github.com/swsad-dalaotelephone/Server/models/user"
	"github.com/swsad-dalaotelephone/Server/modules/log"
	"github.com/swsad-dalaotelephone/Server/modules/util"

	"github.com/gin-gonic/gin"
)

/*
GetSubmittedTasks : get submitted task
require: task_id, cookie
return: submitted task list
*/
func GetSubmittedTasks(c *gin.Context) {

	// taskId := c.Query("task_id")
	// publisherId := c.Query("publisher_id")
	taskId := c.Param("task_id")
	user := c.MustGet("user").(userModel.User)
	publisherId := user.Id

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

	if len(tasks) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid argument",
		})
		log.ErrorLog.Println("invalid argument")
		c.Error(err)
		return
	}

	if publisherId != tasks[0].PublisherId {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "permission denied",
		})
		log.ErrorLog.Println("permission denied")
		c.Error(errors.New("permission denied"))
		return
	}

	// get acceptances list
	acceptances, err := taskModel.GetAcceptancesByStrKey("task_id", taskId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "can not fetch acceptance list",
		})
		log.ErrorLog.Println(err)
		c.Error(err)
		return
	}

	submitted := make([]taskModel.Acceptance, 0)
	for _, item := range acceptances {
		if item.Status >= 1 {
			submitted = append(submitted, item)
		}
	}

	submittedJson, err := util.StructToJsonStr(submitted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "json convert error",
		})
		log.ErrorLog.Println(err)
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"submitted": submittedJson,
	})
	log.InfoLog.Println(publisherId, len(submitted), "success")
}

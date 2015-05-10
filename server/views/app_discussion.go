package views

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/guregu/kami"
	"github.com/shumipro/meetapp/server/models"
	"golang.org/x/net/context"
)

func init() {
	// Discussion API
	kami.Post("/u/api/app/discussion", APIAppDiscussion)
	kami.Delete("/u/api/app/discussion/:id", APIAppDiscussionDelete)
}

type DiscussionRequest struct {
	AppID          string `json:"appId"` // アプリID
	DiscussionInfo models.DiscussionInfo
}

func APIAppDiscussion(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err.Error())
		return
	}
	fmt.Println(string(data))

	// convert request params to struct
	var discussionReq DiscussionRequest
	if err := json.Unmarshal(data, &discussionReq); err != nil {
		log.Println("ERROR! json parse", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	// get appinfo from db
	appInfo, err := models.AppsInfoTable.FindID(ctx, discussionReq.AppID)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	nowTime := time.Now()
	discussionReq.DiscussionInfo.ID = strconv.FormatInt(time.Now().UnixNano(), 10)
	discussionReq.DiscussionInfo.Timestamp = nowTime
	// push a discussionInfo
	appInfo.Discussions = append(appInfo.Discussions, discussionReq.DiscussionInfo)
	appInfo.UpdateAt = nowTime

	if err := models.AppsInfoTable.Upsert(ctx, appInfo); err != nil {
		log.Println("ERROR! discussion", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	notification := models.Notification{}
	notification.NotificationID = discussionReq.DiscussionInfo.ID
	notification.SourceID = discussionReq.DiscussionInfo.ID
	notification.NotificationType = models.NotificationDiscussion
	notification.Message = "メッセージがあります"
	notification.IsRead = false

	go func() {
		for _, m := range appInfo.Members {
			err := models.NotificationTable.AddNotification(ctx, m.UserID, notification)
			if err != nil {
				log.Println("ERROR: AddNotification", m.UserID, notification)
				continue
			}
			log.Println("OK: AddNotification", m.UserID, notification)
		}
	}()

	renderer.JSON(w, 200, appInfo.Discussions)
}

func APIAppDiscussionDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	discussionID := kami.Param(ctx, "id")

	if err := models.AppsInfoTable.DeleteDiscussionByID(ctx, discussionID); err != nil {
		log.Println("ERROR! discussion", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	renderer.JSON(w, 200, discussionID)
}
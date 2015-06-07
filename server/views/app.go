package views

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/guregu/kami"
	"golang.org/x/net/context"

	"github.com/shumipro/meetapp/server/constants"
	"github.com/shumipro/meetapp/server/login"
	"github.com/shumipro/meetapp/server/models"
	"github.com/shumipro/meetapp/server/twitter"
)

var notAdminError = fmt.Errorf("%s", "not admin user")

var sortLabels = map[string]map[string]string{
	"new": {
		"title": "新着アプリ",
	},
	"popular": {
		"title": "人気アプリ",
	},
	"updateAt": {
		"title": "開発アイデアを探す",
	},
}

func init() {
	kami.Get("/app/detail/:id", AppDetail)
	kami.Get("/app/list", AppList)
	kami.Get("/u/app/register", AppRegister)
	kami.Get("/u/app/edit/:id", AppEdit)
	// Apps API
	//	kami.Get("/u/api/app/apps", APIAppGetAll)
	//	kami.Get("/u/api/app/apps/:id", APIAppGet)
	kami.Post("/u/api/app/register", APIAppRegister)   // TODO: [POST] /u/api/app/apps
	kami.Put("/u/api/app/edit/:id", APIAppEdit)        // TODO: [PUT] /u/api/app/apps/:id
	kami.Delete("/u/api/app/delete/:id", APIAppDelete) // TODO: [DELETE] /u/api/app/apps/:id
}

type AppListResponse struct {
	TemplateHeader
	AppInfoList []AppInfoView
	CurrentPage int
	PerPageNum  int
	TotalCount  int
}

const perPageNum int = 10

func AppList(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	orderBy := r.FormValue("orderBy")
	if orderBy == "" {
		orderBy = string(models.OrderByUpdateAt) // デフォルトはUpdateAt
	}

	page, _ := strconv.Atoi(r.FormValue("page"))
	if page > 0 {
		page -= 1 // 1ページは0とする
	}

	platform := r.FormValue("platform")
	occupation := r.FormValue("occupation")
	category := r.FormValue("category")
	pLang := r.FormValue("pLang")
	area := r.FormValue("area")

	filter := models.AppInfoFilter{}
	filter.OccupationType = constants.OccupationType(occupation)
	filter.CategoryType = constants.CategoryType(category)
	filter.LanguageType = constants.LanguageType(pLang)
	filter.AreaType = constants.AreaType(area)
	filter.PlatformType = constants.PlatformType(platform)
	filter.OrderBy = models.AppInfoOrderType(orderBy)

	preload := AppListResponse{}
	preload.TemplateHeader = NewHeader(ctx,
		"MeetApp - "+sortLabels[orderBy]["title"],
		"",
		"気になるアプリ開発に参加しよう",
		false,
		"",
		"",
	)

	// ViewModel変換して詰める
	totalCount, apps, err := models.AppsInfoTable.FindFilter(ctx, filter, page*perPageNum, perPageNum)
	if err != nil {
		panic(err)
	}
	preload.AppInfoList = convertAppInfoViewList(ctx, apps)
	preload.PerPageNum = perPageNum
	preload.CurrentPage = page + 1
	preload.TotalCount = totalCount

	ExecuteTemplate(ctx, w, r, "app/list", preload)
}

type AppDetailResponse struct {
	TemplateHeader
	AppInfo AppInfoView
}

func AppDetail(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	appID := kami.Param(ctx, "id")
	// TODO: とりあえず
	if appID == "favicon.png" || appID == "" {
		return
	}
	appInfo, err := models.AppsInfoTable.FindID(ctx, appID)
	if err != nil {
		http.Redirect(w, r, "/error", 302)
		return
	}

	preload := AppDetailResponse{}
	preload.TemplateHeader = NewHeader(ctx,
		"MeetApp - "+appInfo.Name,
		appInfo.Description,
		appInfo.Name,
		false,
		// set current relative URL
		r.URL.RequestURI(),
		appInfo.MainImage,
	)
	preload.AppInfo = NewAppInfoView(ctx, appInfo)

	ExecuteTemplate(ctx, w, r, "app/detail", preload)
}

type AppRegisterResponse struct {
	TemplateHeader
	AppInfo AppInfoView
}

func AppRegister(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	preload := AppRegisterResponse{}
	preload.TemplateHeader = NewHeader(ctx,
		"MeetApp - アプリの登録",
		"",
		"アプリを登録して仲間を探そう",
		false,
		"",
		"",
	)

	// 自分をデフォルトメンバーとして突っ込んでおく
	a, _ := login.FromContext(ctx)
	appInfo := models.AppInfo{}
	members := []models.Member{
		{
			UserID:     a.UserID,
			Occupation: constants.OccupationType("1"),
			IsAdmin:    true,
		},
	}
	appInfo.Members = members
	// sizeを3にする
	appInfo.ImageURLs = make([]models.URLInfo, 3)

	preload.AppInfo = NewAppInfoView(ctx, appInfo)

	ExecuteTemplate(ctx, w, r, "app/register", preload)
}

func AppEdit(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	appID := kami.Param(ctx, "id")

	// TODO: とりあえず
	if appID == "favicon.png" || appID == "" {
		return
	}

	appInfo, err := models.AppsInfoTable.FindID(ctx, appID)
	if err != nil {
		http.Redirect(w, r, "/error", 302)
		return
	}

	preload := AppRegisterResponse{}
	preload.TemplateHeader = NewHeader(ctx,
		"MeetApp - アプリの編集",
		"",
		"アプリを登録して仲間を探そう",
		false,
		"",
		"",
	)
	// sizeを3にする
	if len(appInfo.ImageURLs) != 3 {
		imageURLs := make([]models.URLInfo, 3)
		for idx, img := range appInfo.ImageURLs {
			imageURLs[idx] = img
		}
		appInfo.ImageURLs = imageURLs
	}
	preload.AppInfo = NewAppInfoView(ctx, appInfo)

	ExecuteTemplate(ctx, w, r, "app/register", preload)
}

func APIAppRegister(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	regAppInfo, err := readBodyAppInfo(r.Body)
	if err != nil {
		renderer.JSON(w, 400, "[ERROR] request param appInfo "+err.Error())
		return
	}

	regAppInfo = convertRegisterAppInfo(ctx, regAppInfo)
	if err := models.AppsInfoTable.Upsert(ctx, regAppInfo); err != nil {
		panic(err)
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		// Initialize the Twitter Client
		twClient, ok := twitter.FromContext(ctx)
		if !ok {
			log.Printf("Failed to initialize twitter client: %s.", err)
			return
		}

		message := fmt.Sprintf(
			"開発アイデアが新規登録されました: MeetApp - %s https://meetapp.tokyo/app/detail/%s #meetapp",
			regAppInfo.Name,
			regAppInfo.ID,
		)
		id, err := twClient.Tweet(message)
		if err != nil {
			log.Printf("Failed to post a tweet for %s: %s.", regAppInfo.ID, err)
			return
		}
		log.Printf("Successfully posted a tweet %s.", id)
	}()

	renderer.JSON(w, 200, regAppInfo)
}

func APIAppEdit(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	regAppInfo, err := readBodyAppInfo(r.Body)
	if err != nil || regAppInfo.ID == "" {
		renderer.JSON(w, 400, "[ERROR] request param appInfo ")
		return
	}

	// アプリが存在しないか
	beforeApp, err := models.AppsInfoTable.FindID(ctx, regAppInfo.ID)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	// 管理者じゃないアプリか
	a, _ := login.FromContext(ctx)
	if !beforeApp.IsAdmin(a.UserID) {
		log.Println("ERROR!", notAdminError)
		renderer.JSON(w, 400, notAdminError.Error())
		return
	}

	regAppInfo = convertEditAppInfo(ctx, regAppInfo, beforeApp)

	if err := models.AppsInfoTable.Upsert(ctx, regAppInfo); err != nil {
		panic(err)
	}

	renderer.JSON(w, 200, regAppInfo)
}

func APIAppDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	a, _ := login.FromContext(ctx)
	appID := kami.Param(ctx, "id")

	app, err := models.AppsInfoTable.FindID(ctx, appID)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	// 管理者のみ削除可能
	if !app.IsAdmin(a.UserID) {
		log.Println("ERROR!", notAdminError)
		renderer.JSON(w, 400, notAdminError.Error())
		return
	}

	if err := models.AppsInfoTable.Delete(ctx, app.ID); err != nil {
		panic(err)
	}

	renderer.JSON(w, 200, appID)
}

package views

import (
	"log"
	"net/http"

	"github.com/guregu/kami"
	"github.com/huandu/facebook"
	"github.com/shumipro/meetapp/server/models"
	"github.com/shumipro/meetapp/server/oauth"
	"golang.org/x/net/context"
)

func init() {
	kami.Get("/mypage/other/:id", MypageOther)

	kami.Get("/u/mypage", Mypage)
	kami.Get("/u/mypage/edit", MypageEdit)
}

type MyPageResponse struct {
	TemplateHeader
	User         models.User
	AdminAppList []AppInfoView
	JoinAppList  []AppInfoView
	IsMe         bool
}

func Mypage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	a, ok := oauth.FromContext(ctx)
	if !ok {
		panic("login error")
	}

	// マイページ表示時はFacebookとかの情報とりなおしてUserテーブル更新する
	_, err := facebook.Get("/me", facebook.Params{
		"access_token": a.AuthToken,
	})
	if err != nil {
		log.Println(err)
	}

	// TODO: Facebook情報に変更あればUserテーブル更新する

	user, err := models.UsersTable.FindID(ctx, a.UserID)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	preload := MyPageResponse{}
	preload.User = user
	preload.TemplateHeader = NewHeader(ctx, "マイページ", "", "", false, "", "")

	adminApps, _ := models.AppsInfoTable.FindByAdminID(ctx, a.UserID)
	joinApps, _ := models.AppsInfoTable.FindByJoinID(ctx, a.UserID)
	preload.AdminAppList = convertAppInfoViewList(ctx, adminApps)
	preload.JoinAppList = convertAppInfoViewList(ctx, joinApps)
	preload.IsMe = true

	ExecuteTemplate(ctx, w, r, "mypage", preload)
}

func MypageOther(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID := kami.Param(ctx, "id")
	if userID == "" || userID == "favicon.png" {
		return
	}

	user, err := models.UsersTable.FindID(ctx, userID)
	if err != nil {
		log.Println(err, userID)
		renderer.JSON(w, 400, err.Error())
		return
	}

	preload := MyPageResponse{}
	preload.User = user
	preload.TemplateHeader = NewHeader(ctx, user.Name, "", "", false, "", "")

	// loginしている状態のみ他人のページとして自分を見に来たときに管理アイデアを表示
	adminApps, _ := models.AppsInfoTable.FindByAdminID(ctx, userID)
	joinApps, _ := models.AppsInfoTable.FindByJoinID(ctx, userID)
	preload.AdminAppList = convertAppInfoViewList(ctx, adminApps)
	preload.JoinAppList = convertAppInfoViewList(ctx, joinApps)

	// loginしてる状態でotherが自分のページであればIsMeにtrueをセット
	a, ok := oauth.FromContext(ctx)
	preload.IsMe = ok && userID == a.UserID

	ExecuteTemplate(ctx, w, r, "mypage", preload)
}

func MypageEdit(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	a, ok := oauth.FromContext(ctx)
	if !ok {
		panic("login error")
	}

	user, err := models.UsersTable.FindID(ctx, a.UserID)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	preload := MyPageResponse{}
	preload.User = user
	preload.TemplateHeader = NewHeader(ctx, "マイページの編集", "", "", false, "", "")

	ExecuteTemplate(ctx, w, r, "mypageEdit", preload)
}

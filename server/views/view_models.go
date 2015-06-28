package views

import (
	"github.com/shumipro/meetapp/server/login"
	"github.com/shumipro/meetapp/server/models"
	"golang.org/x/net/context"
)

type AppInfoView struct {
	models.AppInfo
	Members     []UserMember      // models.Membersを上書きします
	Discussions []UserDiscussions // models.Discussionsを上書きします
	StarUsers   []models.User     // models.Discussionsを上書きします
	Stared      bool              // 現在認証されているユーザーがstarしているかどうか
	IsAdmin     bool              // 管理者かどうか
}

// TODO: あとでRenameする（Emptyというよりは未登録）
func (a AppInfoView) IsEmpty() bool {
	return a.AppInfo.ID == ""
}

// UserMember User情報を持つMember
type UserMember struct {
	models.Member
	models.User
}

type UserDiscussions struct {
	models.DiscussionInfo
	models.User
	Deletable bool // 削除できるかどうか
}

func NewAppInfoView(ctx context.Context, appInfo models.AppInfo) AppInfoView {
	a := AppInfoView{}
	a.AppInfo = appInfo

	account, ok := login.FromContext(ctx)
	if ok {
		a.IsAdmin = a.AppInfo.IsAdmin(account.UserID)
		a.Stared = a.AppInfo.Stared(account.UserID)
	}

	a.Members = make([]UserMember, len(a.AppInfo.Members))
	for idx, m := range appInfo.Members {
		// TODO: あとでIn句にして1クエリにする
		u, _ := models.UsersTable.FindID(ctx, m.UserID)
		a.Members[idx] = UserMember{Member: m, User: u}
	}

	a.Discussions = make([]UserDiscussions, len(a.AppInfo.Discussions))
	for idx, d := range appInfo.Discussions {
		// TODO: あとでIn句にして1クエリにする
		u, _ := models.UsersTable.FindID(ctx, d.UserID)
		a.Discussions[idx] = UserDiscussions{DiscussionInfo: d, User: u, Deletable: d.UserID == account.UserID}
	}

	// starしたユーザーの一覧表示用
	a.StarUsers = make([]models.User, len(a.AppInfo.StarUsers))
	for idx, s := range appInfo.StarUsers {
		// TODO: あとでIn句にして1クエリにする
		u, _ := models.UsersTable.FindID(ctx, s)
		a.StarUsers[idx] = u
	}

	// projectStateあと追加なので未設定の場合は募集中にする対応
	if string(a.ProjectState) == "" {
		a.ProjectState = "1" // "募集中"
	}

	return a
}

func convertAppInfoViewList(ctx context.Context, apps []models.AppInfo) []AppInfoView {
	appViews := make([]AppInfoView, len(apps))
	for idx, app := range apps {
		appViews[idx] = NewAppInfoView(ctx, app)
	}
	return appViews
}

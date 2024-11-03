package admin

import (
	"net/http"

	"github.com/qor5/admin/seo"

	"github.com/qor/oss/filesystem"
	"github.com/qor5/admin/activity"
	"github.com/qor5/admin/l10n"
	l10n_view "github.com/qor5/admin/l10n/views"
	media_view "github.com/qor5/admin/media/views"
	"github.com/qor5/admin/pagebuilder"
	"github.com/qor5/admin/pagebuilder/example"
	"github.com/qor5/admin/presets"
	"github.com/qor5/admin/presets/gorm2op"
	"github.com/qor5/admin/publish"
	publish_view "github.com/qor5/admin/publish/views"
	"github.com/qor5/admin/utils"
	"github.com/qor5/ui/vuetify"
	"github.com/qor5/web"
	"github.com/qor5/x/login"
	"github.com/qor5/x/perm"
	h "github.com/theplant/htmlgo"
	"golang.org/x/text/language"
)

const (
	PublishDir = "./publish"
)

type Config struct {
	pb          *presets.Builder
	pageBuilder *pagebuilder.Builder
}

func InitApp() *http.ServeMux {
	c := newPB()
	mux := SetupRouter(c)

	return mux
}

func newPB() Config {
	db := ConnectDB()

	b := presets.New().VuetifyOptions(`
{
  icons: {
	iconfont: 'md', // 'mdi' || 'mdiSvg' || 'md' || 'fa' || 'fa4'
  },
  theme: {
    themes: {
      light: {
		  primary: "#ed6f2d",
		  secondary: "#009688",
		  accent: "#ff5722",
		  error: "#f44336",
		  warning: "#ff9800",
		  info: "#8bc34a",
		  success: "#4caf50"
      },
    },
  },
}
`)

	b.URIPrefix("/admin").DataOperator(gorm2op.DataOperator(db)).
		BrandFunc(func(ctx *web.EventContext) h.HTMLComponent {
			return vuetify.VCardText(
				h.H1("Admin").Style("color: red;"),
			).Class("pa-0")
		}).
		HomePageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
			r.Body = vuetify.VContainer(
				h.H1("Home"),
				h.P().Text("Change your home page here"))
			return
		})

	b.Permission(
		perm.New().Policies(
			perm.PolicyFor(perm.Anybody).WhoAre(perm.Allowed).ToDo(presets.PermCreate, presets.PermUpdate, presets.PermDelete, presets.PermGet, presets.PermList).On("*"),
			perm.PolicyFor(perm.Anybody).WhoAre(perm.Denied).ToDo(presets.PermCreate, presets.PermUpdate, presets.PermDelete).On("*:activity_logs:*"),
		),
	)

	utils.Configure(b)
	media_view.Configure(b, db)
	ab := activity.New(b, db).SetCreatorContextKey(login.UserKey)
	l10nBuilder := l10n.New()

	pageBuilder := example.ConfigPageBuilder(db, "/admin/page_builder", ``, b.I18n())
	storage := filesystem.New(PublishDir)
	publisher := publish.New(db, storage).WithPageBuilder(pageBuilder)

	seoBuilder := seo.NewBuilder(db)
	pm := pageBuilder.Configure(b, db, l10nBuilder, ab, publisher, seoBuilder)
	tm := pageBuilder.ConfigTemplate(b, db)
	cm := pageBuilder.ConfigCategory(b, db, l10nBuilder)

	ab.RegisterModels(pm, tm, cm)

	publish_view.Configure(b, db, ab, publisher, pm)

	l10nBuilder.
		RegisterLocales("International", "International", "International").
		RegisterLocales("China", "China", "China").
		GetSupportLocaleCodesFromRequestFunc(func(R *http.Request) []string {
			return l10nBuilder.GetSupportLocaleCodes()[:]
		})
	l10n_view.Configure(b, db, l10nBuilder, ab, pm)

	b.I18n().
		SupportLanguages(language.English, language.SimplifiedChinese).
		RegisterForModule(language.English, I18nExampleKey, Messages_en_US).
		RegisterForModule(language.SimplifiedChinese, I18nExampleKey, Messages_zh_CN).
		RegisterForModule(language.SimplifiedChinese, presets.ModelsI18nModuleKey, Messages_zh_CN_ModelsI18nModuleKey).
		GetSupportLanguagesFromRequestFunc(func(r *http.Request) []language.Tag {
			return b.I18n().GetSupportLanguages()
		})

	b.MenuOrder(
		b.MenuGroup("Page Builder").SubItems("pages", "page_templates", "page_categories").Icon("web"),
		"shared_containers",
		"demo_containers",
		"media-library",
	)

	initMediaLibraryData(db)
	initWebsiteData(db)

	return Config{
		pb:          b,
		pageBuilder: pageBuilder,
	}
}

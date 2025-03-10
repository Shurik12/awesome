package admin

import (
	"net/http"

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
	"github.com/qor5/admin/seo"
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
						primary: "#0c80d7", 
						secondary: "#eabd34",
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
			logo := `data:image/jpeg;base64,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`
			return vuetify.VContainer(
				h.Img(logo).Attr("width", "150"),
			).Class("ma-n4")
		}).
		HomePageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
			r.Body = vuetify.VContainer(
				h.H1("Home"),
				h.P().Text("Awesome1"))
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

	seoBuilder := seo.NewBuilder(db, seo.WithLocales("International"))
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

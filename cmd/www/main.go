package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/kammeph/go-htmx/internal/data"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const (
	baseTemplateName              = "base.html"
	baseTemplatePath              = "templates/base.html"
	menuTemplatePath              = "templates/partials/menu.html"
	typesMenuTemplatePath         = "templates/partials/types-menu.html"
	componentsMenuTemplatePath    = "templates/partials/components-menu.html"
	gearboxesTemplatePath         = "templates/pages/gearboxes.html"
	gearboxTypesTemplatePath      = "templates/pages/types/gearbox-types.html"
	gearTypesTemplatePath         = "templates/pages/types/gear-types.html"
	housingTypesTemplatePath      = "templates/pages/types/housing-types.html"
	polygonTypesTemplatePath      = "templates/pages/types/polygon-types.html"
	gearComponentsTemplatePath    = "templates/pages/components/gear-components.html"
	housingComponentsTemplatePath = "templates/pages/components/housing-components.html"
	polygonComponentsTemplatePath = "templates/pages/components/polygon-components.html"
	gearboxEditTemplatePath       = "templates/pages/gearbox-edit.html"
)

func TranslatorMiddleware(t *Translator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			locale := chi.URLParam(r, "locale")
			t.SetLocale(locale)
			next.ServeHTTP(w, r)
		})
	}
}

type Translator struct {
	localizers map[string]*i18n.Localizer
	locale     string
}

func NewTranslator(bundle *i18n.Bundle, locales []string, defaultLocale string) *Translator {
	localizers := map[string]*i18n.Localizer{}
	for _, locale := range locales {
		localizers[locale] = i18n.NewLocalizer(bundle, locale)
	}
	return &Translator{localizers, defaultLocale}
}

func (t *Translator) Translate(messageID string) string {
	translation, _ := t.localizers[t.locale].Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	return translation
}

func (t *Translator) GetLocale() string {
	return t.locale
}

func (t *Translator) SetLocale(locale string) {
	t.locale = locale
}

func main() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile("i18n/en.json")
	bundle.MustLoadMessageFile("i18n/de.json")

	translator := NewTranslator(bundle, []string{"en", "de"}, "en")
	funcs := template.FuncMap{"translate": translator.Translate}

	gearboxesTemplate := template.Must(template.New(baseTemplateName).Funcs(funcs).ParseFiles(baseTemplatePath, menuTemplatePath, gearboxesTemplatePath))
	gearboxTypeTemplate := template.Must(template.New(baseTemplateName).Funcs(funcs).ParseFiles(baseTemplatePath, menuTemplatePath, typesMenuTemplatePath, gearboxTypesTemplatePath))
	gearTypeTemplate := template.Must(template.New(baseTemplateName).Funcs(funcs).ParseFiles(baseTemplatePath, menuTemplatePath, typesMenuTemplatePath, gearTypesTemplatePath))
	housingTypeTemplate := template.Must(template.New(baseTemplateName).Funcs(funcs).ParseFiles(baseTemplatePath, menuTemplatePath, typesMenuTemplatePath, housingTypesTemplatePath))
	polygonTypeTemplate := template.Must(template.New(baseTemplateName).Funcs(funcs).ParseFiles(baseTemplatePath, menuTemplatePath, typesMenuTemplatePath, polygonTypesTemplatePath))
	gearComponentTemplate := template.Must(template.New(baseTemplateName).Funcs(funcs).ParseFiles(baseTemplatePath, menuTemplatePath, componentsMenuTemplatePath, gearComponentsTemplatePath))
	housingComponentTemplate := template.Must(template.New(baseTemplateName).Funcs(funcs).ParseFiles(baseTemplatePath, menuTemplatePath, componentsMenuTemplatePath, housingComponentsTemplatePath))
	polygonComponentTemplate := template.Must(template.New(baseTemplateName).Funcs(funcs).ParseFiles(baseTemplatePath, menuTemplatePath, componentsMenuTemplatePath, polygonComponentsTemplatePath))
	gearboxEditTemplate := template.Must(template.New(baseTemplateName).Funcs(funcs).ParseFiles(baseTemplatePath, menuTemplatePath, gearboxEditTemplatePath))

	r := chi.NewRouter()
	r.Use(middleware.Compress(5))
	fs := http.FileServer(http.Dir("./public"))
	r.Handle("/public/*", http.StripPrefix("/public", fs))
	r.Handle("/manifest.json", fs)
	r.Handle("/service-worker.js", fs)
	r.Handle("/styles/*", http.StripPrefix("/styles", http.FileServer(http.Dir("./styles"))))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/en/gearboxes", http.StatusFound)
	})
	r.Route("/{locale}", func(r chi.Router) {
		r.Use(TranslatorMiddleware(translator))
		r.Route("/types", func(r chi.Router) {
			r.Get("/gearbox", func(w http.ResponseWriter, r *http.Request) {
				boosted := r.Header.Get("HX-Boosted")
				data := map[string]interface{}{"Locale": translator.GetLocale()}
				if boosted == "true" {
					gearboxTypeTemplate.ExecuteTemplate(w, "page", data)
				} else {
					gearboxTypeTemplate.Execute(w, data)
				}
			})
			r.Get("/housing", func(w http.ResponseWriter, r *http.Request) {
				boosted := r.Header.Get("HX-Boosted")
				data := map[string]interface{}{"Locale": translator.GetLocale()}
				if boosted == "true" {
					housingTypeTemplate.ExecuteTemplate(w, "page", data)
				} else {
					housingTypeTemplate.Execute(w, data)
				}
			})
			r.Get("/polygon", func(w http.ResponseWriter, r *http.Request) {
				boosted := r.Header.Get("HX-Boosted")
				data := map[string]interface{}{"Locale": translator.GetLocale()}
				if boosted == "true" {
					polygonTypeTemplate.ExecuteTemplate(w, "page", data)
				} else {
					polygonTypeTemplate.Execute(w, data)
				}
			})
			r.Get("/gear", func(w http.ResponseWriter, r *http.Request) {
				boosted := r.Header.Get("HX-Boosted")
				data := map[string]interface{}{"Locale": translator.GetLocale()}
				if boosted == "true" {
					gearTypeTemplate.ExecuteTemplate(w, "page", data)
				} else {
					gearTypeTemplate.Execute(w, data)
				}
			})
		})
		r.Route("/components", func(r chi.Router) {
			r.Get("/housing", func(w http.ResponseWriter, r *http.Request) {
				boosted := r.Header.Get("HX-Boosted")
				data := map[string]interface{}{"Locale": translator.GetLocale()}
				if boosted == "true" {
					housingComponentTemplate.ExecuteTemplate(w, "page", data)
				} else {
					housingComponentTemplate.Execute(w, data)
				}
			})
			r.Get("/polygon", func(w http.ResponseWriter, r *http.Request) {
				boosted := r.Header.Get("HX-Boosted")
				data := map[string]interface{}{"Locale": translator.GetLocale()}
				if boosted == "true" {
					polygonComponentTemplate.ExecuteTemplate(w, "page", data)
				} else {
					polygonComponentTemplate.Execute(w, data)
				}
			})
			r.Get("/gear", func(w http.ResponseWriter, r *http.Request) {
				boosted := r.Header.Get("HX-Boosted")
				data := map[string]interface{}{"Locale": translator.GetLocale()}
				if boosted == "true" {
					gearComponentTemplate.ExecuteTemplate(w, "page", data)
				} else {
					gearComponentTemplate.Execute(w, data)
				}
			})
		})
		r.Route("/gearboxes", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				// boosted := r.Header.Get("HX-Boosted")
				target := r.Header.Get("HX-Target")
				data := map[string]interface{}{"Gearboxes": data.Gearboxes, "Locale": translator.GetLocale()}
				// if boosted == "true" {
				if target == "page" {
					gearboxesTemplate.ExecuteTemplate(w, "page", data)
				} else {
					gearboxesTemplate.Execute(w, data)
				}
			})
			r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
				target := r.Header.Get("HX-Target")
				id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 0)
				var selectedgearbox data.Gearbox
				for _, gearbox := range data.Gearboxes {
					if gearbox.ID == id {
						selectedgearbox = gearbox
						break
					}
				}
				data := map[string]interface{}{"Gearbox": selectedgearbox, "Locale": translator.GetLocale()}
				if target == "page" {
					gearboxEditTemplate.ExecuteTemplate(w, "page", data)
				} else {
					gearboxEditTemplate.Execute(w, data)
				}
			})
		})
	})
	r.Get("/change-locale", func(w http.ResponseWriter, r *http.Request) {
		locale := r.URL.Query().Get("locale")
		rest := strings.Join(strings.Split(r.Referer(), "/")[4:], "/")
		w.Header().Set("HX-Redirect", fmt.Sprintf("/%s/%s", locale, rest))
		w.WriteHeader(http.StatusFound)
	})
	fmt.Println(http.ListenAndServe(":8080", r))
}

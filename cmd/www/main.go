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
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

const (
	baseTemplateName              = "base.html"
	baseTemplatePath              = "templates/base.html"
	menuTemplatePath              = "templates/partials/menu.html"
	typesMenuTemplatePath         = "templates/partials/types-menu.html"
	componentsMenuTemplatePath    = "templates/partials/components-menu.html"
	gearboxesTemplatePath         = "templates/pages/gearboxes.html"
	gearComponentsTemplatePath    = "templates/pages/components/gear-components.html"
	housingComponentsTemplatePath = "templates/pages/components/housing-components.html"
	polygonComponentsTemplatePath = "templates/pages/components/polygon-components.html"
	gearboxEditTemplatePath       = "templates/pages/gearbox-edit.html"
	housingEditTemplatePath       = "templates/pages/components/housing-edit.html"
	polygonEditTemplatePath       = "templates/pages/components/polygon-edit.html"
	gearEditTemplatePath          = "templates/pages/components/gear-edit.html"
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
	printers   map[string]*message.Printer
	locale     string
}

func NewTranslator(bundle *i18n.Bundle, locales []string, defaultLocale string) *Translator {
	localizers := map[string]*i18n.Localizer{}
	printers := map[string]*message.Printer{}
	for _, locale := range locales {
		localizers[locale] = i18n.NewLocalizer(bundle, locale)
		printers[locale] = message.NewPrinter(language.Make(locale))
	}
	return &Translator{localizers, printers, defaultLocale}
}

func (t *Translator) Translate(messageID string) string {
	translation, err := t.localizers[t.locale].Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	if err != nil {
		fmt.Println(err)
		return messageID
	}
	return translation
}

func (t *Translator) Number(input interface{}) string {
	return t.printers[t.locale].Sprint(number.Decimal(input))
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
	funcs := template.FuncMap{"translate": translator.Translate, "number": translator.Number}

	gearboxesTemplate := template.Must(template.New(baseTemplateName).Funcs(funcs).ParseFiles(baseTemplatePath, menuTemplatePath, gearboxesTemplatePath))
	gearComponentTemplate := template.Must(template.New(baseTemplateName).Funcs(funcs).ParseFiles(baseTemplatePath, menuTemplatePath, componentsMenuTemplatePath, gearComponentsTemplatePath))
	housingComponentTemplate := template.Must(template.New(baseTemplateName).Funcs(funcs).ParseFiles(baseTemplatePath, menuTemplatePath, componentsMenuTemplatePath, housingComponentsTemplatePath))
	polygonComponentTemplate := template.Must(template.New(baseTemplateName).Funcs(funcs).ParseFiles(baseTemplatePath, menuTemplatePath, componentsMenuTemplatePath, polygonComponentsTemplatePath))
	gearboxEditTemplate := template.Must(template.New(baseTemplateName).Funcs(funcs).ParseFiles(baseTemplatePath, menuTemplatePath, gearboxEditTemplatePath))
	housingEditTemplate := template.Must(template.New(baseTemplateName).Funcs(funcs).ParseFiles(baseTemplatePath, menuTemplatePath, housingEditTemplatePath))
	polygonEditTemplate := template.Must(template.New(baseTemplateName).Funcs(funcs).ParseFiles(baseTemplatePath, menuTemplatePath, polygonEditTemplatePath))
	gearEditTemplate := template.Must(template.New(baseTemplateName).Funcs(funcs).ParseFiles(baseTemplatePath, menuTemplatePath, gearEditTemplatePath))

	r := chi.NewRouter()
	r.Use(middleware.Logger)
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
		r.Route("/components", func(r chi.Router) {
			r.Route("/housing", func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					boosted := r.Header.Get("HX-Boosted")
					data := map[string]interface{}{"Housings": data.Housings, "Locale": translator.GetLocale()}
					if boosted == "true" {
						housingComponentTemplate.ExecuteTemplate(w, "page", data)
					} else {
						housingComponentTemplate.Execute(w, data)
					}
				})
				r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
					target := r.Header.Get("HX-Target")
					id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 0)
					var selectedHousing data.Housing
					for _, housing := range data.Housings {
						if housing.ID == id {
							selectedHousing = housing
							break
						}
					}
					data := map[string]interface{}{"Housing": selectedHousing, "Locale": translator.GetLocale()}
					if target == "page" {
						housingEditTemplate.ExecuteTemplate(w, "page", data)
					} else {
						housingEditTemplate.Execute(w, data)
					}
				})
			})
			r.Route("/polygon", func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					boosted := r.Header.Get("HX-Boosted")
					data := map[string]interface{}{"Polygons": data.Polygons, "Locale": translator.GetLocale()}
					if boosted == "true" {
						polygonComponentTemplate.ExecuteTemplate(w, "page", data)
					} else {
						polygonComponentTemplate.Execute(w, data)
					}
				})
				r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
					target := r.Header.Get("HX-Target")
					id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 0)
					var selectedPolygon data.Polygon
					for _, polygon := range data.Polygons {
						if polygon.ID == id {
							selectedPolygon = polygon
							break
						}
					}
					data := map[string]interface{}{"Polygon": selectedPolygon, "Locale": translator.GetLocale()}
					if target == "page" {
						polygonEditTemplate.ExecuteTemplate(w, "page", data)
					} else {
						polygonEditTemplate.Execute(w, data)
					}
				})
			})
			r.Route("/gear", func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					boosted := r.Header.Get("HX-Boosted")
					data := map[string]interface{}{"Gears": data.Gears, "Locale": translator.GetLocale()}
					if boosted == "true" {
						gearComponentTemplate.ExecuteTemplate(w, "page", data)
					} else {
						gearComponentTemplate.Execute(w, data)
					}
				})
				r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
					target := r.Header.Get("HX-Target")
					id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 0)
					var selectedGear data.Gear
					for _, gear := range data.Gears {
						if gear.ID == id {
							selectedGear = gear
							break
						}
					}
					data := map[string]interface{}{"Gear": selectedGear, "Locale": translator.GetLocale()}
					if target == "page" {
						gearEditTemplate.ExecuteTemplate(w, "page", data)
					} else {
						gearEditTemplate.Execute(w, data)
					}
				})
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
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		fmt.Println(r.URL.Query())
	})
	fmt.Println(http.ListenAndServe(":8080", r))
}

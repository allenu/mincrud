package main

import (
    "net/http"
    "html/template"
    "path"
	"google.golang.org/appengine"
    "strings"

    "github.com/allenu/minauth"
    "github.com/allenu/minauth/auth"
    _ "github.com/allenu/mincrud/api"
    "github.com/allenu/mincrud/core"
    "github.com/allenu/mincrud/core/model"
)

type ReadPageData struct {
    UserInfo auth.UserInfo
    EntityId string
    EditMode bool
    IsSignedIn bool
}

type EntitiesListPageData struct {
    UserInfo auth.UserInfo
    Entities []model.Entity
}

func init() {
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/edit/", editHandler)
    http.HandleFunc("/read/", readHandler)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    userInfo := minauth.GetUserInfo(r)
    adminMode := false

    const resultsPerPage int = 100
    const page int32 = 0

    ctx := appengine.NewContext(r)
    db := model.NewGoogleEntityDatabase(ctx)
    sc := core.NewEntityController(db)
    entities, err := sc.ReadEntities(userInfo, adminMode)

    if err == nil {
        fp := path.Join("templates", "index.html")
        if tmpl, err := template.ParseFiles(fp); err == nil {
            pageData := EntitiesListPageData {
                UserInfo: userInfo,
                Entities: entities,
            }
            tmpl.Execute(w, pageData)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    } else {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func readEditHandler(w http.ResponseWriter, r *http.Request, editMode bool) {
    userInfo := auth.GetUserInfo(r)
    var urlPath string
    if (editMode) {
        urlPath = "/edit/"
    } else {
        urlPath = "/read/"
    }
    entityId := strings.TrimPrefix(r.URL.Path, urlPath)

    fp := path.Join("templates", "readEdit.html")
    if tmpl, err := template.ParseFiles(fp); err == nil {
        if err == nil {
            pageData := ReadPageData{UserInfo: userInfo, EntityId: entityId, EditMode: editMode, IsSignedIn: userInfo.IsSignedIn}
            err = tmpl.Execute(w, pageData)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    } else {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func readHandler(w http.ResponseWriter, r *http.Request) {
    readEditHandler(w, r, false)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
    readEditHandler(w, r, true)
}


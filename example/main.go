package main

import (
    "net/http"
    "html/template"
    "path"

    "github.com/allenu/minauth"
    _ "github.com/allenu/mincrud/api"
)

func init() {
    http.HandleFunc("/", indexHandler)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    userInfo := minauth.GetUserInfo(r)

    fp := path.Join("templates", "index.html")
    if tmpl, err := template.ParseFiles(fp); err == nil {
        if err == nil {
            tmpl.Execute(w, userInfo)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    } else {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}


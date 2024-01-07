package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/dapings/examples/go-programing-tour-2023/chatroom/global"
	"github.com/dapings/examples/go-programing-tour-2023/chatroom/logic"
)

func homeHandleFunc(w http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles(filepath.Join(global.RootDir, "template/home.html"))
	if err != nil {
		_, _ = fmt.Fprint(w, "home template parse error")

		return
	}

	err = tpl.Execute(w, nil)
	if err != nil {
		_, _ = fmt.Fprint(w, "home template execute error")

		return
	}
}

func userListHandleFunc(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	userList := logic.Broadcaster.GetUserList()
	b, err := json.Marshal(userList)

	if err != nil {
		_, _ = fmt.Fprint(w, `[]`)
	} else {
		_, _ = fmt.Fprint(w, string(b))
	}
}

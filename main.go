package main

import (
	"Game-Application/repository/mongo"
	"Game-Application/service/user"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/user/register", userRegisterHandler)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return
	}
}

func userRegisterHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		_, err := fmt.Fprintf(w, `{"error":"method not allowed"}`)
		if err != nil {
			panic("error writing to testWriter")
		}
		return

	}
	data, cErr := io.ReadAll(req.Body)
	if cErr != nil {
		_, _ = fmt.Fprintf(w, `{"error":"reading body error"}`)
		return
	}
	var request user.RegisterRequest
	err := json.Unmarshal(data, &request)
	if err != nil {
		_, _ = fmt.Fprintf(w, `{"error":"unmarshal json error"}`)
		return
	}
	repo, Merr := mongo.New("mongodb://localhost:27017", "game")
	if Merr != nil {
		_, _ = fmt.Fprintf(w, `{"error":"mongodb connect error"}`)
		return
	}
	UserSvc := user.New(repo)
	_, RErr := UserSvc.Register(request)
	if RErr != nil {
		_, _ = fmt.Fprintf(w, `{"error":"register error"}`+RErr.Error())
		return
	}
	_, _ = fmt.Fprintf(w, `{"success":true}`)
	return
}

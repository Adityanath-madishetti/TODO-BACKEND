package controller

import (
	"net/http"

	"github.com/adityanath-madishetti/todo/backend/utils"
)


func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	utils.SendJSONError(w,http.StatusNotFound,"No page with this route")
}
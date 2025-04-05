package handlers

import (

	"net/http"
)



func (h *RouterStruct) CreateConversation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(200 << 20)
	if err != nil {
		http.Error(w, "can't parse form", http.StatusBadRequest)
		return
	}

	conversationsName := r.FormValue("conversation_name")

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file not found", http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = h.Service.Conversations.CreateConversation(r.Context(), conversationsName, header.Filename, file)
	if err != nil {
		http.Error(w, "failed to create", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
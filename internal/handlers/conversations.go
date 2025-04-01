package handlers

import "net/http"

func (router *RouterStruct) createConversationHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		respondWithError(w, "form parse error: "+err.Error(), err, http.StatusBadRequest)
		return
	}

	conversationName := r.FormValue("conversation_name")
	if conversationName == "" {
		respondWithError(w, "conversation_name required", err, http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		respondWithError(w, "file receive error: "+err.Error(), err, http.StatusBadRequest)
		return
	}

	if err := router.Service.Conversation.Create(r.Context(), file,
		conversationName, header.Filename); err != nil {
		respondWithError(w, "create failed: "+err.Error(), err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

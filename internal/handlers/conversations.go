package handlers

import "net/http"

func (router *RouterStruct) createConversationHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "form parse error: "+err.Error(), http.StatusBadRequest)
		return
	}

	conversationName := r.FormValue("conversation_name")
	if conversationName == "" {
		http.Error(w, "conversation_name required", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file receive error: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := router.Service.Conversation.CreateConversation(r.Context(), file,
		conversationName, header.Filename); err != nil {
		http.Error(w, "create failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

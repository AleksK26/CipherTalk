package api

import (
	"net/http"

	"github.com/AleksK26/WASA_AleksK_2024-25/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) commentMessage(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
	ctx reqcontext.RequestContext,
) {
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conversationID := ps.ByName("conversationId")
	// Verify user is a member of this conversation
	isMember, err := rt.db.IsUserInConversation(conversationID, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to check conversation membership")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !isMember {
		http.Error(w, "Forbidden: You are not a member of this conversation", http.StatusForbidden)
		return
	}

	commentID, err := generateNewID()
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to generate comment ID")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := rt.db.CommentMessage(commentID, ps.ByName("messageId"), userID); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) uncommentMessage(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
	ctx reqcontext.RequestContext,
) {
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conversationID := ps.ByName("conversationId")
	// Verify user is a member of this conversation
	isMember, err := rt.db.IsUserInConversation(conversationID, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to check conversation membership")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !isMember {
		http.Error(w, "Forbidden: You are not a member of this conversation", http.StatusForbidden)
		return
	}

	if err := rt.db.UncommentMessage(ps.ByName("messageId"), userID); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

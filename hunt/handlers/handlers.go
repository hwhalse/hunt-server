package handlers

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"hunt/collections"
	"hunt/db"
	"hunt/utils"
	"net/http"
)

type Uid struct {
	Uid string `json:"uid"`
}

func HandleTestPing(w http.ResponseWriter, r *http.Request, log zerolog.Logger) {
	_, err := w.Write([]byte("hi"))
	if err != nil {
		log.Error().Err(err).Msg("unable to write test ping")
	}
}

func HandleDelete(w http.ResponseWriter, r *http.Request, log zerolog.Logger) {
	var uid Uid
	err := json.NewDecoder(r.Body).Decode(&uid)
	if err != nil {
		log.Error().Err(err).Msg("Unable to decode request body")
	}
	filter := db.CreateFilter("uid", uid.Uid)
	err = collections.CommandNodeCollection.Delete(r.Context(), filter)
	if err != nil {
		utils.SendApiFailure(w, err.Error())
		log.Error().Err(err).Msg("unable to delete command")
		return
	}
	utils.SendApiSuccess(w, nil, "Deleted node")
}

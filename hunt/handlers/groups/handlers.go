package groups

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"hunt/collections"
	"hunt/structs"
	"hunt/utils"
	"net/http"
)

func handleGetAllGroups(w http.ResponseWriter, r *http.Request, log zerolog.Logger) {
	allUnits, err := collections.GroupCollection.FindAll(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("unable to get all groups")
		utils.SendApiFailure(w, err.Error())
		return
	}
	utils.SendApiSuccess(w, allUnits, "Units returned")
}

func handleAddGroup(w http.ResponseWriter, r *http.Request, log zerolog.Logger) {
	var unit structs.Unit
	err := json.NewDecoder(r.Body).Decode(&unit)
	if err != nil {
		log.Error().Err(err).Msg("Unable to decode request")
	}
	err = collections.GroupCollection.InsertUnit(r.Context(), unit)
	if err != nil {
		log.Error().Err(err).Msg("Unable to insert unit")
		utils.SendApiFailure(w, err.Error())
		return
	}
	utils.SendApiSuccess(w, unit, "Unit inserted")
}

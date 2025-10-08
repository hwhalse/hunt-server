package location

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"hunt/collections"
	"hunt/structs"
	"hunt/utils"
	"net/http"
)

func getByUid(w http.ResponseWriter, r *http.Request, log zerolog.Logger) {
	uid := r.PathValue("uid")
	location, err := collections.LocationCollection.FindById(r.Context(), uid)
	if err != nil {
		log.Error().Err(err).Msg("unable to get location by id")
		utils.SendApiFailure(w, err.Error())
		return
	}
	utils.SendApiSuccess(w, location, "Found location")
}

func upsertLocation(w http.ResponseWriter, r *http.Request, log zerolog.Logger) {
	var usr structs.HuntUser
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		log.Error().Err(err).Msg("Unable to decode json")
	}
	err = collections.LocationCollection.UpdateUser(r.Context(), usr)
	if err != nil {
		log.Error().Err(err).Msg("Unable to insert location")
		utils.SendApiFailure(w, err.Error())
		return
	}
	utils.SendApiSuccess(w, nil, "Successfully inserted location")
}

func getAll(w http.ResponseWriter, r *http.Request, log zerolog.Logger) {
	allLocations, err := collections.LocationCollection.FindAllActive(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("unable to get all locations")
		utils.SendApiFailure(w, err.Error())
		return
	}
	utils.SendApiSuccess(w, allLocations, "Got all locations")
}

func readMany(w http.ResponseWriter, r *http.Request, log zerolog.Logger) {
	readManyResponse := &ReadManyLocations{}
	err := json.NewDecoder(r.Body).Decode(&readManyResponse)
	if err != nil {
		log.Error().Err(err).Msg("Unable to decode json")
	}
	many, err := collections.LocationCollection.ReadManyLocations(r.Context(), readManyResponse.Uids, readManyResponse.PartitionKeys)
	if err != nil {
		log.Error().Err(err).Msg("Unable to read locations")
		utils.SendApiFailure(w, err.Error())
		return
	}
	utils.SendApiSuccess(w, many, "Successfully read many locations")
	w.Header().Set("Content-Type", "application/json")
}

type ReadManyLocations struct {
	Uids          []string `json:"uids"`
	PartitionKeys []string `json:"partitionKeys"`
}

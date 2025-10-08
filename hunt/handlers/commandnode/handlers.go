package commandnode

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"hunt/collections"
	"hunt/db"
	"hunt/socket"
	"hunt/structs"
	"hunt/utils"
	"io"
	"net/http"
)

func handleGetByUid(w http.ResponseWriter, r *http.Request, log zerolog.Logger) {
	uid := r.PathValue("uid")
	node, err := collections.CommandNodeCollection.FindById(r.Context(), uid)
	if err != nil {
		log.Error().Err(err).Msg("unable to get command node by uid")
		utils.SendApiFailure(w, err.Error())
		return
	}
	utils.SendApiSuccess(w, node, "Found command node")
}

func handleUpsert(w http.ResponseWriter, r *http.Request, manager *socket.Manager, log zerolog.Logger) {
	var node structs.CommandNode
	err := json.NewDecoder(r.Body).Decode(&node)
	if err != nil {
		log.Error().Err(err).Msg("Unable to read request body")
	}
	if node.Uid == "" {
		node.Uid = uuid.NewString()
	}
	filter := db.CreateFilter("uid", node.Uid)
	updatedNode, err := collections.CommandNodeCollection.FindOneAndUpdate(r.Context(), filter, node)
	if err != nil {
		log.Error().Err(err).Msg("Unable to insert node")
		utils.SendApiFailure(w, err.Error())
		return
	}
	utils.SendApiSuccess(w, updatedNode, "Updated node")
	nodeBytes, err := json.Marshal(node)
	if err != nil {
		log.Error().Err(err).Msg("Unable to marshal json")
	}
	err = manager.BroadcastMessage(socket.Event{
		Type:    6,
		Payload: string(nodeBytes),
	})
	if err != nil {
		log.Error().Err(err).Msg("unable to broadcast message")
	}
}

func getAllNodes(w http.ResponseWriter, r *http.Request, log zerolog.Logger) {
	log.Info().Msg("Getting all command nodes")
	allNodes, err := collections.CommandNodeCollection.FindAll(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Unable to find all nodes")
		utils.SendApiFailure(w, err.Error())
		return
	}
	utils.SendApiSuccess(w, allNodes, "Found all nodes")
}

type ReadManyNodes struct {
	Uids          []string `json:"uids"`
	PartitionKeys []string `json:"partitionKeys"`
}

func handleReadMany(w http.ResponseWriter, r *http.Request, log zerolog.Logger) {
	var readMany ReadManyNodes
	err := json.NewDecoder(r.Body).Decode(&readMany)
	if err != nil {
		log.Error().Err(err).Msg("Unable to read request body")
	}
	nodes, err := collections.CommandNodeCollection.ReadManyNodes(r.Context(), readMany.Uids, readMany.PartitionKeys)
	if err != nil {
		log.Error().Err(err).Msg("Unable to read request body")
		utils.SendApiFailure(w, err.Error())
		return
	}
	utils.SendApiSuccess(w, nodes, "Read nodes from database")
}

type Uid struct {
	Uid string `json:"uid"`
}

func handleDelete(w http.ResponseWriter, r *http.Request, manager *socket.Manager, log zerolog.Logger) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse delete request body")
	}
	uid := Uid{}
	err = json.Unmarshal(body, &uid)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse delete msg")
	}
	log.Info().Str("uid", uid.Uid).Msg("incoming delete req")
	filter := db.CreateFilter("uid", uid.Uid)
	err = collections.CommandNodeCollection.Delete(r.Context(), filter)
	if err != nil {
		log.Error().Err(err).Msg("unable to delete command node")
		utils.SendApiFailure(w, err.Error())
		return
	}
	utils.SendApiSuccess(w, nil, "Deleted node")
	err = manager.BroadcastMessage(socket.Event{
		Type:    socket.CommandNodeDelete,
		Payload: uid.Uid,
	})
	if err != nil {
		log.Error().Err(err).Msg("unable to broadcast msg")
	}
}

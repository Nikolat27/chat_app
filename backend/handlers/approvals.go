package handlers

import (
	"chat_app/utils"
	"errors"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func (handler *Handler) CreateApproval(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	inviteLink := chi.URLParam(r, "invite_link")
	if inviteLink == "" {
		utils.WriteError(w, http.StatusBadRequest, "urlParam", "invite-link is missing")
		return
	}

	groupId, err := handler.getIdByInviteLink(inviteLink)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "getIdByInviteLink", err.Error())
		return
	}

	var input struct {
		Reason string `json:"reason"`
	}

	if err := utils.ParseJSON(r.Body, 1_000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "parseJson", err.Error())
		return
	}

	filter := bson.M{
		"group_id":     groupId,
		"requester_id": payload.UserId,
	}

	projection := bson.M{
		"_id": 1,
	}

	if _, err := handler.Models.Approval.Get(filter, projection); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteError(w, http.StatusBadRequest, "getApproval", err.Error())
			return
		}
	} else {
		// If no error returned, it means approval exists → block duplicate
		utils.WriteError(w, http.StatusBadRequest, "getApproval", "You only can create one approval per group")
		return
	}

	filter = bson.M{
		"_id": groupId,
	}

	projection = bson.M{
		"owner_id": 1,
	}

	groupInstance, err := handler.Models.Group.Get(filter, projection)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "getGroup", err.Error())
		return
	}

	newApproval, err := handler.Models.Approval.Create(groupId, groupInstance.OwnerId, payload.UserId, input.Reason)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "createApproval", err.Error())
		return
	}

	resp := map[string]string{
		"approval_id": newApproval.InsertedID.(primitive.ObjectID).Hex(),
	}

	utils.WriteJSON(w, http.StatusCreated, resp)
}

func (handler *Handler) EditApprovalStatus(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	approvalId := chi.URLParam(r, "approval_id")
	if approvalId == "" {
		utils.WriteError(w, http.StatusBadRequest, "urlParam", "approval id is missing")
		return
	}

	approvalObjectId, errResp := utils.ToObjectId(approvalId)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	var input struct {
		Status string `json:"status" bson:"status"`
	}

	if err := utils.ParseJSON(r.Body, 1_000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "parseJson", err.Error())
		return
	}

	if input.Status != "pending" && input.Status != "rejected" && input.Status != "approved" {
		utils.WriteError(w, http.StatusBadRequest, "approvalStatus", "status is invalid. Must be either pending, rejected or approved")
		return
	}

	filter := bson.M{
		"_id":            approvalObjectId,
		"group_owner_id": payload.UserId,
	}

	updates := bson.M{
		"status":      input.Status,
		"reviewed_at": time.Now(),
	}

	if _, err := handler.Models.Approval.Update(filter, updates); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "deleteApproval", err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, "approval updated successfully")
}

// GetReceivedApprovals -> You are the group owner and other people have requested you to approve them
func (handler *Handler) GetReceivedApprovals(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	filter := bson.M{
		"group_owner_id": payload.UserId,
	}

	page, pageLimit, errResp := utils.ParsePageAndLimitQueryParams(r.URL)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	approvals, err := handler.Models.Approval.GetAll(filter, bson.M{}, page, pageLimit)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "getAllApprovals", err.Error())
		return
	}

	resp := map[string]any{
		"approvals": approvals,
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

// GetSentApprovals -> The approvals you have sent to other group owners
func (handler *Handler) GetSentApprovals(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	filter := bson.M{
		"requested_id": payload.UserId,
	}

	page, pageLimit, errResp := utils.ParsePageAndLimitQueryParams(r.URL)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	approvals, err := handler.Models.Approval.GetAll(filter, bson.M{}, page, pageLimit)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "getAllApprovals", err.Error())
		return
	}

	resp := map[string]any{
		"approvals": approvals,
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (handler *Handler) DeleteApproval(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	approvalId := chi.URLParam(r, "approval_id")
	if approvalId == "" {
		utils.WriteError(w, http.StatusBadRequest, "urlParam", "approval id is missing")
		return
	}

	approvalObjectId, errResp := utils.ToObjectId(approvalId)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	filter := bson.M{
		"_id":            approvalObjectId,
		"group_owner_id": payload.UserId,
	}

	if _, err := handler.Models.Approval.Delete(filter); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "deleteApproval", err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, "approval deleted successfully")
}

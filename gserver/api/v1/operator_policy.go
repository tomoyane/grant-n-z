package v1

import (
	"encoding/json"
	"net/http"

	"github.com/tomoyane/grant-n-z/gserver/api"
	"github.com/tomoyane/grant-n-z/gserver/common/property"
	"github.com/tomoyane/grant-n-z/gserver/entity"
	"github.com/tomoyane/grant-n-z/gserver/log"
	"github.com/tomoyane/grant-n-z/gserver/model"
	"github.com/tomoyane/grant-n-z/gserver/service"
)

var omhInstance OperatorPolicy

type OperatorPolicy interface {
	Api(w http.ResponseWriter, r *http.Request)

	get(w http.ResponseWriter, r *http.Request)

	post(w http.ResponseWriter, r *http.Request)

	put(w http.ResponseWriter, r *http.Request)

	delete(w http.ResponseWriter, r *http.Request)
}

type OperatorPolicyImpl struct {
	Request               api.Request
	OperatorPolicyService service.OperatorPolicyService
}

func GetOperatorPolicyInstance() OperatorPolicy {
	if omhInstance == nil {
		omhInstance = NewOperatorPolicy()
	}
	return omhInstance
}

func NewOperatorPolicy() OperatorPolicy {
	log.Logger.Info("New `OperatorPolicy` instance")
	log.Logger.Info("Inject `Request`, `operatorMemberRoleService` to `OperatorPolicy`")
	return OperatorPolicyImpl{
		Request:               api.GetRequestInstance(),
		OperatorPolicyService: service.NewOperatorPolicyServiceService(),
	}
}

func (rmrhi OperatorPolicyImpl) Api(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := rmrhi.Request.VerifyToken(w, r, property.AuthOperator)
	if err != nil {
		return
	}

	switch r.Method {
	case http.MethodGet:
		rmrhi.get(w, r)
	case http.MethodPost:
		rmrhi.post(w, r)
	case http.MethodPut:
		rmrhi.put(w, r)
	case http.MethodDelete:
		rmrhi.delete(w, r)
	default:
		err := model.MethodNotAllowed()
		model.Error(w, err.ToJson(), err.Code)
	}
}

func (rmrhi OperatorPolicyImpl) get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(entity.OperatorPolicyUserId.String())

	roleMemberEntities, err := rmrhi.OperatorPolicyService.Get(id)
	if err != nil {
		model.Error(w, err.ToJson(), err.Code)
		return
	}

	res, _ := json.Marshal(roleMemberEntities)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (rmrhi OperatorPolicyImpl) post(w http.ResponseWriter, r *http.Request) {
	var roleMemberEntity *entity.OperatorPolicy

	body, err := rmrhi.Request.Intercept(w, r)
	if err != nil {
		return
	}

	json.Unmarshal(body, &roleMemberEntity)
	if err := rmrhi.Request.ValidateBody(w, roleMemberEntity); err != nil {
		return
	}

	roleMemberEntity, err = rmrhi.OperatorPolicyService.Insert(roleMemberEntity)
	if err != nil {
		model.Error(w, err.ToJson(), err.Code)
		return
	}

	res, _ := json.Marshal(roleMemberEntity)
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (rmrhi OperatorPolicyImpl) put(w http.ResponseWriter, r *http.Request) {
}

func (rmrhi OperatorPolicyImpl) delete(w http.ResponseWriter, r *http.Request) {
}
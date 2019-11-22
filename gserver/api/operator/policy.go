package operator

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

var operatorPlhInstance OperatorPolicy

type OperatorPolicy interface {
	// Implement policy api
	Api(w http.ResponseWriter, r *http.Request)

	// Http GET method
	get(w http.ResponseWriter, r *http.Request)

	// Http POST method
	post(w http.ResponseWriter, r *http.Request, body []byte)

	// Http PUT method
	put(w http.ResponseWriter, r *http.Request)

	// Http DELETE method
	delete(w http.ResponseWriter, r *http.Request)
}

type PolicyImpl struct {
	Request       api.Request
	PolicyService service.PolicyService
}

func GetOperatorPolicyInstance() OperatorPolicy {
	if operatorPlhInstance == nil {
		operatorPlhInstance = NewOperatorPolicy()
	}
	return operatorPlhInstance
}

func NewOperatorPolicy() OperatorPolicy {
	log.Logger.Info("New `OperatorPolicy` instance")
	log.Logger.Info("Inject `Request`, `PolicyService` to `OperatorPolicy`")
	return PolicyImpl{
		Request:       api.GetRequestInstance(),
		PolicyService: service.GetPolicyServiceInstance(),
	}
}

func (ph PolicyImpl) Api(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ph.Request.Intercept(w, r, property.AuthOperator)
	if err != nil {
		return
	}

	switch r.Method {
	case http.MethodGet:
		ph.get(w, r)
	case http.MethodPost:
		ph.post(w, r, body)
	case http.MethodPut:
		ph.put(w, r)
	case http.MethodDelete:
		ph.delete(w, r)
	default:
		err := model.MethodNotAllowed()
		model.WriteError(w, err.ToJson(), err.Code)
	}
}

func (ph PolicyImpl) get(w http.ResponseWriter, r *http.Request) {
	roleMemberEntities, err := ph.PolicyService.GetPolicies()
	if err != nil {
		model.WriteError(w, err.ToJson(), err.Code)
		return
	}

	res, _ := json.Marshal(roleMemberEntities)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (ph PolicyImpl) post(w http.ResponseWriter, r *http.Request, body []byte) {
	var policyEntity *entity.Policy

	json.Unmarshal(body, &policyEntity)
	if err := ph.Request.ValidateBody(w, policyEntity); err != nil {
		return
	}

	policy, err := ph.PolicyService.InsertPolicy(policyEntity)
	if err != nil {
		model.WriteError(w, err.ToJson(), err.Code)
		return
	}

	res, _ := json.Marshal(policy)
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (ph PolicyImpl) put(w http.ResponseWriter, r *http.Request) {
}

func (ph PolicyImpl) delete(w http.ResponseWriter, r *http.Request) {
}

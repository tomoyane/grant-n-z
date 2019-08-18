package service

import (
	"github.com/tomoyane/grant-n-z/gserver/entity"
	"github.com/tomoyane/grant-n-z/gserver/model"
)

type PolicyService interface {
	Get(queryParam string) ([]*entity.Policy, *model.ErrorResponse)

	GetPolicies() ([]*entity.Policy, *model.ErrorResponse)

	GetPoliciesByRoleId(roleId int) ([]*entity.Policy, *model.ErrorResponse)

	InsertPolicy(policy *entity.Policy) (*entity.Policy, *model.ErrorResponse)

	EncryptData(data string) (*string, error)

	DecryptData(data string) (*string, error)
}
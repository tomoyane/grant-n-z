package service

import (
	"github.com/satori/go.uuid"

	"github.com/tomoyane/grant-n-z/gserver/common/ctx"
	"github.com/tomoyane/grant-n-z/gserver/common/driver"
	"github.com/tomoyane/grant-n-z/gserver/common/property"
	"github.com/tomoyane/grant-n-z/gserver/data"
	"github.com/tomoyane/grant-n-z/gserver/entity"
	"github.com/tomoyane/grant-n-z/gserver/log"
	"github.com/tomoyane/grant-n-z/gserver/model"
)

var gsInstance GroupService

type GroupService interface {
	// Get all groups
	GetGroups() ([]*entity.Group, *model.ErrorResBody)

	// Get group that has the user
	GetGroupOfUser() ([]*entity.Group, *model.ErrorResBody)

	// Insert group
	InsertGroup(group *entity.Group) (*entity.Group, *model.ErrorResBody)
}

// GroupService struct
type GroupServiceImpl struct {
	groupRepository      data.GroupRepository
	userGroupRepository  data.UserGroupRepository
	roleRepository       data.RoleRepository
	permissionRepository data.PermissionRepository
}

// Get Policy instance.
// If use singleton pattern, call this instance method
func GetGroupServiceInstance() GroupService {
	if gsInstance == nil {
		gsInstance = NewGroupService()
	}
	return gsInstance
}

// Constructor
func NewGroupService() GroupService {
	log.Logger.Info("New `GroupService` instance")
	log.Logger.Info("Inject `GroupRepository`, `UserGroupRepository`, `RoleRepository`, `PermissionRepository` to `GroupService`")
	return GroupServiceImpl{
		groupRepository:      data.GetGroupRepositoryInstance(driver.Db),
		userGroupRepository:  data.GetUserGroupRepositoryInstance(driver.Db),
		roleRepository:       data.GetRoleRepositoryInstance(driver.Db),
		permissionRepository: data.GetPermissionRepositoryInstance(driver.Db),
	}
}

func (gs GroupServiceImpl) GetGroups() ([]*entity.Group, *model.ErrorResBody) {
	return gs.groupRepository.FindAll()
}

func (gs GroupServiceImpl) GetGroupOfUser() ([]*entity.Group, *model.ErrorResBody) {
	if ctx.GetUserId().(int) == 0 {
		return nil, model.BadRequest("Required user id")
	}
	return gs.userGroupRepository.FindGroupsByUserId(ctx.GetUserId().(int))
}

func (gs GroupServiceImpl) InsertGroup(group *entity.Group) (*entity.Group, *model.ErrorResBody) {
	group.Uuid, _ = uuid.NewV4()
	role, err := gs.roleRepository.FindByName(property.Admin)
	if err != nil {
		log.Logger.Info("Failed to get role for insert groups process")
		return nil, model.InternalServerError()
	}

	permission, err := gs.permissionRepository.FindByName(property.Admin)
	if err != nil {
		log.Logger.Info("Failed to get permission for insert groups process")
		return nil, model.InternalServerError()
	}

	return gs.groupRepository.SaveWithRelationalData(*group, role.Id, permission.Id)
}

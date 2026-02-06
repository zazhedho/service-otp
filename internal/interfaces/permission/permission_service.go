package interfacepermission

import (
	domainpermission "service-sender/internal/domain/permission"
	"service-sender/internal/dto"
	"service-sender/pkg/filter"
)

type ServicePermissionInterface interface {
	Create(req dto.PermissionCreate) (domainpermission.Permission, error)
	GetByID(id string) (domainpermission.Permission, error)
	GetAll(params filter.BaseParams) ([]domainpermission.Permission, int64, error)
	GetByResource(resource string) ([]domainpermission.Permission, error)
	GetUserPermissions(userId string) ([]domainpermission.Permission, error)
	Update(id string, req dto.PermissionUpdate) (domainpermission.Permission, error)
	Delete(id string) error
}

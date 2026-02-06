package interfacemenu

import (
	domainmenu "service-sender/internal/domain/menu"
	"service-sender/internal/dto"
	"service-sender/pkg/filter"
)

type ServiceMenuInterface interface {
	Create(req dto.MenuCreate) (domainmenu.MenuItem, error)
	GetByID(id string) (domainmenu.MenuItem, error)
	GetAll(params filter.BaseParams) ([]domainmenu.MenuItem, int64, error)
	GetActiveMenus() ([]domainmenu.MenuItem, error)
	GetUserMenus(userId string) ([]domainmenu.MenuItem, error)
	Update(id string, req dto.MenuUpdate) (domainmenu.MenuItem, error)
	Delete(id string) error
}

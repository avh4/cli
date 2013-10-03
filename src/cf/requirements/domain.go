package requirements

import (
	"cf"
	"cf/api"
	"cf/net"
	"cf/terminal"
)

type DomainRequirement interface {
	Requirement
	GetDomain() cf.Domain
}

type DomainApiRequirement struct {
	name       string
	ui         terminal.UI
	domainRepo api.DomainRepository
	domain     cf.Domain
}

func NewDomainRequirement(name string, ui terminal.UI, domainRepo api.DomainRepository) (req *DomainApiRequirement) {
	req = new(DomainApiRequirement)
	req.name = name
	req.ui = ui
	req.domainRepo = domainRepo
	return
}

func (req *DomainApiRequirement) Execute() bool {
	var apiStatus net.ApiStatus
	req.domain, apiStatus = req.domainRepo.FindByName(req.name)

	if apiStatus.IsError() {
		req.ui.Failed(apiStatus.Message)
		return false
	}

	if apiStatus.IsNotFound() {
		req.ui.Failed("Domain not found")
		return false
	}

	return true
}

func (req *DomainApiRequirement) GetDomain() cf.Domain {
	return req.domain
}
package space

import (
	"cf/api"
	"cf/configuration"
	"cf/models"
	"cf/requirements"
	"cf/terminal"
	"github.com/codegangsta/cli"
)

type ListSpaces struct {
	ui        terminal.UI
	config    configuration.Reader
	spaceRepo api.SpaceRepository
}

func NewListSpaces(ui terminal.UI, config configuration.Reader, spaceRepo api.SpaceRepository) (cmd ListSpaces) {
	cmd.ui = ui
	cmd.config = config
	cmd.spaceRepo = spaceRepo
	return
}

func (cmd ListSpaces) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	reqs = []requirements.Requirement{
		reqFactory.NewLoginRequirement(),
		reqFactory.NewTargetedOrgRequirement(),
	}
	return
}

func (cmd ListSpaces) Run(c *cli.Context) {
	cmd.ui.Say("Getting spaces in org %s as %s...\n",
		terminal.EntityNameColor(cmd.config.OrganizationFields().Name),
		terminal.EntityNameColor(cmd.config.Username()))

	foundSpaces := false
	table := cmd.ui.Table([]string{"name"})
	apiResponse := cmd.spaceRepo.ListSpaces(func(space models.Space) bool {
		table.Print([][]string{{space.Name}})
		foundSpaces = true
		return true
	})

	if apiResponse.IsNotSuccessful() {
		cmd.ui.Failed("Failed fetching spaces.\n%s", apiResponse.Message)
		return
	}

	if !foundSpaces {
		cmd.ui.Say("No spaces found")
	}
}

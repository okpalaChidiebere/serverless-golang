package groups

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/udacity/serverless-golang/src/dataLayer/groupsAccess"
	"github.com/udacity/serverless-golang/src/models"
	"github.com/udacity/serverless-golang/src/requests"
)

/*Other developers might call this Service*/
type GroupAccess interface {
	GetAllGroups() []models.Group
	CreateGroup(c *requests.CreateGroupRequest) (models.Group, error)
}

//This 'groupAccess' businessLogic can only coomunitcate with the external service through the Port(groupsAccess.Repository)
type groupAccess struct {
	groupRepo groupsAccess.Repository
}

func NewGroupAccess(r groupsAccess.Repository) GroupAccess {
	return &groupAccess{r}
}

func (g *groupAccess) GetAllGroups() []models.Group {
	return g.groupRepo.GetAllGroups()
}

func (g *groupAccess) CreateGroup(createReq *requests.CreateGroupRequest) (models.Group, error) {
	id := uuid.Must(uuid.NewV4(), nil).String() //create a new id

	// Initialize group
	group := models.Group{
		Id:          id,
		Name:        createReq.Name,
		Description: createReq.Description,
		Timestamp:   time.Now().String(),
	}

	return g.groupRepo.CreateGroup(group)
}

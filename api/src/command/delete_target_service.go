package command

import "github.com/google/uuid"

type DeleteTargetServiceCommand struct {
	Id uuid.UUID
}

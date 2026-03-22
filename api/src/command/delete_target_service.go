package command

import "github.com/google/uuid"

type DeleteTargetServiceCommand struct {
	ID uuid.UUID
}

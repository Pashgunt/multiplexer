package command

import "github.com/google/uuid"

type GetTargetServiceCommand struct {
	ID uuid.UUID
}

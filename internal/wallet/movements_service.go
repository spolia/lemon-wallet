package wallet

import (
	"context"

	"github.com/spolia/lemon-wallet/internal/wallet/movement"
	"github.com/spolia/lemon-wallet/internal/wallet/user"
)

type MovementService struct {
	userRepo     user.Repository
	movementRepo movement.Repository
}

// New creates a MovementService implementation.
func NewMovementService(userRepo user.Repository, movRepo movement.Repository) *MovementService {
	return &MovementService{userRepo: userRepo, movementRepo: movRepo}
}

func (m *MovementService) Create(ctx context.Context, movement movement.Movement) (int64, error) {
	movementID, err := m.movementRepo.Save(ctx, movement)
	if err != nil {
		return 0, err
	}

	return movementID, nil
}

func (m *MovementService) Search(ctx context.Context, limit, offset uint64, movType, currencyName string, userID int64) ([]movement.Row, error) {
	movements, err := m.movementRepo.Search(ctx, limit, offset, movType, currencyName, userID)
	if err != nil {
		return []movement.Row{}, err
	}

	return movements, nil
}

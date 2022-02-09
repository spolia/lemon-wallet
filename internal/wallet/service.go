package wallet

import (
	"context"
	"strings"

	"github.com/spolia/lemon-wallet/internal/wallet/movement"
	"github.com/spolia/lemon-wallet/internal/wallet/user"
)

type Service struct {
	userRepo     user.Repository
	movementRepo movement.Repository
}

// New creates a Service implementation.
func New(userRepo user.Repository, movRepo movement.Repository) *Service {
	return &Service{userRepo: userRepo, movementRepo: movRepo}
}

// CreateUser saves a new user
func (s *Service) CreateUser(ctx context.Context, name, lastName, alias, email string) (int64, error) {
	userID, err := s.userRepo.Save(ctx, name, lastName, alias, email)
	if err != nil {
		println(err.Error())
		return 0, err
	}

	// every time that a new user is saved is necessary init movements
	err = s.movementRepo.InitSave(ctx, movement.Movement{
		Type:   "init",
		UserID: userID,
	})
	if err != nil {
		// delete the created user
		println("Delete", err.Error())
		return 0, s.userRepo.Delete(ctx, userID)
	}

	return userID, nil
}

// GetUser returns an user
func (s *Service) GetUser(ctx context.Context, id int64) (user.User, error) {
	userResult, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return user.User{}, err
	}

	// now we need to get the account extract with the latest movements per currency
	accountExtract, err := s.movementRepo.GetAccountExtract(ctx, userResult.ID)
	if err != nil {
		return user.User{}, err
	}

	userResult.WalletStatement = accountExtract

	return userResult, nil
}

// CreateMovement saves a movement
func (s *Service) CreateMovement(ctx context.Context, movement movement.Movement) (int64, error) {
	movement.CurrencyName = strings.ToUpper(movement.CurrencyName)
	movementID, err := s.movementRepo.Save(ctx, movement)
	if err != nil {
		return 0, err
	}

	return movementID, nil
}

// SearchMovement returns the user movements given certain filters
func (s *Service) SearchMovement(ctx context.Context, userID int64, limit, offset uint64, movType, currencyName string) ([]movement.Row, error) {
	movements, err := s.movementRepo.Search(ctx, userID, limit, offset, movType, strings.ToUpper(currencyName))
	if err != nil {
		return []movement.Row{}, err
	}

	return movements, nil
}

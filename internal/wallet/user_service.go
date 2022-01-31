package wallet

import (
	"context"

	"github.com/spolia/lemon-wallet/internal/wallet/movement"
	"github.com/spolia/lemon-wallet/internal/wallet/user"
)

type UserService struct {
	userRepo     user.Repository
	movementRepo movement.Repository
}

// New creates a service implementation.
func NewUserService(userRepo user.Repository, movRepo movement.Repository) *UserService {
	return &UserService{userRepo: userRepo, movementRepo: movRepo}
}

func (u *UserService) Create(ctx context.Context, name, lastName, alias, email string) error {
	userID, err := u.userRepo.Save(ctx, name, lastName, alias, email)
	if err != nil {
		return err
	}

	//every time that i save a new user i have to create the initial movement with 3 currencies in 0
	m := movement.Movement{
		Type:        movement.DepositMov,
		Amount:      0,
		TotalAmount: 0,
		UserID:      userID,
	}

	err = u.movementRepo.InitInsert(ctx, m)
	if err != nil {
		// delete the created user
		return u.userRepo.Delete(ctx, userID)
	}

	return nil
}

func (u *UserService) Get(ctx context.Context, id int64) (user.User, error) {
	userResult, err := u.userRepo.Get(ctx, id)
	if err != nil {
		return user.User{}, err
	}

	// now we need to get the account extract with the latest movements per currency
	accountExtract, err := u.movementRepo.GetAccountExtract(ctx, userResult.ID)
	if err != nil {
		return user.User{}, err
	}

	userResult.WalletStatement = accountExtract

	return userResult, nil
}

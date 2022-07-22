package userrepository

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-kivik/kivik/v3"
	usermodel "github.com/meneketehe/hehe/app/model/user"
)

type userAccountRepository struct {
	Couch *kivik.Client
}

func NewUserAccountRepository(couch *kivik.Client) usermodel.UserAccountRepository {
	return &userAccountRepository{
		Couch: couch,
	}
}

func (r *userAccountRepository) FindByID(ID string) (*usermodel.UserAccount, error) {
	db := r.Couch.DB(context.TODO(), "users")
	row := db.Get(context.TODO(), ID)

	var account usermodel.UserAccount
	if err := row.ScanDoc(&account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *userAccountRepository) FindByEmail(email string) (*usermodel.UserAccount, error) {
	db := r.Couch.DB(context.TODO(), "users")
	row, err := db.Find(
		context.TODO(),
		gin.H{
			"limit": 1,
			"selector": gin.H{
				"email": email,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	if !row.Next() {
		return nil, nil
	}

	var account usermodel.UserAccount
	if err := row.ScanDoc(&account); err != nil {
		return nil, err
	}
	if row.Err() != nil {
		return nil, row.Err()
	}

	return &account, nil
}

func (r *userAccountRepository) Create(account *usermodel.UserAccount) (*usermodel.UserAccount, error) {
	db := r.Couch.DB(context.TODO(), "users")
	rev, err := db.Put(context.TODO(), account.ID, account)
	if err != nil {
		return nil, err
	}
	account.Rev = rev

	return account, nil
}

func (r *userAccountRepository) Update(account *usermodel.UserAccount) (*usermodel.UserAccount, error) {
	db := r.Couch.DB(context.TODO(), "users")
	rev, err := db.Put(context.TODO(), account.ID, account)
	if err != nil {
		return nil, err
	}
	account.Rev = rev

	return account, nil
}

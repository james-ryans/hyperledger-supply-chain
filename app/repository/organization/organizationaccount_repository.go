package repository

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-kivik/kivik/v3"
	"github.com/meneketehe/hehe/app/model"
)

type organizationAccountRepository struct {
	Couch *kivik.Client
}

func NewOrganizationAccountRepository(couch *kivik.Client) model.OrganizationAccountRepository {
	return &organizationAccountRepository{
		Couch: couch,
	}
}

func (r *organizationAccountRepository) FindByID(ID string) (*model.OrganizationAccount, error) {
	db := r.Couch.DB(context.TODO(), "admins")
	row := db.Get(context.TODO(), ID)

	var account model.OrganizationAccount
	if err := row.ScanDoc(&account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *organizationAccountRepository) FindByEmail(email string) (*model.OrganizationAccount, error) {
	db := r.Couch.DB(context.TODO(), "admins")
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

	var account model.OrganizationAccount
	if err := row.ScanDoc(&account); err != nil {
		return nil, err
	}
	if row.Err() != nil {
		return nil, row.Err()
	}

	return &account, nil
}

func (r *organizationAccountRepository) Update(account *model.OrganizationAccount) (*model.OrganizationAccount, error) {
	db := r.Couch.DB(context.TODO(), "admins")
	rev, err := db.Put(context.TODO(), account.ID, account)
	if err != nil {
		return nil, err
	}
	account.Rev = rev

	return account, nil
}

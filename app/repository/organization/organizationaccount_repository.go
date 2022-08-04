package repository

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-kivik/kivik/v3"
	"github.com/meneketehe/hehe/app/model"
	"github.com/meneketehe/hehe/app/model/enum"
)

type organizationAccountRepository struct {
	Couch *kivik.Client
}

func NewOrganizationAccountRepository(couch *kivik.Client) model.OrganizationAccountRepository {
	return &organizationAccountRepository{
		Couch: couch,
	}
}

func (r *organizationAccountRepository) FindAllUser() ([]*model.OrganizationAccount, error) {
	db := r.Couch.DB(context.TODO(), "admins")
	row, err := db.Find(
		context.TODO(),
		gin.H{
			"selector": gin.H{
				"type": enum.OrgAccUser,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	accs := make([]*model.OrganizationAccount, 0)
	for row.Next() {
		var acc model.OrganizationAccount
		if err := row.ScanDoc(&acc); err != nil {
			return nil, err
		}
		accs = append(accs, &acc)
	}
	if row.Err() != nil {
		return nil, row.Err()
	}

	return accs, nil
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

func (r *organizationAccountRepository) Create(account *model.OrganizationAccount) (*model.OrganizationAccount, error) {
	db := r.Couch.DB(context.TODO(), "admins")
	rev, err := db.Put(context.TODO(), account.ID, account)
	if err != nil {
		return nil, err
	}
	account.Rev = rev

	return account, nil
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

func (r *organizationAccountRepository) Delete(account *model.OrganizationAccount) error {
	db := r.Couch.DB(context.TODO(), "admins")
	_, err := db.Delete(context.TODO(), account.ID, account.Rev)
	if err != nil {
		return err
	}

	return nil
}

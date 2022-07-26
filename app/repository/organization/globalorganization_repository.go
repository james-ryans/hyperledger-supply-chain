package repository

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-kivik/kivik/v3"
	"github.com/meneketehe/hehe/app/model"
)

type globalOrganizationRepository struct {
	Couch *kivik.Client
}

func NewGlobalOrganizationRepository(couch *kivik.Client) model.GlobalOrganizationRepository {
	return &globalOrganizationRepository{
		Couch: couch,
	}
}

func (r *globalOrganizationRepository) FindAll() ([]*model.GlobalOrganization, error) {
	db := r.Couch.DB(context.TODO(), "organizations")
	row, err := db.Find(
		context.TODO(),
		gin.H{
			"selector": gin.H{},
		},
	)
	if err != nil {
		return nil, err
	}

	orgs := make([]*model.GlobalOrganization, 0)
	for row.Next() {
		var org model.GlobalOrganization
		if err := row.ScanDoc(&org); err != nil {
			return nil, err
		}

		orgs = append(orgs, &org)
	}
	if row.Err() != nil {
		return nil, row.Err()
	}

	return orgs, nil
}

func (r *globalOrganizationRepository) FindByID(ID string) (*model.GlobalOrganization, error) {
	db := r.Couch.DB(context.TODO(), "organizations")
	row := db.Get(context.TODO(), ID)

	var org model.GlobalOrganization
	if err := row.ScanDoc(&org); err != nil {
		return nil, err
	}

	return &org, nil
}

func (r *globalOrganizationRepository) FindByCode(code string) (*model.GlobalOrganization, error) {
	db := r.Couch.DB(context.TODO(), "organizations")
	row, err := db.Find(
		context.TODO(),
		gin.H{
			"limit": 1,
			"selector": gin.H{
				"code": code,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	if !row.Next() {
		return nil, nil
	}

	var org model.GlobalOrganization
	if err := row.ScanDoc(&org); err != nil {
		return nil, err
	}
	if row.Err() != nil {
		return nil, row.Err()
	}

	return &org, nil
}

func (r *globalOrganizationRepository) FindByMSPID(MSPID string) (*model.GlobalOrganization, error) {
	db := r.Couch.DB(context.TODO(), "organizations")
	row, err := db.Find(
		context.TODO(),
		gin.H{
			"limit": 1,
			"selector": gin.H{
				"msp_id": MSPID,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	if !row.Next() {
		return nil, nil
	}

	var org model.GlobalOrganization
	if err := row.ScanDoc(&org); err != nil {
		return nil, err
	}
	if row.Err() != nil {
		return nil, row.Err()
	}

	return &org, nil
}

func (r *globalOrganizationRepository) FindByDomain(domain string) (*model.GlobalOrganization, error) {
	db := r.Couch.DB(context.TODO(), "organizations")
	row, err := db.Find(
		context.TODO(),
		gin.H{
			"limit": 1,
			"selector": gin.H{
				"domain": domain,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	if !row.Next() {
		return nil, nil
	}

	var org model.GlobalOrganization
	if err := row.ScanDoc(&org); err != nil {
		return nil, err
	}
	if row.Err() != nil {
		return nil, row.Err()
	}

	return &org, nil
}

func (r *globalOrganizationRepository) Create(org *model.GlobalOrganization) (*model.GlobalOrganization, error) {
	db := r.Couch.DB(context.TODO(), "organizations")
	rev, err := db.Put(context.TODO(), org.ID, org)
	if err != nil {
		return nil, err
	}
	org.Rev = rev

	return org, nil
}

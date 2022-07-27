package repository

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-kivik/kivik/v3"
	"github.com/meneketehe/hehe/app/model"
)

type globalChannelRepository struct {
	Couch *kivik.Client
}

func NewGlobalChannelRepository(couch *kivik.Client) model.GlobalChannelRepository {
	return &globalChannelRepository{
		Couch: couch,
	}
}

func (r *globalChannelRepository) FindAll() ([]*model.GlobalChannel, error) {
	db := r.Couch.DB(context.TODO(), "channels")
	row, err := db.Find(
		context.TODO(),
		gin.H{
			"selector": gin.H{},
		},
	)
	if err != nil {
		return nil, err
	}

	chs := make([]*model.GlobalChannel, 0)
	for row.Next() {
		var ch model.GlobalChannel
		if err := row.ScanDoc(&ch); err != nil {
			return nil, err
		}

		chs = append(chs, &ch)
	}
	if row.Err() != nil {
		return nil, row.Err()
	}

	return chs, nil
}

func (r *globalChannelRepository) FindByID(ID string) (*model.GlobalChannel, error) {
	db := r.Couch.DB(context.TODO(), "channels")
	row := db.Get(context.TODO(), ID)

	var ch model.GlobalChannel
	if err := row.ScanDoc(&ch); err != nil {
		return nil, err
	}

	return &ch, nil
}

func (r *globalChannelRepository) FindByName(name string) (*model.GlobalChannel, error) {
	db := r.Couch.DB(context.TODO(), "channels")
	row, err := db.Find(
		context.TODO(),
		gin.H{
			"limit": 1,
			"selector": gin.H{
				"name": name,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	if !row.Next() {
		return nil, nil
	}

	var ch model.GlobalChannel
	if err := row.ScanDoc(&ch); err != nil {
		return nil, err
	}
	if row.Err() != nil {
		return nil, row.Err()
	}

	return &ch, nil
}

func (r *globalChannelRepository) Create(ch *model.GlobalChannel) (*model.GlobalChannel, error) {
	db := r.Couch.DB(context.TODO(), "channels")
	rev, err := db.Put(context.TODO(), ch.ID, ch)
	if err != nil {
		return nil, err
	}
	ch.Rev = rev

	return ch, nil
}

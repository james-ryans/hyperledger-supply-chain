package model

type Producer struct {
	Organization
}

type ProducerService interface {
	GetAllProducers(channelID string) ([]*Producer, error)
	GetProducerByID(channelID, ID string) (*Producer, error)
	CreateProducer(channelID string, producer *Producer) (*Producer, error)
}

type ProducerRepository interface {
	FindAll(channelID string) ([]*Producer, error)
	FindByID(channelID, ID string) (*Producer, error)
	Create(channelID string, producer *Producer) error
}

package response

import "github.com/meneketehe/hehe/app/model"

type globalChannelsResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type globalChannelResponse struct {
	ID            string                      `json:"id"`
	Name          string                      `json:"name"`
	Suppliers     []*liteOrganizationResponse `json:"suppliers"`
	Producers     []*liteOrganizationResponse `json:"producers"`
	Manufacturers []*liteOrganizationResponse `json:"manufacturers"`
	Distributors  []*liteOrganizationResponse `json:"distributors"`
	Retailers     []*liteOrganizationResponse `json:"retailers"`
}

type liteOrganizationResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GlobalChannelsResponse(chs []*model.GlobalChannel) []*globalChannelsResponse {
	res := make([]*globalChannelsResponse, 0)
	for _, ch := range chs {
		res = append(res, &globalChannelsResponse{
			ID:   ch.ID,
			Name: ch.Name,
		})
	}

	return res
}

func GlobalChannelResponse(ch *model.GlobalChannel, orgs []*model.GlobalOrganization) *globalChannelResponse {
	return &globalChannelResponse{
		ID:            ch.ID,
		Name:          ch.Name,
		Suppliers:     mapOrgsFromIDs(ch.SuppliersID, orgs),
		Producers:     mapOrgsFromIDs(ch.ProducersID, orgs),
		Manufacturers: mapOrgsFromIDs(ch.ManufacturersID, orgs),
		Distributors:  mapOrgsFromIDs(ch.DistributorsID, orgs),
		Retailers:     mapOrgsFromIDs(ch.RetailersID, orgs),
	}
}

func mapOrgsFromIDs(IDs []string, orgs []*model.GlobalOrganization) []*liteOrganizationResponse {
	res := make([]*liteOrganizationResponse, 0)
	for _, ID := range IDs {
		for _, org := range orgs {
			if org.ID == ID {
				res = append(res, &liteOrganizationResponse{
					ID:   org.ID,
					Name: org.Name,
				})
				break
			}
		}
	}

	return res
}

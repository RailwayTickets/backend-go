package mongo

import (
	"time"

	"github.com/RailwayTickets/backend-go/entity"
	"gopkg.in/mgo.v2/bson"
)

type ticket struct{}

func (ticket) Search(params *entity.TicketSearchParams) ([]entity.Ticket, error) {
	var found []entity.Ticket
	query := bson.M{
		"owner": nil,
	}
	if params.From != "" {
		query["from"] = params.From
	}
	if params.To != "" {
		query["to"] = params.To
	}
	departure := time.Time(params.Date)
	if !departure.IsZero() {
		query["$and"] = []bson.M{
			{
				"departure": bson.M{
					"$gte": departure,
				},
			},
			{
				"departure": bson.M{
					"$lte": departure.Add(time.Hour * 24),
				},
			},
		}
	}
	err := tickets.Find(query).All(&found)
	return found, err
}

func (ticket) AllDirections() ([]string, error) {
	var directions []string
	err := tickets.Find(nil).Distinct("to", &directions)
	return directions, err
}

func (ticket) AllDepartures() ([]string, error) {
	var departures []string
	err := tickets.Find(nil).Distinct("from", &departures)
	return departures, err
}

func (ticket) Buy(login, id string) error {
	err := tickets.Update(
		bson.M{
			"_id":   bson.ObjectIdHex(id),
			"owner": nil,
		},
		bson.M{
			"$set": bson.M{
				"owner": login,
			},
		},
	)
	return err
}

func (ticket) ByID(id string) (*entity.Ticket, error) {
	t := new(entity.Ticket)
	err := tickets.Find(bson.M{
		"_id": bson.ObjectIdHex(id),
	}).One(t)
	return t, err
}

func (ticket) Return(login, id string) error {
	err := tickets.Update(
		bson.M{
			"_id":   bson.ObjectIdHex(id),
			"owner": login,
		},
		bson.M{
			"$set": bson.M{
				"owner": nil,
			},
		},
	)
	return err
}

func (ticket) ForUser(login string) ([]entity.Ticket, error) {
	var found []entity.Ticket
	query := bson.M{
		"owner": login,
	}
	err := tickets.Find(query).All(&found)
	return found, err
}

func (ticket) ValidReturnForUser(login string) ([]entity.Ticket, error) {
	var found []entity.Ticket
	query := bson.M{
		"owner": login,
		"departure": bson.M{
			"$gt": bson.Now(),
		},
	}
	err := tickets.Find(query).All(&found)
	return found, err
}

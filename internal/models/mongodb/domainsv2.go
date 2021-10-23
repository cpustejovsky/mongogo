package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DomainModelV2 struct {
	DB *mongo.Database
}

func (d *DomainModelV2) UpdateDelivered(name string) error {
	deliveredDomains := d.DB.Collection("delivered_domains")
	item := bson.D{{Key: "domain_name", Value: name}}

	_, err := deliveredDomains.InsertOne(context.TODO(), item)
	if err != nil {
		return err
	}
	return nil
}
func (d *DomainModelV2) UpdateBounced(name string) error {
	bouncedDomains := d.DB.Collection("bounced_domains")
	item := bson.D{{Key: "domain_name", Value: name}}
	_, err := bouncedDomains.InsertOne(context.TODO(), item)
	if err != nil {
		return err
	}
	return nil
}

func (d *DomainModelV2) CheckStatus(name string) (string, error) {
	//Create Collections
	deliveredDomains := d.DB.Collection("delivered_domains")
	bouncedDomains := d.DB.Collection("bounced_domains")
	fmt.Println("created collections")
	bouncedDomainsList := QueryDomains(bouncedDomains, "domain_name", name)
	deliveredDomainsList := QueryDomains(deliveredDomains, "domain_name", name)
	fmt.Println("len(bouncedDomainsList)", len(bouncedDomainsList))
	fmt.Println("len(deliveredDomainsList)", len(deliveredDomainsList))
	if len(bouncedDomainsList) >= 1 {
		return "not a catch-all", nil
	}
	if len(deliveredDomainsList) >= 1000 {
		return "catch-all", nil
	} else {
		return "unknown", nil
	}
}

func QueryDomains(collection *mongo.Collection, filter, name string) []bson.M {
	fmt.Println("querying collection %v", collection.Name())
	filterCursor, err := collection.Find(context.TODO(), bson.M{filter: name})
	if err != nil {
		log.Fatal(err)
	}
	var domain []bson.M
	if err = filterCursor.All(context.TODO(), &domain); err != nil {
		log.Fatal(err)
	}
	fmt.Println("successfully found domains")
	return domain
}

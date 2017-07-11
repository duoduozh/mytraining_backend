package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"models"
	"time"
)

type Training struct {
	Id            bson.ObjectId `bson:"_id,omitempty"`
	Name          string        `bson:"name"`
	Overview      string        `bson:"overview"`
	BasicInfo     string        `bson:"basic_info"`
	Commitment    string        `bson:"commitment"`
	HowToPass     string        `bson:"how_to_pass"`
	AverageRating int           `bson:"average_rating"`
	Icon          string        `bson:"icon"`
	Specification string        `bson:"specification_info"`
	Forum         string        `bson:"forum"`
	Language      []string      `bson:"language"`
}

type Lecture struct {
	Id           bson.ObjectId `bson:"_id,omitempty"`
	Name         string        `bson:"name"`
	Organization string        `bson:"organization"`
}

type Module struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	Week        int           `bson:"week"`
	Module      string        `bson:"module"`
	Title       string        `bson:"title"`
	Description string        `bson:"description"`
	Video       int           `bson:"video"`
	Reading     int           `bson:"reading"`
	Practice    int           `bson:"practice"`
	Duration    int           `bson:"duration"`
	Grade       int           `bson:"grade"`
}

var (
	IsDrop = true
)

func main() {
	session, err := mgo.Dial("127.0.0.1")

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("connect success\n")
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	// Drop DB
	if IsDrop {
		err = session.DB("test").DropDatabase()
		if err != nil {
			panic(err)
		}
	}

	// Collection People
	c := session.DB("test").C("people")

	// Index
	index := mgo.Index{
		Key:        []string{"name", "phone"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err = c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	person1 := models.Person{
		Name:      "Ale",
		Phone:     "+55 53 1234 4321",
		Timestamp: time.Now(),
		FAQList: []models.FAQ{
			models.FAQ{Question: "q1", Answer: "a1"},
			models.FAQ{Question: "q2", Answer: "a2"},
		},
	}
	// Insert Datas
	err = c.Insert(&person1,
		&models.Person{Name: "Cla", Phone: "+66 33 1234 5678", Timestamp: time.Now()})

	if err != nil {
		panic(err)
	}

	// Query One
	result := models.Person{}
	err = c.Find(bson.M{"name": "Ale"}).Select(bson.M{"phone": 0}).One(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println("Result Person", result)

	// Query All
	var resultPersonList []models.Person
	err = c.Find(bson.M{"name": "Ale"}).Sort("-timestamp").All(&resultPersonList)
	if err != nil {
		panic(err)
	}
	fmt.Println("Result Person List", resultPersonList)

	// Update
	colQuerier := bson.M{"name": "Ale"}
	change := bson.M{"$set": bson.M{"phone": "+86 99 8888 7776", "timestamp": time.Now()}}
	err = c.Update(colQuerier, change)
	if err != nil {
		panic(err)
	}

	// Query All
	err = c.Find(bson.M{"name": "Ale"}).Sort("-timestamp").All(&resultPersonList)

	if err != nil {
		panic(err)
	}
	fmt.Println("Results All: ", resultPersonList)

}

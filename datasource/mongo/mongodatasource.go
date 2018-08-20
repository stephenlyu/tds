package mongodatasource

import (
	"gopkg.in/mgo.v2"
	. "github.com/stephenlyu/tds/entity"
	. "github.com/stephenlyu/tds/period"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"github.com/stephenlyu/tds/datasource"
	"strings"
	"math"
)

type _MongoDataSource struct {
	session *mgo.Session
	dbName string
}

func NewMongoDataSource(dbUrl string, dbName string) datasource.BaseDataSource {
	session, err := mgo.Dial(dbUrl)
	if err != nil {
		panic(err)
	}
	return &_MongoDataSource{session: session, dbName: dbName}
}

func (this *_MongoDataSource) collectionName(security *Security, period Period) string {
	return strings.ToLower(fmt.Sprintf("%s_%s_%s", period.ShortName(), security.Code, security.Exchange))
}

func (this *_MongoDataSource) GetData(security *Security, period Period) (error, []Record) {
	colName := this.collectionName(security, period)
	l := []Record{}
	err := this.session.DB(this.dbName).C(colName).Find(bson.M{}).Sort("_id").All(&l)
	if err != nil {
		return err, nil
	}
	return nil, l
}

func (this *_MongoDataSource) GetDataEx(security *Security, period Period, startDate uint64, count int) (error, []Record) {
	colName := this.collectionName(security, period)
	l := []Record{}
	err := this.session.DB(this.dbName).C(colName).Find(bson.M{"_id": bson.M{"$gte": startDate}}).Sort("_id").Limit(count).All(&l)
	if err != nil {
		return err, nil
	}
	return nil, l
}

func (this *_MongoDataSource) GetRangeData(security *Security, period Period, startDate, endDate uint64) (error, []Record) {
	colName := this.collectionName(security, period)

	cond := bson.M{}
	if startDate != 0 {
		cond["$gte"] = startDate
	}
	if endDate != 0 {
		cond["$lte"] = endDate
	}
	query := bson.M{}

	if len(cond) > 0 {
		query["_id"] = cond
	}

	l := []Record{}
	err := this.session.DB(this.dbName).C(colName).Find(query).Sort("_id").All(&l)
	if err != nil {
		return err, nil
	}
	return nil, l
}

func (this *_MongoDataSource) GetDataFromLast(security *Security, period Period, endDate uint64, count int) (error, []Record) {
	if endDate == 0 {
		endDate = math.MaxUint64
	}

	colName := this.collectionName(security, period)

	cond := bson.M{}
	if endDate != 0 {
		cond["$lte"] = endDate
	}
	query := bson.M{}

	if len(cond) > 0 {
		query["_id"] = cond
	}

	l := []Record{}
	err := this.session.DB(this.dbName).C(colName).Find(query).Sort("-_id").Limit(count).All(&l)
	if err != nil {
		return err, nil
	}
	return nil, l
}

func (this *_MongoDataSource) GetLastRecord(security *Security, period Period) (error, *Record) {
	colName := this.collectionName(security, period)
	var record *Record
	err := this.session.DB(this.dbName).C(colName).Find(bson.M{}).Sort("-_id").Limit(1).One(&record)
	if err != nil {
		return err, nil
	}
	return nil, record
}

func (this *_MongoDataSource) AppendData(security *Security, period Period, data []Record) error {
	return this.SaveData(security, period, data)
}

func (this *_MongoDataSource) SaveData(security *Security, period Period, data []Record) error {
	colName := this.collectionName(security, period)
	for i := range data {
		_, err := this.session.DB(this.dbName).C(colName).UpsertId(data[i].Date, bson.M{"$set": &data[i]})
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *_MongoDataSource) RemoveData(security *Security, period Period, startDate, endDate uint64) error {
	colName := this.collectionName(security, period)
	cond := bson.M{}
	if startDate != 0 {
		cond["$gte"] = startDate
	}
	if endDate != 0 {
		cond["$lte"] = endDate
	}
	query := bson.M{}

	if len(cond) > 0 {
		query["_id"] = cond
	}

	_, err := this.session.DB(this.dbName).C(colName).RemoveAll(query)
	return err
}

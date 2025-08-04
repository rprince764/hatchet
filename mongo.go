/*
 * Copyright 2022-present Kuei-chun Chen. All rights reserved.
 * mongo.go
 */

package hatchet

import (
	"context"
	"log"
	"net/url"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MAX_DOC_SIZE = 16 * (1024 * 1024)
	BATCH_SIZE   = 1000
)

type MongoDB struct {
	db          *mongo.Database
	hatchetName string
	url         string
	verbose     bool

	clients []interface{}
	drivers []interface{}
	logs    []interface{}
}

func NewMongoDB(connstr string, hatchetName string) (*MongoDB, error) {
	var err error
	mongodb := &MongoDB{url: connstr, hatchetName: hatchetName}
	clientOptions := options.Client().ApplyURI(connstr)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return mongodb, err
	}
	u, err := url.Parse(connstr)
	if err != nil {
		return nil, err
	}
	dbName := u.Path[1:]
	if dbName == "" || dbName == "admin" {
		dbName = "logdb"
	}
	mongodb.db = client.Database(dbName)
	return mongodb, err
}

func (ptr *MongoDB) GetVerbose() bool {
	return ptr.verbose
}

func (ptr *MongoDB) SetVerbose(b bool) {
	ptr.verbose = b
}

func (ptr *MongoDB) Begin() error {
	var err error
	log.Println("creating hatchet", ptr.hatchetName)
	collName := ptr.hatchetName
	for _, keys := range []bson.D{
		{{Key: "component", Value: 1}},
		{{Key: "context", Value: 1}, {Key: "date", Value: 1}},
		{{Key: "severity", Value: 1}},
		{{Key: "op", Value: 1}, {Key: "ns", Value: 1}, {Key: "filter", Value: 1}},
	} {
		index := mongo.IndexModel{
			Keys:    keys,
			Options: options.Index().SetUnique(false),
		}
		_, err = ptr.db.Collection(collName).Indexes().CreateOne(context.TODO(), index)
		if err != nil {
			return err
		}
	}

	collName = ptr.hatchetName + "_clients"
	index := mongo.IndexModel{
		Keys:    bson.D{{Key: "context", Value: 1}, {Key: "ip", Value: 1}},
		Options: options.Index().SetUnique(false),
	}
	_, err = ptr.db.Collection(collName).Indexes().CreateOne(context.TODO(), index)
	if err != nil {
		return err
	}

	return err
}

func (ptr *MongoDB) Commit() error {
	if len(ptr.logs) > 0 {
		if _, err := ptr.db.Collection(ptr.hatchetName).InsertMany(context.TODO(), ptr.logs); err != nil {
			log.Println("Commit logs failed:", err)
			return err
		}
		ptr.logs = []interface{}{}
	}
	if len(ptr.clients) > 0 {
		if _, err := ptr.db.Collection(ptr.hatchetName+"_clients").InsertMany(context.TODO(), ptr.clients); err != nil {
			log.Println("Commit clients failed:", err)
			return err
		}
		ptr.clients = []interface{}{}
	}
	if len(ptr.drivers) > 0 {
		if _, err := ptr.db.Collection(ptr.hatchetName+"_drivers").InsertMany(context.TODO(), ptr.drivers); err != nil {
			log.Println("Commit drivers failed:", err)
			return err
		}
		ptr.drivers = []interface{}{}
	}
	return nil
}

func (ptr *MongoDB) Close() error {
	return ptr.db.Client().Disconnect(context.TODO())
}

// Drop drops all tables of a hatchet
func (ptr *MongoDB) Drop() error {
	if err := ptr.db.Collection(ptr.hatchetName + "_audit").Drop(context.TODO()); err != nil {
		log.Println("Drop audit failed:", err)
	}
	if err := ptr.db.Collection(ptr.hatchetName + "_clients").Drop(context.TODO()); err != nil {
		log.Println("Drop clients failed:", err)
	}
	if err := ptr.db.Collection(ptr.hatchetName + "_drivers").Drop(context.TODO()); err != nil {
		log.Println("Drop drivers failed:", err)
	}
	if err := ptr.db.Collection(ptr.hatchetName + "_ops").Drop(context.TODO()); err != nil {
		log.Println("Drop ops failed:", err)
	}
	if err := ptr.db.Collection(ptr.hatchetName).Drop(context.TODO()); err != nil {
		log.Println("Drop logs failed:", err)
	}
	if _, err := ptr.db.Collection("hatchet").DeleteOne(context.TODO(), bson.M{"name": ptr.hatchetName}); err != nil {
		log.Println("Drop hatchet failed:", err)
	}
	return nil
}

func (ptr *MongoDB) InsertLog(index int, end string, doc *Logv2Info, stat *OpStat) error {
	var err error
	data := bson.M{
		"_id":                index,
		"date":               end,
		"severity":           doc.Severity,
		"component":          doc.Component,
		"context":            doc.Context,
		"msg":                doc.Msg,
		"plan":               doc.Attributes.PlanSummary,
		"type":               BsonD2M(doc.Attr)["type"],
		"ns":                 doc.Attributes.NS,
		"message":            doc.Message,
		"op":                 stat.Op,
		"filter":             stat.QueryPattern,
		"_index":             stat.Index,
		"milli":              doc.Attributes.Milli,
		"reslen":             doc.Attributes.Reslen,
		"remoteOpWaitMillis": doc.Attributes.RemoteOpWaitMillis,
		"resolvedViews":      doc.Attributes.ResolvedViews,
		"authorization":      doc.Attributes.Authorization,
		"catalogCacheDatabaseLookupDurationMillis":   doc.Attributes.CatalogCacheDatabaseLookupDurationMillis,
		"catalogCacheCollectionLookupDurationMillis": doc.Attributes.CatalogCacheCollectionLookupDurationMillis,
		"catalogCacheIndexLookupDurationMillis":      doc.Attributes.CatalogCacheIndexLookupDurationMillis,
		"placementVersionRefreshMillis":              doc.Attributes.PlacementVersionRefreshMillis,
		"queryFramework":                             doc.Attributes.QueryFramework,
		"cpuNanos":                                   doc.Attributes.CPUNanos,
		"totalOplogSlotDurationMicros":               doc.Attributes.TotalOplogSlotDurationMicros,
		"queues":                                     doc.Attributes.Queues,
		"planCacheShapeHash":                         doc.Attributes.PlanCacheShapeHash,
		"workingMillis":                              doc.Attributes.WorkingMillis,
	}
	ptr.logs = append(ptr.logs, data)
	if len(ptr.logs) > BATCH_SIZE {
		collName := ptr.hatchetName
		_, err = ptr.db.Collection(collName).InsertMany(context.TODO(), ptr.logs)
		ptr.logs = []interface{}{}
	}
	return err
}

func (ptr *MongoDB) InsertClientConn(index int, doc *Logv2Info) error {
	var err error
	client := doc.Client
	data := bson.M{
		"_id": index, "ip": client.IP, "port": client.Port, "conns": client.Conns, "accepted": client.Accepted,
		"ended": client.Ended, "context": doc.Context}
	ptr.clients = append(ptr.clients, data)
	if len(ptr.clients) > BATCH_SIZE {
		collName := ptr.hatchetName + "_clients"
		_, err = ptr.db.Collection(collName).InsertMany(context.TODO(), ptr.clients)
		ptr.clients = []interface{}{}
	}
	return err
}

func (ptr *MongoDB) InsertDriver(index int, doc *Logv2Info) error {
	var err error
	client := doc.Client
	data := bson.M{
		"_id": index, "ip": client.IP, "driver": client.Driver, "version": client.Version}
	ptr.drivers = append(ptr.drivers, data)
	if len(ptr.drivers) > BATCH_SIZE {
		collName := ptr.hatchetName + "_drivers"
		_, err = ptr.db.Collection(collName).InsertMany(context.TODO(), ptr.drivers)
		ptr.drivers = []interface{}{}
	}
	return err
}

func (ptr *MongoDB) UpdateHatchetInfo(info HatchetInfo) error {
	filter := bson.M{"name": ptr.hatchetName}
	update := bson.M{"$set": bson.M{"version": info.Version, "module": info.Module, "arch": info.Arch, "os": info.OS, "start": info.Start, "end": info.End}}
	upsertOptions := options.Update().SetUpsert(true)
	_, err := ptr.db.Collection("hatchet").UpdateOne(context.TODO(), filter, update, upsertOptions)
	return err
}

func (ptr *MongoDB) CreateMetaData() error {
	var err error
	log.Printf("insert ops into %v_ops\n", ptr.hatchetName)
	pipeline := []bson.M{
		{"$match": bson.M{
			"op": bson.M{
				"$nin": []interface{}{nil, ""},
			},
		}},
		{"$group": bson.M{
			"_id": bson.M{
				"op":     "$op",
				"ns":     "$ns",
				"filter": "$filter",
				"_index": "$_index",
			},
			"count":    bson.M{"$sum": 1},
			"avg_ms":   bson.M{"$avg": "$milli"},
			"max_ms":   bson.M{"$max": "$milli"},
			"total_ms": bson.M{"$sum": "$milli"},
			"reslen":   bson.M{"$sum": "$reslen"},
		}},
		{"$project": bson.M{
			"_id":      0,
			"op":       "$_id.op",
			"count":    1,
			"avg_ms":   bson.M{"$round": []interface{}{"$avg_ms", 0}},
			"max_ms":   1,
			"total_ms": 1,
			"ns":       "$_id.ns",
			"_index":   "$_id._index",
			"reslen":   1,
			"filter":   "$_id.filter",
		}},
		{"$merge": bson.M{
			"into": ptr.hatchetName + "_ops",
		}},
	}
	if _, err = ptr.db.Collection(ptr.hatchetName).Aggregate(context.TODO(), pipeline); err != nil {
		log.Println("CreateMetaData ops failed:", err)
		return err
	}

	log.Printf("insert [exception] into %v_audit\n", ptr.hatchetName)
	pipeline = []bson.M{
		{"$match": bson.M{
			"severity": bson.M{
				"$in": []interface{}{"W", "E", "F"},
			},
		}},
		{"$group": bson.M{
			"_id": bson.M{
				"severity": "$severity",
			},
			"count": bson.M{"$sum": 1},
		}},
		{"$project": bson.M{
			"_id":   0,
			"type":  "exception",
			"name":  "$_id.severity",
			"value": "$count",
		}},
		{"$merge": bson.M{
			"into": ptr.hatchetName + "_audit",
		}},
	}
	if _, err = ptr.db.Collection(ptr.hatchetName).Aggregate(context.TODO(), pipeline); err != nil {
		log.Println("CreateMetaData exception failed:", err)
		return err
	}

	log.Printf("insert [failed] into %v_audit\n", ptr.hatchetName)
	pipeline = []bson.M{
		{"$match": bson.M{
			"message": bson.M{
				"$regex": `(\w\sfailed\s)`,
			},
		}},
		{"$group": bson.M{
			"_id": bson.M{
				"$substr": bson.A{"$message", 0, bson.M{
					"$add": bson.A{
						bson.M{
							"$indexOfBytes": bson.A{"$message", "failed"},
						},
						6,
					},
				}},
			},
			"count": bson.M{"$sum": 1},
		}},
		{"$project": bson.M{
			"_id":   0,
			"type":  "failed",
			"name":  "$_id",
			"value": "$count",
		}},
		{"$merge": bson.M{
			"into": ptr.hatchetName + "_audit",
		}},
	}
	if _, err = ptr.db.Collection(ptr.hatchetName).Aggregate(context.TODO(), pipeline); err != nil {
		log.Println("CreateMetaData failed failed:", err)
		return err
	}

	log.Printf("insert [op] into %v_audit\n", ptr.hatchetName)
	pipeline = []bson.M{
		{"$match": bson.M{
			"op": bson.M{
				"$nin": []interface{}{nil, ""},
			},
		}},
		{"$group": bson.M{
			"_id": bson.M{
				"op": "$op",
			},
			"count": bson.M{"$sum": 1},
		}},
		{"$project": bson.M{
			"_id":   0,
			"type":  "op",
			"name":  "$_id.op",
			"value": "$count",
		}},
		{"$merge": bson.M{
			"into": ptr.hatchetName + "_audit",
		}},
	}
	if _, err = ptr.db.Collection(ptr.hatchetName).Aggregate(context.TODO(), pipeline); err != nil {
		log.Println("CreateMetaData op failed:", err)
		return err
	}

	log.Printf("insert [ip] into %v_audit\n", ptr.hatchetName)
	pipeline = []bson.M{
		{"$group": bson.M{
			"_id": bson.M{
				"ip": "$ip",
			},
			"open": bson.M{"$sum": "$accepted"},
		}},
		{"$project": bson.M{
			"_id":   0,
			"type":  "ip",
			"name":  "$_id.ip",
			"value": "$open",
		}},
		{"$merge": bson.M{
			"into": ptr.hatchetName + "_audit",
		}},
	}
	if _, err = ptr.db.Collection(ptr.hatchetName+"_clients").Aggregate(context.TODO(), pipeline); err != nil {
		log.Println("CreateMetaData ip failed:", err)
		return err
	}

	log.Printf("insert [reslen-ip] into %v_audit\n", ptr.hatchetName)
	pipeline = []bson.M{
		{"$match": bson.M{
			"op": bson.M{
				"$nin": []interface{}{nil, ""},
			},
			"reslen": bson.M{"$gt": 0},
		}},
		{"$lookup": bson.M{
			"from": ptr.hatchetName + "_clients",
			"let":  bson.M{"context": "$context"},
			"pipeline": []bson.M{
				{"$match": bson.M{
					"$expr": bson.M{"$eq": []interface{}{"$context", "$$context"}},
				}},
				{"$project": bson.M{
					"_id": 0,
					"ip":  1,
				}},
			},
			"as": "clients",
		}},
		{"$unwind": bson.M{"path": "$clients"}},
		{"$group": bson.M{
			"_id":    "$clients.ip",
			"reslen": bson.M{"$sum": "$reslen"},
		}},
		{"$project": bson.M{
			"_id":   0,
			"type":  "reslen-ip",
			"name":  "$_id",
			"value": "$reslen",
		}},
		{"$merge": bson.M{
			"into": ptr.hatchetName + "_audit",
		}},
	}
	if _, err = ptr.db.Collection(ptr.hatchetName).Aggregate(context.TODO(), pipeline); err != nil {
		log.Println("CreateMetaData reslen-ip failed:", err)
		return err
	}

	log.Printf("insert [ns] into %v_audit\n", ptr.hatchetName)
	pipeline = []bson.M{
		{"$match": bson.M{
			"op": bson.M{
				"$nin": []interface{}{nil, ""},
			},
		}},
		{"$group": bson.M{
			"_id": bson.M{
				"ns": "$ns",
			},
			"count": bson.M{"$sum": 1},
		}},
		{"$project": bson.M{
			"_id":   0,
			"type":  "ns",
			"name":  "$_id.ns",
			"value": "$count",
		}},
		{"$merge": bson.M{
			"into": ptr.hatchetName + "_audit",
		}},
	}
	if _, err = ptr.db.Collection(ptr.hatchetName).Aggregate(context.TODO(), pipeline); err != nil {
		log.Println("CreateMetaData ns failed:", err)
		return err
	}

	log.Printf("insert [reslen-ns] into %v_audit\n", ptr.hatchetName)
	pipeline = []bson.M{
		{"$match": bson.M{
			"ns": bson.M{
				"$nin": []interface{}{nil, ""},
			},
			"reslen": bson.M{"$gt": 0},
		}},
		{"$group": bson.M{
			"_id": bson.M{
				"ns": "$ns",
			},
			"reslen": bson.M{"$sum": "$reslen"},
		}},
		{"$project": bson.M{
			"_id":   0,
			"type":  "reslen-ns",
			"name":  "$_id.ns",
			"value": "$reslen",
		}},
		{"$merge": bson.M{
			"into": ptr.hatchetName + "_audit",
		}},
	}
	if _, err = ptr.db.Collection(ptr.hatchetName).Aggregate(context.TODO(), pipeline); err != nil {
		log.Println("CreateMetaData reslen-ns failed:", err)
		return err
	}
	return nil
}

func (ptr *MongoDB) GetQueryFrameworkCounts(duration string) ([]NameValue, error) {
	collName := ptr.hatchetName
	pipeline := []bson.M{
		{"$match": bson.M{"queryFramework": bson.M{"$ne": nil}}},
		{"$group": bson.M{"_id": "$queryFramework", "count": bson.M{"$sum": 1}}},
		{"$project": bson.M{"_id": 0, "name": "$_id", "value": "$count"}},
	}
	cursor, err := ptr.db.Collection(collName).Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	var results []NameValue
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (ptr *MongoDB) InsertFailedMessages(m *FailedMessages) error {
	return nil
}

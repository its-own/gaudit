package mongo

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/its-own/gaudit/db"
	in "github.com/its-own/gaudit/in"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Mongo holds necessary fields and mongo Database session to connect
type Mongo struct {
	*mongo.Client
	Database *mongo.Database
	hook     in.Hook
}

var instance *Mongo

func InitMongo(cl *mongo.Client, database *mongo.Database, hook in.Hook) *Mongo {
	instance = &Mongo{
		Client:   cl,
		Database: database,
		hook:     hook,
	}
	return instance
}
func GetDbConnection() *Mongo {
	return instance
}

func (d *Mongo) Ping(ctx context.Context) error {
	return d.Client.Ping(ctx, readpref.Primary())
}

func (d *Mongo) Disconnect(ctx context.Context) error {
	return d.Client.Disconnect(ctx)
}

// EnsureIndices creates indices for collection col
func (d *Mongo) EnsureIndices(ctx context.Context, col string, index []db.Index) error {
	_db := d.Database
	var indexModels []mongo.IndexModel
	for _, ind := range index {
		keys := bson.D{}
		for _, k := range ind.Keys {
			keys = append(keys, bson.E{Key: k.Key, Value: k.Asc})
		}
		opts := options.Index()
		if ind.Unique != nil {
			opts.SetUnique(*ind.Unique)
		}
		if ind.Sparse != nil {
			opts.SetSparse(*ind.Sparse)
		}
		if ind.Name != "" {
			opts.SetName(ind.Name)
		}
		if ind.ExpireAfter != nil {
			opts.SetExpireAfterSeconds(int32(ind.ExpireAfter.Seconds()))
		}
		im := mongo.IndexModel{
			Keys:    keys,
			Options: opts,
		}
		indexModels = append(indexModels, im)
	}
	if _, err := _db.Collection(col).Indexes().CreateMany(ctx, indexModels); err != nil {
		return err
	}
	return nil
}

// DropIndices drops indices from collection col
func (d *Mongo) DropIndices(ctx context.Context, col string, index []db.Index) error {
	if _, err := d.Database.Collection(col).Indexes().DropAll(ctx); err != nil {
		return err
	}
	return nil
}

// Insert inserts doc into collection
func (d *Mongo) Insert(ctx context.Context, col string, doc interface{}) error {
	var (
		err    error
		insRes *mongo.InsertOneResult
	)
	d.hook.PreSave(ctx, doc, nil, col, "insert", "")
	if insRes, err = d.Database.Collection(col).InsertOne(ctx, doc); err != nil {
		return err
	}
	d.hook.PostSave(ctx, doc, nil, col, "insert", insRes.InsertedID.(primitive.ObjectID).Hex())
	return nil
}

func (d *Mongo) InsertMany(ctx context.Context, col string, docs []interface{}) error {
	if _, err := d.Database.Collection(col).InsertMany(ctx, docs); err != nil {
		return err
	}
	return nil
}

// FindOne finds a doc by query
func (d *Mongo) FindOne(ctx context.Context, col string, q interface{}, v interface{}, sort ...interface{}) error {
	findOneOpts := options.FindOne()
	if len(sort) > 0 {
		findOneOpts = findOneOpts.SetSort(sort[0])
	}

	err := d.Database.Collection(col).FindOne(ctx, q, findOneOpts).Decode(v)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return db.ErrNotFound
		}
		return err
	}
	return nil
}

// List finds list of docs that matches query with skip and limit
func (d *Mongo) List(ctx context.Context, col string, filter interface{}, skip, limit int64, v interface{}, sort ...interface{}) error {
	findOpts := options.Find().SetSkip(skip).SetLimit(limit)
	if len(sort) > 0 {
		findOpts = findOpts.SetSort(sort[0])
	}
	cursor, err := d.Database.Collection(col).Find(ctx, filter, findOpts)
	if err != nil {
		return err
	}
	if err := cursor.All(ctx, v); err != nil {
		return err
	}

	return nil
}

// Aggregate runs aggregation q on docs and store the result on v
func (d *Mongo) Aggregate(ctx context.Context, col string, q []interface{}, v interface{}) error {
	cursor, err := d.Database.Collection(col).Aggregate(ctx, q)
	if err != nil {
		return err
	}
	if err := cursor.All(ctx, v); err != nil {
		return err
	}
	return nil
}

func (d *Mongo) AggregateWithDiskUse(ctx context.Context, col string, q []interface{}, v interface{}) error {
	opt := options.Aggregate().SetAllowDiskUse(true)
	cursor, err := d.Database.Collection(col).Aggregate(ctx, q, opt)
	if err != nil {
		return err
	}
	if err := cursor.All(ctx, v); err != nil {
		return err
	}
	return nil
}

func (d *Mongo) Distinct(ctx context.Context, col, field string, q interface{}, v interface{}) error {
	interfaces, err := d.Database.Collection(col).Distinct(ctx, field, q)
	if err != nil {
		return err
	}
	data, err := json.Marshal(interfaces)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func (d *Mongo) PartialUpdateMany(ctx context.Context, col string, filter interface{}, data interface{}) error {
	_, err := d.Database.Collection(col).UpdateMany(ctx, filter, bson.M{"$set": data})
	if err != nil {
		return err
	}
	return nil
}

func (d *Mongo) PartialUpdateManyByQuery(ctx context.Context, col string, filter interface{}, query db.UnorderedDbQuery) error {
	_, err := d.Database.Collection(col).UpdateMany(ctx, filter, query)
	if err != nil {
		return err
	}
	return nil
}

func (d *Mongo) BulkUpdate(ctx context.Context, col string, models []mongo.WriteModel) error {
	_, err := d.Database.Collection(col).BulkWrite(ctx, models)
	return err
}

func (d *Mongo) DeleteMany(ctx context.Context, col string, filter interface{}) error {
	_, err := d.Database.Collection(col).DeleteMany(ctx, filter)
	return err
}

func (d *Mongo) Count(ctx context.Context, col string, q interface{}) (int64, error) {
	cnt, err := d.Database.Collection(col).CountDocuments(ctx, q)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, db.ErrNotFound
		}
		return 0, err
	}
	return cnt, nil
}

func (d *Mongo) Update(ctx context.Context, col string, filter interface{}, data interface{}) error {
	var (
		err  error
		opts = options.FindOneAndUpdate().SetReturnDocument(options.After)
		res  bson.M
	)
	update := bson.M{
		"$set": data,
	}
	d.hook.PreSave(ctx, data, filter, col, "update", "")
	if err = d.Database.Collection(col).FindOneAndUpdate(ctx, filter, update, opts).Decode(&res); err != nil {
		return err
	}
	if id, ok := res["_id"].(primitive.ObjectID); ok {
		d.hook.PostSave(ctx, data, filter, col, "update", id.Hex())
	}
	return nil
}

package source

import (
	"context"
	"encoding/json"
	"net/url"
	"path"
	"strings"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoSource struct {
	Key         string   `json:"key" mapstructure:"key"`
	DSN         string   `json:"dsn" mapstructure:"dsn"`
	Collections []string `json:"collections" mapstructure:"collections"`
	DB          string
	Client      *mongo.Client
}

var key2MongoSourceTable = make(map[string]*MongoSource)

func InitMongoSources(ctx context.Context, sources []*MongoSource) error {
	for _, source := range sources {
		dsn := source.DSN
		u, err := url.Parse(dsn)
		if err != nil {
			return errors.Wrap(err, "mongo dsn parse error")
		}

		dbName := strings.TrimPrefix(u.Path, "/")
		if dbName == "" {
			return errors.New("mongo db name is empty")
		}

		option := options.Client().ApplyURI(dsn)
		if err = option.Validate(); err != nil {
			return errors.Wrap(err, "mongo dsn validate error")
		}
		client, err := mongo.Connect(ctx, option)
		if err != nil {
			return errors.Wrap(err, "mongo connect error")
		}

		source.DB = dbName
		source.Client = client

		key2MongoSourceTable[source.Key] = source
	}
	return nil
}

func GetMongoSource(key string) *MongoSource {
	return key2MongoSourceTable[key]
}

var _ Source = (*MongoSource)(nil)

func (s *MongoSource) Export(ctx context.Context) ([]Data, error) {
	data := make([]Data, 0)
	db := s.Client.Database(s.DB)
	for _, collection := range s.Collections {
		cur, err := db.Collection(collection).Find(ctx, bson.M{})
		if err != nil {
			return nil, errors.Wrapf(err, "mongo find error, collection: %s", collection)
		}

		docs := make([]map[string]interface{}, 0)
		if err = cur.All(ctx, &docs); err != nil {
			return nil, errors.Wrapf(err, "mongo decode error, collection: %s", collection)
		}

		collData, err := json.Marshal(docs)
		if err != nil {
			return nil, errors.Wrapf(err, "mongo marshal error, collection: %s", collection)
		}

		data = append(data, Data{
			Key:     path.Join(s.DB, collection),
			Content: collData,
		})
	}
	return data, nil
}

func (s *MongoSource) Keys() []string {
	keys := make([]string, 0)

	for _, collection := range s.Collections {
		keys = append(keys, path.Join(s.DB, collection))
	}
	return keys
}

func (s *MongoSource) Import(ctx context.Context, data Data) error {
	db := s.Client.Database(s.DB)
	coll := db.Collection(data.Key)
	tempDocs := make([]map[string]interface{}, 0)
	err := json.Unmarshal(data.Content, &tempDocs)
	if err != nil {
		return errors.Wrapf(err, "mongo unmarshal error, collection: %s", data.Key)
	}

	if len(tempDocs) == 0 {
		return nil
	}

	docs := make([]interface{}, 0)
	for _, doc := range tempDocs {
		docs = append(docs, doc)
	}

	opts := options.Replace().SetUpsert(true)
	for _, doc := range docs {
		_, err = coll.ReplaceOne(ctx, doc, doc, opts)
		if err != nil {
			return errors.Wrapf(err, "mongo insert error, collection: %s, doc: %+v", data.Key, doc)
		}
	}
	return nil
}

package source

import (
	"context"
	"data-backup/component/target"
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
	Id          string   `json:"id" mapstructure:"id"`
	DSN         string   `json:"dsn" mapstructure:"dsn"`
	Collections []string `json:"collections" mapstructure:"collections"`
	DB          string
	Client      *mongo.Client
}

func RegisterMongoSources(ctx context.Context, sources []*MongoSource) error {
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

		registerSource(source.Id, source)
	}
	return nil
}

var _ Source = (*MongoSource)(nil)

func (s *MongoSource) Backup(ctx context.Context, target target.Target) error {
	db := s.Client.Database(s.DB)
	for _, collection := range s.Collections {
		cur, err := db.Collection(collection).Find(ctx, bson.M{})
		if err != nil {
			return errors.Wrapf(err, "mongo find error, collection: %s", collection)
		}

		docs := make([]map[string]interface{}, 0)
		if err = cur.All(ctx, &docs); err != nil {
			return errors.Wrapf(err, "mongo decode error, collection: %s", collection)
		}

		collData, err := json.Marshal(docs)
		if err != nil {
			return errors.Wrapf(err, "mongo marshal error, collection: %s", collection)
		}

		err = target.Import(ctx, path.Join(s.DB, collection), collData)
		if err != nil {
			return errors.Wrapf(err, "target import error, collection: %s", collection)
		}

	}
	return nil
}

func (s *MongoSource) Restore(ctx context.Context, target target.Target) error {
	db := s.Client.Database(s.DB)

	for _, collection := range s.Collections {
		exportData, err := target.Export(ctx, path.Join(s.DB, collection))
		if err != nil {
			return errors.Wrapf(err, "target export error, collection: %s", collection)
		}
		tempDocs := make([]map[string]interface{}, 0)
		err = json.Unmarshal(exportData, &tempDocs)
		if err != nil {
			return errors.Wrapf(err, "mongo unmarshal error, collection: %s", collection)
		}
		if len(tempDocs) == 0 {
			return nil
		}

		docs := make([]interface{}, 0)
		for _, doc := range tempDocs {
			docs = append(docs, doc)
		}

		coll := db.Collection(collection)
		opts := options.Replace().SetUpsert(true)
		for _, doc := range docs {
			_, err = coll.ReplaceOne(ctx, doc, doc, opts)
			if err != nil {
				return errors.Wrapf(err, "mongo insert error, collection: %s, doc: %+v", collection, doc)
			}
		}
	}

	return nil
}

//    Copyright 2021 Anderson Rodrigues do Livramento

//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at

//        http://www.apache.org/licenses/LICENSE-2.0

//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package model

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type MongoDatabase struct {
	Conn     *mongo.Client
	Database string
	URI      string
}

func NewMongoDatabase(dbURI string) *MongoDatabase {
	md := &MongoDatabase{}
	cs, err := connstring.Parse(dbURI)
	if err != nil {
		log.Fatal(err)
	}
	md.Database = cs.Database
	md.URI = dbURI
	return md
}

func (md *MongoDatabase) Connect() error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	var err error
	md.Conn, err = mongo.Connect(ctx, options.Client().ApplyURI(md.URI))
	if err != nil {
		log.Fatal("Connection Error:", err)
	}
	return err
}

func (md *MongoDatabase) Close() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	if err := md.Conn.Disconnect(ctx); err != nil {
		log.Fatal(err)
	}
}

func (md *MongoDatabase) getCollection(collectionName string) *mongo.Collection {
	return md.Conn.Database(md.Database).Collection(collectionName)
}

func (md *MongoDatabase) FindOneUpsert(collection string, filter, update bson.D) *mongo.SingleResult {
	updOptions := options.FindOneAndUpdate().SetUpsert(true)
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	coll := md.getCollection(collection)
	return coll.FindOneAndUpdate(ctx, filter, update, updOptions)
}

func (md *MongoDatabase) FindOne(collection string, filter bson.D) *mongo.SingleResult {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	coll := md.getCollection(collection)
	return coll.FindOne(ctx, filter)
}

// Interface IDataStorage
func (md *MongoDatabase) FindOneUpsertParameter(data Parameter) (Parameter, error) {
	filter := bson.D{
		{
			Key:   "_id",
			Value: data.ID,
		},
	}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "value",
					Value: data.Value,
				},
			},
		},
	}
	var result Parameter
	err := md.FindOneUpsert("parameters", filter, update).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = ErrNoRows
	}
	return result, err
}

func (md *MongoDatabase) FindOneUpsertStatusDescription(data StatusDescription) (StatusDescription, error) {
	filter := bson.D{
		{
			Key:   "_id",
			Value: data.ID,
		},
	}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "motivo",
					Value: data.Motivo,
				},
			},
		},
	}
	var result StatusDescription
	err := md.FindOneUpsert("situacao_motivos", filter, update).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = ErrNoRows
	}
	return result, err
}

func (md *MongoDatabase) FindOneStatusDescriptionById(ID int64) (StatusDescription, error) {
	filter := bson.D{
		{
			Key:   "_id",
			Value: ID,
		},
	}
	var result StatusDescription
	err := md.FindOne("situacao_motivos", filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = ErrNoRows
	}
	return result, err
}

func (md *MongoDatabase) FindOneUpsertRiskLevel(data RiskLevel) (RiskLevel, error) {
	filter := bson.D{
		{
			Key:   "_id",
			Value: data.ID,
		},
	}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "grau_risco",
					Value: data.GrauRisco,
				},
			},
		},
	}
	var result RiskLevel
	err := md.FindOneUpsert("graus_risco", filter, update).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = ErrNoRows
	}
	return result, err
}

func (md *MongoDatabase) FindOneRiskLevelById(ID string) (RiskLevel, error) {
	filter := bson.D{
		{
			Key:   "_id",
			Value: ID,
		},
	}
	var result RiskLevel
	err := md.FindOne("graus_risco", filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = ErrNoRows
	}
	return result, err
}

func (md *MongoDatabase) FindOneUpsertBaseCompany(data BaseCompany) (BaseCompany, error) {
	filter := bson.D{
		{
			Key:   "_id",
			Value: data.ID,
		},
	}
	update := bson.D{
		{
			Key:   "$set",
			Value: data,
		},
	}
	var result BaseCompany
	err := md.FindOneUpsert("base_empresas", filter, update).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = ErrNoRows
	}
	return result, err
}

func (md *MongoDatabase) FindOneBaseCompanyById(ID string) (BaseCompany, error) {
	filter := bson.D{
		{
			Key:   "_id",
			Value: ID,
		},
	}
	var result BaseCompany
	err := md.FindOne("base_empresas", filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = ErrNoRows
	}
	return result, err
}

func (md *MongoDatabase) FindOneUpsertCompany(data Company) (Company, error) {
	filter := bson.D{
		{
			Key:   "_id",
			Value: data.ID,
		},
	}
	update := bson.D{
		{
			Key:   "$set",
			Value: data,
		},
	}
	var result Company
	err := md.FindOneUpsert("empresas", filter, update).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = ErrNoRows
	}
	return result, err
}

func (md *MongoDatabase) FindOneCompanyById(ID string) (Company, error) {
	filter := bson.D{
		{
			Key:   "_id",
			Value: ID,
		},
	}
	var result Company
	err := md.FindOne("empresas", filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = ErrNoRows
	}
	return result, err
}

func (md *MongoDatabase) FindOneUpsertCity(data City) (City, error) {
	filter := bson.D{
		{
			Key:   "_id",
			Value: data.ID,
		},
	}
	update := bson.D{
		{
			Key:   "$set",
			Value: data,
		},
	}
	var result City
	err := md.FindOneUpsert("municipios", filter, update).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = ErrNoRows
	}
	return result, err
}

func (md *MongoDatabase) FindOneCityById(ID int64) (City, error) {
	filter := bson.D{
		{
			Key:   "_id",
			Value: ID,
		},
	}
	var result City
	err := md.FindOne("municipios", filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = ErrNoRows
	}
	return result, err
}

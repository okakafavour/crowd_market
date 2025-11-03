package services

import (
	"context"
	"errors"
	"time"

	"CROWD_MARKET/config"
	"CROWD_MARKET/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var productCollection *mongo.Collection

func InitProductService() {
	productCollection = config.DB.Collection("products")
}

// Helper: Convert string to ObjectID
func toObjectID(id string) (primitive.ObjectID, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, errors.New("invalid ID")
	}
	return objID, nil
}

// Helper: Filter by ID and optionally user ownership
func buildFilter(productID, userID string) (bson.M, error) {
	objID, err := toObjectID(productID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}

	if userID != "" {
		userObjID, err := toObjectID(userID)
		if err != nil {
			return nil, errors.New("invalid user ID")
		}
		filter["user_id"] = userObjID
	}

	return filter, nil
}

// ✅ Add a new product
func AddProduct(product model.Product) (model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := productCollection.InsertOne(ctx, product)
	if err != nil {
		return model.Product{}, err
	}

	product.ID = result.InsertedID.(primitive.ObjectID)
	return product, nil
}

// ✅ Get all products owned by a specific user
func GetAllProducts(userID string) ([]model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userObjID, err := toObjectID(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	filter := bson.M{"user_id": userObjID}

	cursor, err := productCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []model.Product
	for cursor.Next(ctx) {
		var p model.Product
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// ✅ Get a single product by its ID
func GetProductByID(id string) (*model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter, err := buildFilter(id, "")
	if err != nil {
		return nil, err
	}

	var product model.Product
	err = productCollection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &product, nil
}

// ✅ Update a product (ensures ownership)
func UpdateProductByUser(id string, userID string, fields map[string]interface{}) (*model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter, err := buildFilter(id, userID)
	if err != nil {
		return nil, err
	}

	update := bson.M{"$set": fields}
	result := productCollection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updated model.Product
	if err := result.Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("product not found or not owned by user")
		}
		return nil, err
	}

	return &updated, nil
}

// ✅ Delete a product (ensures ownership)
func DeleteProductByUser(id string, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter, err := buildFilter(id, userID)
	if err != nil {
		return err
	}

	result, err := productCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("product not found or not owned by user")
	}

	return nil
}

// ✅ Delete product for admin (no ownership check)
func DeleteProduct(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := toObjectID(id)
	if err != nil {
		return err
	}

	_, err = productCollection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

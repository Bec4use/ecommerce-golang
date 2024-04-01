package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Bec4use/ecommerce-golang/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct    = errors.New("Can't find the product")
	ErrCantDecodeProducts = errors.New("Can't find the products")
	ErrUserIdIsNotValid   = errors.New("User id is not valid")
	ErrCantUpdateUser     = errors.New("Can't add this product to the cart")
	ErrCantRemoveItemCart = errors.New("Can't remove this item from the cart")
	ErrCantGetItem        = errors.New("Was unable to get the item from the cart")
	ErrCantBuyCartItem    = errors.New("Can't update the purchase of the item")
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	searchfromdb, err := prodCollection.Find(ctx, bson.M{"_id": productID})
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}
	var productCart []models.Product
	err = searchfromdb.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}
	return nil

}

func RemoveCartItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	//update := bson.D{{Key: "$pull", Value: bson.M{"usercart": bson.M{"_id": productID}}}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrCantRemoveItemCart
	}

	return nil

}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	//Fetch the cart of the user , find the cart total and create an order with the cart items then remove the cart items
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var getcartitems models.User
	var ordercart models.Order

	ordercart.Order_ID = primitive.NewObjectID()
	ordercart.Ordered_At = time.Now()
	ordercart.Order_Cart = make([]models.ProductUser, 0)
	ordercart.Payment_Method.COD = true

	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

	currentresult, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	ctx.Done()
	if err != nil {
		log.Println(err)
	}

}

func InstantBuyer() {

}

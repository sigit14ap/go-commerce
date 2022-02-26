package repository

import (
	"context"

	"github.com/sigit14ap/go-commerce/internal/domain"
	"github.com/sigit14ap/go-commerce/internal/domain/dto"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CartsRepo struct {
	db *mongo.Collection
}

func (c *CartsRepo) FindAll(ctx context.Context) ([]domain.Cart, error) {
	cursor, err := c.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var cartArray []domain.Cart
	err = cursor.All(ctx, &cartArray)
	return cartArray, err
}

func (c *CartsRepo) FindByID(ctx context.Context, userID primitive.ObjectID) (domain.Cart, error) {
	result := c.db.FindOne(ctx, bson.M{"userID": userID})

	var cart domain.Cart
	err := result.Decode(&cart)

	return cart, err
}

func (c *CartsRepo) FindCartItems(ctx context.Context, userID primitive.ObjectID) ([]domain.CartItem, error) {
	result := c.db.FindOne(ctx, bson.M{"userID": userID}, options.FindOne().SetProjection(bson.M{"cartItems": 1}))

	var cart domain.Cart
	err := result.Decode(&cart)

	return cart.CartItems, err
}

func (c *CartsRepo) FindItem(ctx context.Context, userID primitive.ObjectID, productID primitive.ObjectID) (domain.CartItem, error) {

	opts := options.FindOne().SetProjection(bson.M{"cartItems.$": 1})
	result := c.db.FindOne(ctx, bson.M{"userID": userID, "cartItems": bson.M{"$elemMatch": bson.M{"productID": productID}}}, opts)

	var cart domain.Cart
	err := result.Decode(&cart)

	if len(cart.CartItems) == 0 {
		return domain.CartItem{}, err
	} else {
		return cart.CartItems[0], err
	}
}

func (c *CartsRepo) AddCartItem(ctx context.Context, cartItem domain.CartItem, userID primitive.ObjectID) (domain.CartItem, error) {

	findCart := c.db.FindOne(ctx, bson.M{"userID": userID})
	var cart domain.Cart
	err := findCart.Decode(&cart)

	if err != nil {
		tempCart := domain.Cart{
			ID:        primitive.NewObjectID(),
			UserID:    userID,
			CartItems: []domain.CartItem{},
		}

		cart, err = c.Create(ctx, tempCart)
	}

	opts := options.FindOne().SetProjection(bson.M{"cartItems.$": 1})
	result := c.db.FindOne(ctx, bson.M{"userID": userID, "cartItems": bson.M{"$elemMatch": bson.M{"productID": cartItem.ProductID}}}, opts)

	var cartData domain.Cart
	_ = result.Decode(&cartData)

	if len(cartData.CartItems) == 0 {
		_, err := c.db.UpdateOne(ctx, bson.M{"userID": userID}, bson.M{"$addToSet": bson.M{"cartItems": cartItem}})
		return cartItem, err
	} else {
		item := cartData.CartItems[0]

		quantity := item.Quantity + cartItem.Quantity

		updateOptions := bson.M{"$set": bson.M{"cartItems.$.quantity": quantity}}
		_, err := c.db.UpdateOne(ctx, bson.M{"userID": userID, "cartItems.productID": cartItem.ProductID}, updateOptions)
		return cartItem, err
	}
}

func (c *CartsRepo) UpdateCartItem(ctx context.Context, cartItem domain.CartItem, userID primitive.ObjectID) (domain.CartItem, error) {
	updateOptions := bson.M{"$set": bson.M{"cartItems.$.quantity": cartItem.Quantity}}
	_, err := c.db.UpdateOne(ctx, bson.M{"userID": userID, "cartItems.productID": cartItem.ProductID}, updateOptions)
	return cartItem, err
}

func (c *CartsRepo) DeleteCartItem(ctx context.Context, productID primitive.ObjectID, userID primitive.ObjectID) error {
	updateOptions := bson.M{"$pull": bson.M{"cartItems": bson.M{"productID": productID}}}
	_, err := c.db.UpdateOne(ctx, bson.M{"userID": userID}, updateOptions)
	return err
}

func (c *CartsRepo) ClearCart(ctx context.Context, userID primitive.ObjectID) error {
	emptyArray := make([]domain.CartItem, 0)
	_, err := c.db.UpdateOne(ctx, bson.M{"userID": userID}, bson.M{"$set": bson.M{"cartItems": emptyArray}})
	return err
}

func (c *CartsRepo) Create(ctx context.Context, cart domain.Cart) (domain.Cart, error) {
	cart.ID = primitive.NewObjectID()
	if cart.CartItems == nil {
		cart.CartItems = make([]domain.CartItem, 0)
	}

	_, err := c.db.InsertOne(ctx, cart)
	return cart, err
}

func (c *CartsRepo) Update(ctx context.Context, cartInput dto.UpdateCartInput, cartID primitive.ObjectID) (domain.Cart, error) {
	updateQuery := bson.M{}

	if cartInput.CartItems != nil {
		updateQuery["products"] = cartInput.CartItems
	}

	_, err := c.db.UpdateOne(ctx, bson.M{"_id": cartID}, bson.M{"$set": updateQuery})
	findResult := c.db.FindOne(ctx, bson.M{"_id": cartID})

	var cart domain.Cart
	err = findResult.Decode(&cart)

	return cart, err
}

func (c *CartsRepo) Delete(ctx context.Context, cartID primitive.ObjectID) error {
	_, err := c.db.DeleteOne(ctx, bson.M{"_id": cartID})
	return err
}

func NewCartsRepo(db *mongo.Database) *CartsRepo {
	collection := db.Collection(cartsCollection)
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"UserID": 1},
		Options: options.Index().SetExpireAfterSeconds(0),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("unable to create cart collection index, %v", err)
	}

	return &CartsRepo{
		db: collection,
	}
}

package user

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Login string `json:"login,omitempty" bson:"login,omitempty"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	BillingAddress string `json:"billingAddress,omitempty" bson:"billingAddress,omitempty"`
	ShippingAddress string `json:"shippingAddress,omitempty" bson:"shippingAddress,omitempty"`
	Phone string `json:"phone,omitempty" bson:"phone,omitempty"`
}

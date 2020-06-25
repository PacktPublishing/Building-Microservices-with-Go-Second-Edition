package main

import (
	"fmt"

	"github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/currency/proto_examples/example"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/genproto/googleapis/rpc/status"
)

func main() {
	// define a Product message
	product := example.Product{}
	product.Id = 3
	product.Name = "Latte"
	product.Description = "Milky coffee"
	product.Price = 2.25
	product.Sku = "abc-123"

	// define an Order messages including sub messages
	order := example.Order{}
	order.Product = &product
	order.Customer = &example.Order_Customer{
		Name:  "Nic",
		Phone: "+44 7122 121 121",
	}

	// define a message containing a list
	serverOrder := example.ServerOrders{}
	serverOrder.Server = &example.ServerOrders_Server{
		Name: "Gordon",
	}
	serverOrder.Order = []*example.Order{&order}

	// define a message containing a map
	productPrice := example.RestaurantOrders{}
	productPrice.RestaurantOrderCount = map[string]float64{
		"Chiswick": 121,
		"Richmond": 332,
	}
	productPrice.RestaurantServerOrders = map[string]*example.ServerOrders{
		"Chiswick": &serverOrder,
	}

	// enum values
	enumDiet := example.Diet_DAIRY_FREE

	// enum to string
	dietString := enumDiet.String()

	// enum value from string
	dietValue := example.Diet_value[dietString]

	// enum from string
	enumDiet = example.Diet(dietValue)

	// list of enums
	recipie := example.Recipie{}
	recipie.Suitability = []example.Diet{
		example.Diet_DAIRY_FREE,
		example.Diet_GLUTEN_FREE,
	}

	// messages containing external types
	em := example.ErrorMessage{
		Error: &status.Status{ //	"google.golang.org/genproto/googleapis/rpc/status"
			Code: 3,
		},
	}

	// handling timestamps using ptypes
	em.Time = ptypes.TimestampNow()

	// convert a protobuf into an any type
	// must be a message
	a, _ := ptypes.MarshalAny(&recipie)
	b, _ := ptypes.MarshalAny(&product)

	// messages containing annonymous types
	ingredients := example.Ingredients{}
	ingredients.Name = "Coffee"
	ingredients.NutricionalInfo = []*any.Any{
		a,
		b,
	}

	// use the main type from Any

	// loop over the collection and cast as necessary
	for _, v := range ingredients.GetNutricionalInfo() {
		fmt.Println(v.GetTypeUrl())
		switch v.GetTypeUrl() {
		//The default type URL for a given message type is type.googleapis.com/_packagename_._messagename_.
		case "type.googleapis.com/example.Product":
			prod := example.Product{}
			ptypes.UnmarshalAny(v, &prod)
			fmt.Println(prod)
		case "type.googleapis.com/example.Recipie":
			recipie := example.Recipie{}
			ptypes.UnmarshalAny(v, &recipie)
			fmt.Println(recipie)
		}
	}

	// or you can use Is to type check
	for _, v := range ingredients.GetNutricionalInfo() {
		switch {
		case ptypes.Is(v, &example.Product{}):
			fmt.Println("Is Product")
		case ptypes.Is(v, &example.Recipie{}):
			fmt.Println("Is Recipie")
		}
	}

	sor := example.SubscribeOrderResponse{}
	sor.Message = &example.SubscribeOrderResponse_Order{Order: &example.Order{}}

	// to determine the underlying message type
	switch sor.Message.(type) {
	case *example.SubscribeOrderResponse_Order:
		fmt.Println("Is Order")
	case *example.SubscribeOrderResponse_Error:
		fmt.Println("Is Order")
	}

	// or you can use if as both Order and Error will not be set at the same time
	if err := sor.GetError(); err != nil {
		// do something with error
	}

	if order := sor.GetOrder(); order != nil {
		// do something with order
	}
}

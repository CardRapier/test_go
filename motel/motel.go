package motel

import (
	"fmt"

	"github.com/cardrapier/hello-fiber/database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var motelCol mongo.Collection

func SetupRoutes(app *fiber.App) {
	motelCol = database.Collections.Motel
	motelGroup := app.Group("/api/v1/motel")
	motelGroup.Post("", createMotel)
	motelGroup.Get("", getMotels)
}

type Pagination struct {
	Page int8   `query:"page" default:"1"`
	Name string `query:"name"`
}

func createMotel(c *fiber.Ctx) error {
	motel := CreateMotel{}
	if err := c.BodyParser(&motel); err != nil {
		return err
	}
	// fmt.Println(motel.Rooms)
	_, err := motelCol.InsertOne(c.Context(), motel)

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.ErrInternalServerError.Code).SendString("Error")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "success",
	})
}

func getMotels(c *fiber.Ctx) error {
	p := new(Pagination)

	if err := c.QueryParser(p); err != nil {
		return err
	}
	fmt.Print(p)
	// total := getCount(c, p.Name)
	room_lookup := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "rooms"}, {Key: "localField", Value: "rooms"},
			{Key: "foreignField", Value: "_id"}, {Key: "as", Value: "rooms"},
		}},
	}
	// skip := bson.D{{Key: "$skip", Value: 0}}
	// limit := bson.D{{Key: "$limit", Value: 20}}
	cursor, err := motelCol.Aggregate(c.Context(), mongo.Pipeline{room_lookup})
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).SendString(err.Error())
	}

	var motels []Motel
	for cursor.Next(c.Context()) {
		var motel Motel
		cursor.Decode(&motel)
		motels = append(motels, motel)
	}

	// count := len(motels)

	return c.Status(fiber.StatusCreated).JSON(motels)
}

func getCount(c *fiber.Ctx, name string) int64 {
	filter := bson.D{{Key: "name", Value: bson.D{{Key: "$regex", Value: name}}}}
	count, err := motelCol.CountDocuments(c.Context(), filter)
	if err != nil {
		fmt.Println("Error fetching count")
		return 0
	}
	return count
}

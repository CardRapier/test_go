package motel

import (
	"fmt"

	"github.com/cardrapier/hello-fiber/database"
	"github.com/cardrapier/hello-fiber/models"
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
	Page uint16 `query:"page" default:"1"`
	Name string `query:"name"`
}

func (pg *Pagination) Init() *Pagination {
	pg.Page = 1
	return pg
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
	p := new(Pagination).Init()

	if err := c.QueryParser(p); err != nil {
		return err
	}
	fmt.Print(p)
	total := getCount(c, p.Name)
	pages := uint16(total / models.Limit)
	if pages == 0 {
		pages = 1
	}

	// Aggregations
	room_name := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "name", Value: bson.M{
				"$regex": p.Name,
			}},
		}},
	}
	skip_value := (p.Page - 1) * models.Limit
	skip := bson.D{{Key: "$skip", Value: skip_value}}
	limit := bson.D{{Key: "$limit", Value: models.Limit}}
	room_lookup := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "rooms"}, {Key: "localField", Value: "rooms"},
			{Key: "foreignField", Value: "_id"}, {Key: "as", Value: "rooms"},
		}},
	}
	cursor, err := motelCol.Aggregate(c.Context(), mongo.Pipeline{room_name, skip, limit, room_lookup})
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).SendString(err.Error())
	}

	var motels []Motel
	for cursor.Next(c.Context()) {
		var motel Motel
		cursor.Decode(&motel)
		motels = append(motels, motel)
	}

	count := uint16(len(motels))
	data := models.PaginationResult[Motel]{
		Total:   total,
		Count:   count,
		Page:    p.Page,
		Pages:   pages,
		Results: motels,
	}

	return c.Status(fiber.StatusCreated).JSON(data)
}

func getCount(c *fiber.Ctx, name string) uint16 {
	filter := bson.D{{Key: "name", Value: bson.D{{Key: "$regex", Value: name}}}}
	count, err := motelCol.CountDocuments(c.Context(), filter)
	if err != nil {
		fmt.Println("Error fetching count")
		return 0
	}
	return uint16(count)
}

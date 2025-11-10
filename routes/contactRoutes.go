package routes

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"Go-Server/database"
	"Go-Server/models"
)

func ContactRoutes(app *fiber.App) {
	collection := database.GetCollection("contacts")

	app.Post("/contacts", func(c *fiber.Ctx) error {
		var contact models.Contact
		if err := c.BodyParser(&contact); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "cannot parse JSON"})
		}
		contact.ID = primitive.NewObjectID()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := collection.InsertOne(ctx, contact)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to insert contact"})
		}
		return c.Status(201).JSON(contact)
	})

	app.Get("/contacts", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to fetch contacts"})
		}
		var contacts []models.Contact
		if err = cursor.All(ctx, &contacts); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "cursor error"})
		}
		return c.JSON(contacts)
	})

	app.Get("/contacts/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		objID, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var contact models.Contact
		err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&contact)
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "contact not found"})
		}
		return c.JSON(contact)
	})

	app.Put("/contacts/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		objID, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
		}

		var contact models.Contact
		if err := c.BodyParser(&contact); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		update := bson.M{"$set": contact}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to update"})
		}
		return c.JSON(fiber.Map{"message": "contact updated"})
	})

app.Delete("/contacts/:id", func(c *fiber.Ctx) error {
	idParam := c.Params("id")
	fmt.Println("Deleting contact with ID:", idParam) // 👈 add this
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete"})
	}
	if result.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "contact not found"})
	}

	return c.JSON(fiber.Map{"message": "contact deleted"})
})

}

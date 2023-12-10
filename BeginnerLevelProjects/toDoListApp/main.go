package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Task represents a task item
type Task struct {
	ID     string `json:"id" bson:"_id,omitempty"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

var collection *mongo.Collection

func main() {
	app := fiber.New()

	//MongoDB connection options
	clientOptions := options.Client().ApplyURI(getMongoURI())
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %s", err)
	}
	defer client.Disconnect(context.Background())

	//Ping MongoDB to check the connection status
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %s", err)
	}
	log.Println("Connected to MongoDB!")

	db := client.Database(("todo-app"))
	collection = db.Collection("tasks")

	//Routes
	setupRoutes(app)

	//Start the server
	log.Fatal(app.Listen(":8080"))
}

func setupRoutes(app *fiber.App) {
	app.Get("/tasks", GetTasks)
	app.Post("/tasks", CreateTask)
	app.Post("/tasks/:id", UpdateTask)
	app.Get("/tasks/:id", GetTaskByID)
	app.Delete("/tasks/:id", DeleteTask)
}

func GetTasks(c *fiber.Ctx) error {
	cursor, err := collection.Find(context.Background(), nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer cursor.Close(context.Background())
	var tasks []Task
	if err := cursor.All(context.Background(), &tasks); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(tasks)
}

func CreateTask(c *fiber.Ctx) error {
	var task Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	_, err := collection.InsertOne(context.Background(), task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(task)
}
func UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	var updateTask Task

	if err := c.BodyParser(&updateTask); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": updateTask},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(updateTask)
}

func GetTaskByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var task Task

	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&task)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString(err.Error())
	}
	return c.JSON(task)
}

func DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := collection.DeleteOne(context.Background(), bson.M{"id": id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func getMongoURI() string {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}
	return uri
}

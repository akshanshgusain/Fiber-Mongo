package repository

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"mongo/db"
	"mongo/models"
	"net/http"
)

type Repository struct {
	Instance db.MongoInstance
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/employees", r.GetEmployees)
	api.Get("/employee/:id", r.GetEmployee)
	api.Post("/employee/", r.CreateEmployee)
	api.Put("/employee/:id", r.UpdateEmployee)
	api.Delete("/employee/:id", r.DeleteEmployee)
}

func (r *Repository) GetEmployees(c *fiber.Ctx) error {
	query := bson.D{{}}

	cursor, err := r.Instance.Db.Collection("employees").Find(c.Context(), query)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": err.Error(),
			})
	}
	var employees []models.Employee = make([]models.Employee, 0)

	if err := cursor.All(c.Context(), &employees); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": err.Error(),
			})
	}

	return c.JSON(employees)
}

func (r *Repository) GetEmployee(c *fiber.Ctx) error {
	return nil
}

func (r *Repository) CreateEmployee(ctx *fiber.Ctx) error {
	collection := r.Instance.Db.Collection("employees")

	employee := models.Employee{}

	if err := ctx.BodyParser(&employee); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message": "bad input",
			})

	}

	insertionResult, err := collection.InsertOne(ctx.Context(), employee)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": err.Error(),
			})
	}

	// Query Construction
	filter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
	createdRecord := collection.FindOne(ctx.Context(), filter)

	createdEmployee := &models.Employee{}
	err = createdRecord.Decode(createdEmployee)

	if err != nil {

		return ctx.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": err.Error(),
			})
	}

	return ctx.Status(http.StatusCreated).JSON(createdEmployee)
}

func (r *Repository) UpdateEmployee(ctx *fiber.Ctx) error {
	collection := r.Instance.Db.Collection("employees")
	idParam := ctx.Params("id")

	employeeID, err := primitive.ObjectIDFromHex(idParam)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "can't convert the given ID",
		})
	}

	employee := new(models.Employee)

	if err := ctx.BodyParser(employee); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	// Building Query
	query := bson.D{{Key: "_id", Value: employeeID}}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "name", Value: employee.Name},
				{Key: "age", Value: employee.Age},
				{Key: "salary", Value: employee.Salary},
			},
		},
	}

	err = collection.FindOneAndUpdate(ctx.Context(), query, update).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ctx.Status(http.StatusBadRequest).JSON(
				&fiber.Map{
					"message": err.Error(),
				})
		}

		return ctx.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": err.Error(),
			})
	}

	employee.ID = idParam
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "invalid Id",
			})
	}

	return ctx.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "Success",
		})
}

func (r *Repository) DeleteEmployee(ctx *fiber.Ctx) error {
	employeeID, err := primitive.ObjectIDFromHex(ctx.Params("id"))

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "problem with ID",
		})
	}

	// Build query

	query := bson.D{{Key: "_id", Value: employeeID}}
	collection := r.Instance.Db.Collection("employees")
	result, err := collection.DeleteOne(ctx.Context(), &query)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	if result.DeletedCount < 1 {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(
		fiber.Map{
			"message": "record deleted",
		})
}

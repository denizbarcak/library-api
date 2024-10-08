package handlers

import (
	"context"
	"library-api/configs"
	"library-api/models"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Kullanıcı oluşturma
func CreateUser(c *fiber.Ctx) error {
    var user models.User
    if err := c.BodyParser(&user); err != nil {
        return c.Status(http.StatusBadRequest).SendString(err.Error())
    }
    user.ID = primitive.NewObjectID()
    user.Books = []primitive.ObjectID{}

    collection := configs.DB.Database("library").Collection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := collection.InsertOne(ctx, user)
    if err != nil {
        return c.Status(http.StatusInternalServerError).SendString(err.Error())
    }
    return c.Status(http.StatusCreated).JSON(user)
}

// Kullanıcıları listeleme
func GetUsers(c *fiber.Ctx) error {
    collection := configs.DB.Database("library").Collection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        return c.Status(http.StatusInternalServerError).SendString(err.Error())
    }

    var users []models.User
    if err = cursor.All(ctx, &users); err != nil {
        return c.Status(http.StatusInternalServerError).SendString(err.Error())
    }

    return c.Status(http.StatusOK).JSON(users)
}

// Kullanıcı silme
func DeleteUser(c *fiber.Ctx) error {
    userID := c.Params("id")
    objID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return c.Status(http.StatusBadRequest).SendString("Invalid user ID")
    }

    collection := configs.DB.Database("library").Collection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
    if err != nil || result.DeletedCount == 0 {
        return c.Status(http.StatusNotFound).SendString("User not found")
    }

    return c.SendString("User deleted successfully")
}

// Kullanıcıya kitap ekleme (kitap alma)
func AddBookToUser(c *fiber.Ctx) error {
    userID := c.Params("id")
    var bookID primitive.ObjectID
    if err := c.BodyParser(&bookID); err != nil {
        return c.Status(http.StatusBadRequest).SendString("Invalid book ID")
    }

    objID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return c.Status(http.StatusBadRequest).SendString("Invalid user ID")
    }

    collection := configs.DB.Database("library").Collection("users")
    bookCollection := configs.DB.Database("library").Collection("books")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var user models.User
    err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
    if err != nil {
        return c.Status(http.StatusNotFound).SendString("User not found")
    }

    if len(user.Books) >= 5 {
        return c.Status(http.StatusForbidden).SendString("User has already borrowed 5 books")
    }

    var book models.Book
    err = bookCollection.FindOne(ctx, bson.M{"_id": bookID}).Decode(&book)
    if err != nil {
        return c.Status(http.StatusNotFound).SendString("Book not found")
    }
    if book.Borrowed {
        return c.Status(http.StatusForbidden).SendString("This book has already been borrowed by another user")
    }

    user.Books = append(user.Books, bookID)
    _, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"books": user.Books}})
    if err != nil {
        return c.Status(http.StatusInternalServerError).SendString("Failed to add book to user")
    }

    _, err = bookCollection.UpdateOne(ctx, bson.M{"_id": bookID}, bson.M{"$set": bson.M{"borrowed": true}})
    if err != nil {
        return c.Status(http.StatusInternalServerError).SendString("Failed to mark book as borrowed")
    }

    return c.SendString("Book added to user successfully")
}
// Kullanıcı girişi (Login)
func Login(c *fiber.Ctx) error {
    var input struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := c.BodyParser(&input); err != nil {
        return c.Status(http.StatusBadRequest).SendString(err.Error())
    }

    collection := configs.DB.Database("library").Collection("users")
    var user models.User
    err := collection.FindOne(c.Context(), bson.M{"email": input.Email, "password": input.Password}).Decode(&user)
    if err != nil {
        return c.Status(http.StatusUnauthorized).SendString("Invalid email or password")
    }

    // Kullanıcı doğrulandıysa başarılı yanıt döner
    return c.SendString("Login successful")
}


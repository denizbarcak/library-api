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


// Kitap oluşturma
func CreateBook(c *fiber.Ctx) error {
    var book models.Book
    if err := c.BodyParser(&book); err != nil { // c.BodyParser(&book) fonksiyonu, gelen HTTP isteği gövdesini okuyarak book değişkenine atar.
        return c.Status(http.StatusBadRequest).SendString(err.Error())
    }
    book.ID = primitive.NewObjectID()
    book.Borrowed = false

    collection := configs.DB.Database("library").Collection("books")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()  // c.BodyParser(&book) fonksiyonu, gelen HTTP isteği gövdesini okuyarak book değişkenine atar.

    _, err := collection.InsertOne(ctx, book)
    if err != nil {
        return c.Status(http.StatusInternalServerError).SendString(err.Error())
    }
    return c.Status(http.StatusCreated).JSON(book)
}

// Kitapları listeleme
func GetBooks(c *fiber.Ctx) error {
    collection := configs.DB.Database("library").Collection("books")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        return c.Status(http.StatusInternalServerError).SendString(err.Error())
    }

    var books []models.Book
    if err = cursor.All(ctx, &books); err != nil {
        return c.Status(http.StatusInternalServerError).SendString(err.Error())
    }

    return c.Status(http.StatusOK).JSON(books)
}

// Kitap silme
func DeleteBook(c *fiber.Ctx) error {
    bookID := c.Params("id")
    objID, err := primitive.ObjectIDFromHex(bookID)
    if err != nil {
        return c.Status(http.StatusBadRequest).SendString("Invalid book ID")
    }

    collection := configs.DB.Database("library").Collection("books")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
    if err != nil || result.DeletedCount == 0 {
        return c.Status(http.StatusNotFound).SendString("Book not found")
    }

    return c.SendString("Book deleted successfully")
}

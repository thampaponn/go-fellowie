package controller

import (
	"net/http"
	"strings"

	"fmt"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/thampaponn/go-fellowie/models"
	"github.com/thampaponn/go-fellowie/repository"
)

func CreateUser(c echo.Context) error {
	db, err := repository.InitDB()
	if err != nil {
		return err
	}
	var user models.User

	// Bind the incoming JSON to the User struct
	if err := c.Bind(&user); err != nil {
		fmt.Println("Error binding user:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if user.Name == "" || user.Email == "" || user.Password == "" {
		fmt.Println("All fields are required")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "All fields are required"})
	}
	// Insert the user into the database
	err = db.QueryRow(
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		user.Name, user.Email, user.Password,
	).Scan(&user.ID)
	if err != nil {
		fmt.Println("Error inserting user:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}
	return c.JSON(http.StatusCreated, user)
}

func GetUsers(c echo.Context) error {
	db, err := repository.InitDB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to connect to database"})
	}

	var users []models.User
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		fmt.Println("Error querying users:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve users"})
	}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			fmt.Println("Error scanning user:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to scan user"})
		}
		users = append(users, user)
	}

	m := make(map[string]int)
	m["total"] = len(users)
	m["number-1"] = 1

	fmt.Println("Total users:", m["total"])
	fmt.Println("Number 1:", m["number-1"])

	return c.JSON(http.StatusOK, users)
}

func UpdateUser(c echo.Context) error {
	db, err := repository.InitDB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to connect to database"})
	}

	id := c.Param("id")
	var user models.User

	if err := c.Bind(&user); err != nil {
		fmt.Println("Error binding user:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	setClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if user.Name != "" {
		//Add param to the query and +1 to argPos to ensure correct parameter indexing
		// Use fmt.Sprintf to format the string with the current argPos
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argPos))
		args = append(args, user.Name)
		argPos++
	}
	if user.Email != "" {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", argPos))
		args = append(args, user.Email)
		argPos++
	}
	if user.Password != "" {
		setClauses = append(setClauses, fmt.Sprintf("password = $%d", argPos))
		args = append(args, user.Password)
		argPos++
	}

	if len(setClauses) == 0 {
		return nil // nothing to update
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argPos)

	_, err = db.Exec(query, args...)
	if err != nil {
		fmt.Println("Error executing update query:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "User updated successfully"})
}

func DeleteUser(c echo.Context) error {
	db, err := repository.InitDB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to connect to database"})
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User ID is required"})
	}

	_, err = db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		fmt.Println("Error deleting user:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete user"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func saveEmail(email string) error {
	if emailExists(email) {
        return fmt.Errorf("Email already exists")
    }
	file, err := os.OpenFile("subscribers.txt",os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
        return err
    }

	defer file.Close()
	_, err = file.WriteString(email + "\n")
    return err
}

func emailExists(email string) bool {
	file, err := os.Open("subscribers.txt")
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == email {
			return true
		}
	}
	return false
}

func subscribe(c echo.Context) error {
	email := c.FormValue("email")
	if email == "" {
		return echo.NewHTTPError(400,"Email is required")
	}

	if err := saveEmail(email); err != nil {
		if err.Error() == "Email already exists" {
			return echo.NewHTTPError(409, "Email already subscribed")
		}
		return echo.NewHTTPError(500, "Failed to save email")
	}
	return c.String(http.StatusOK, "Subscribed successfully")
}

func main() {
	app := echo.New()
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	app.POST("/subscribe",subscribe)
	app.Logger.Fatal(app.Start(":8080"))
}
package handlers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/emarifer/gofiber-htmx-todolist/models"
	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/flash"
	"golang.org/x/crypto/bcrypt"
)

/********** Handlers for Auth Views **********/

// Render Home Page
func HandleViewHome(c *fiber.Ctx) error {

	return c.Render("home", fiber.Map{
		"FromProtected": fromProtected,
	})
}

// Render Login Page with success/error messages & session management
func HandleViewLogin(c *fiber.Ctx) error {

	if c.Method() == "POST" {
		var (
			user models.User
			err  error
		)
		fm := fiber.Map{
			"type": "alert-error",
		}

		// notice: in production you should not inform the user
		// with detailed messages about login failures
		if user, err = models.CheckEmail(c.FormValue("email")); err != nil {
			fm["message"] = "There is no user with that email"

			return flash.WithError(c, fm).Redirect("/login")
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(c.FormValue("password")),
		)
		if err != nil {
			fm["message"] = "Incorrect password"

			return flash.WithError(c, fm).Redirect("/login")
		}

		session, err := store.Get(c)
		if err != nil {
			fm["message"] = fmt.Sprintf("something went wrong: %s", err)

			return flash.WithError(c, fm).Redirect("/login")
		}

		session.Set(AUTH_KEY, true)
		session.Set(USER_ID, user.ID)

		err = session.Save()
		if err != nil {
			fm["message"] = fmt.Sprintf("something went wrong: %s", err)

			return flash.WithError(c, fm).Redirect("/login")
		}

		fm = fiber.Map{
			"type":    "alert-success",
			"message": "You have successfully logged in!!",
		}

		return flash.WithSuccess(c, fm).Redirect("/todo/list")
	}

	return c.Render("login", fiber.Map{
		"Page":          "login",
		"FromProtected": fromProtected,
		"Message":       flash.Get(c),
	})
}

// Render Register Page with success/error messages
func HandleViewRegister(c *fiber.Ctx) error {

	if c.Method() == "POST" {
		user := models.User{
			Email:    c.FormValue("email"),
			Password: c.FormValue("password"),
			Username: c.FormValue("username"),
		}

		err := models.CreateUser(user)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				err = errors.New("the email is already in use")
			}
			fm := fiber.Map{
				"type":    "alert-error",
				"message": fmt.Sprintf("something went wrong: %s", err),
			}

			return flash.WithError(c, fm).Redirect("/register")
		}

		fm := fiber.Map{
			"type":    "alert-success",
			"message": "You have successfully registered!!",
		}

		return flash.WithSuccess(c, fm).Redirect("/login")
	}

	return c.Render("register", fiber.Map{
		"Page":          "Register",
		"FromProtected": fromProtected,
		"Message":       flash.Get(c),
	})
}

// Authentication Middleware
func AuthMiddleware(c *fiber.Ctx) error {
	fm := fiber.Map{
		"type": "alert-error",
	}

	session, err := store.Get(c)
	if err != nil {
		fm["message"] = "You are not authorized"

		return flash.WithError(c, fm).Redirect("/login")
	}

	if session.Get(AUTH_KEY) == nil {
		fm["message"] = "You are not authorized"

		return flash.WithError(c, fm).Redirect("/login")
	}

	userId := session.Get(USER_ID)
	if userId == nil {
		fm["message"] = "You are not authorized"

		return flash.WithError(c, fm).Redirect("/login")
	}

	user, err := models.GetUserById(fmt.Sprint(userId.(uint64)))
	if err != nil {
		fm["message"] = "You are not authorized"

		return flash.WithError(c, fm).Redirect("/login")
	}

	c.Locals("userId", userId)
	c.Locals("username", user.Username)
	fromProtected = true

	return c.Next()
}

// Logout Handler
func HandleLogout(c *fiber.Ctx) error {
	fm := fiber.Map{
		"type": "alert-error",
	}

	session, err := store.Get(c)
	if err != nil {
		fm["message"] = "logged out (no session)"

		return flash.WithError(c, fm).Redirect("/login")
	}

	err = session.Destroy()
	if err != nil {
		fm["message"] = fmt.Sprintf("something went wrong: %s", err)

		return flash.WithError(c, fm).Redirect("/login")
	}

	fm = fiber.Map{
		"type":    "alert-success",
		"message": "You have successfully logged out!!",
	}

	fromProtected = false

	return flash.WithSuccess(c, fm).Redirect("/login")
}

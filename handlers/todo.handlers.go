package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/emarifer/gofiber-htmx-todolist/models"
	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/flash"
)

/********** Handlers for Todo Views **********/

// Render List Page with success/error messages
func HandleViewList(c *fiber.Ctx) error {
	todo := new(models.Todo)
	todo.CreatedBy = c.Locals("userId").(uint64)

	fm := fiber.Map{
		"type": "alert-error",
	}

	todosSlice, err := todo.GetAllTodos()
	if err != nil {
		fm["message"] = fmt.Sprintf("something went wrong: %s", err)

		return flash.WithError(c, fm).Redirect("/todo/list")
	}

	return c.Render("todo/index", fiber.Map{
		"Page":          "Tasks List",
		"FromProtected": fromProtected,
		"Todos":         todosSlice,
		"UserID":        c.Locals("userId").(uint64),
		"Username":      c.Locals("username").(string),
		"Message":       flash.Get(c),
	})
}

// Render Create Todo Page with success/error messages
func HandleViewCreatePage(c *fiber.Ctx) error {

	if c.Method() == "POST" {
		fm := fiber.Map{
			"type": "alert-error",
		}

		newTodo := new(models.Todo)
		newTodo.CreatedBy = c.Locals("userId").(uint64)
		newTodo.Title = strings.Trim(c.FormValue("title"), " ")
		newTodo.Description = strings.Trim(c.FormValue("description"), " ")

		if _, err := newTodo.CreateTodo(); err != nil {
			fm["message"] = fmt.Sprintf("something went wrong: %s", err)

			return flash.WithError(c, fm).Redirect("/todo/list")
		}

		return c.Redirect("/todo/list")
	}

	return c.Render("todo/create", fiber.Map{
		"Page":          "Create Todo",
		"FromProtected": fromProtected,
		"UserID":        c.Locals("userId").(uint64),
		"Username":      c.Locals("username").(string),
	})
}

// Render Edit Todo Page with success/error messages
func HandleViewEditPage(c *fiber.Ctx) error {
	idParams, _ := strconv.Atoi(c.Params("id"))
	todoId := uint64(idParams)

	todo := new(models.Todo)
	todo.ID = todoId
	todo.CreatedBy = c.Locals("userId").(uint64)

	fm := fiber.Map{
		"type": "alert-error",
	}

	recoveredTodo, err := todo.GetNoteById()
	if err != nil {
		fm["message"] = fmt.Sprintf("something went wrong: %s", err)

		return flash.WithError(c, fm).Redirect("/todo/list")
	}

	if c.Method() == "POST" {
		todo.Title = strings.Trim(c.FormValue("title"), " ")
		todo.Description = strings.Trim(c.FormValue("description"), " ")
		if c.FormValue("status") == "on" {
			todo.Status = true
		} else {
			todo.Status = false
		}

		_, err := todo.UpdateTodo()
		if err != nil {
			fm["message"] = fmt.Sprintf("something went wrong: %s", err)

			return flash.WithError(c, fm).Redirect("/todo/list")
		}

		fm = fiber.Map{
			"type":    "alert-success",
			"message": "Task successfully updated!!",
		}

		return flash.WithSuccess(c, fm).Redirect("/todo/list")
	}

	return c.Render("todo/update", fiber.Map{
		"Page":          fmt.Sprintf("Edit Todo #%d", recoveredTodo.ID),
		"FromProtected": fromProtected,
		"ID":            recoveredTodo.ID,
		"Title":         recoveredTodo.Title,
		"Description":   recoveredTodo.Description,
		"Status":        recoveredTodo.Status,
		"UserID":        c.Locals("userId").(uint64),
		"Username":      c.Locals("username").(string),
	})
}

// Handler Remove Todo
func HandleDeleteTodo(c *fiber.Ctx) error {
	idParams, _ := strconv.Atoi(c.Params("id"))
	todoId := uint64(idParams)

	todo := new(models.Todo)
	todo.ID = todoId
	todo.CreatedBy = c.Locals("userId").(uint64)

	fm := fiber.Map{
		"type": "alert-error",
	}

	if err := todo.DeleteTodo(); err != nil {
		fm["message"] = fmt.Sprintf("something went wrong: %s", err)

		return flash.WithError(c, fm).Redirect(
			"/todo/list",
			fiber.StatusSeeOther,
		)
	}

	fm = fiber.Map{
		"type":    "alert-success",
		"message": "Task successfully deleted!!",
	}

	return flash.WithSuccess(c, fm).Redirect("/todo/list", fiber.StatusSeeOther)
}

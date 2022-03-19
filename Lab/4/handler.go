package main

import (
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

type Handler struct {
	service Service
}

func (h *Handler) indexTask(c echo.Context) error {
	tasks, err := h.service.getAllTask()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, response{
			Data: tasks,
		})
	}

	responses := make([]*taskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = &taskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Completed:   task.Completed,
		}
	}

	return c.JSON(http.StatusOK, response{
		Data: responses,
	})
}

func (h *Handler) viewTask(c echo.Context) error {
	task, err := h.service.getTask(c.Get("id").(int))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, response{
			Error: "Cannot get task",
		})
	}

	return c.JSON(http.StatusOK, response{
		Data: taskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Completed:   task.Completed,
		},
	})
}

type createTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (h *Handler) createTask(c echo.Context) error {
	var req createTaskRequest
	if err := c.Bind(&req); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, response{
			Error: "Cannot bind request",
		})
	}

	task, err := h.service.addTask(&Task{
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, response{
			Error: "Cannot create task",
		})
	}

	return c.JSON(http.StatusCreated, response{
		Data: taskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Completed:   task.Completed,
		},
	})
}

type updateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func (h *Handler) updateTask(c echo.Context) error {
	var req updateTaskRequest
	if err := c.Bind(&req); err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnprocessableEntity, response{
			Error: "Cannot bind request",
		})
	}

	task, err := h.service.updateTask(c.Get("id").(int), &Task{
		Title:       req.Title,
		Description: req.Description,
		Completed:   req.Completed,
	})
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, response{
			Error: "Cannot update task",
		})
	}

	return c.JSON(http.StatusOK, response{
		Data: taskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Completed:   task.Completed,
		},
	})
}

func (h *Handler) deleteTask(c echo.Context) error {
	if err := h.service.deleteTask(c.Get("id").(int)); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, response{
			Error: "Cannot delete task",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) uploadFile(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnprocessableEntity, response{
			Error: "Cannot get file",
		})
	}

	var closingErrors []error
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			closingErrors = append(closingErrors, err)
		}
	}(src)

	// Destination
	dst, err := os.Create(filepath.Join("./files", filepath.Base(file.Filename)))
	if err != nil {
		return err
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			closingErrors = append(closingErrors, err)
		}
	}(dst)

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, response{
			Error: "Cannot copy file",
		})
	}

	if closingErrors != nil {
		log.Println(closingErrors[0])
		log.Println(closingErrors[1])
		return c.JSON(http.StatusInternalServerError, response{
			Error: "Cannot close file",
		})
	}
	return c.NoContent(http.StatusNoContent)
}

type taskResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"name"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type response struct {
	Data  interface{} `json:"data,omitempty"`
	Error interface{} `json:"error,omitempty"`
}


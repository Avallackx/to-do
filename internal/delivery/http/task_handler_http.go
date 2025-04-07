package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"todo-app/internal/model"
	"todo-app/internal/utils"
)

type TaskHTTPHandler struct {
	TaskUsecase model.TaskUsecase
}

func NewTaskHTTPHandler(e *echo.Echo, tu model.TaskUsecase) {
	handler := TaskHTTPHandler{TaskUsecase: tu}

	g := e.Group("/v1")
	g.POST("/tasks", handler.CreateTask)
	g.GET("/tasks", handler.FetchTasks)
	g.GET("/tasks/:ID", handler.FetchTaskByID)
	g.PUT("/tasks", handler.UpdateTask)
	g.DELETE("/tasks/:ID", handler.DeleteTaskByID)
}

func (th *TaskHTTPHandler) CreateTask(c echo.Context) error {
	input := new(model.CreateTaskInput)
	if err := c.Bind(input); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	task, err := th.TaskUsecase.Create(c.Request().Context(), input.ToModel())
	if err != nil {
		logrus.Error(err)
		return c.JSON(utils.ParseHTTPErrorStatusCode(err), err.Error())
	}

	return c.JSON(http.StatusCreated, task)
}

func (th *TaskHTTPHandler) DeleteTaskByID(c echo.Context) error {
	ID, err := strconv.ParseInt(c.Param("ID"), 10, 64)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, "ID param is invalid")
	}

	err = th.TaskUsecase.DeleteByID(c.Request().Context(), ID)
	if err != nil {
		logrus.Error(err)
		return c.JSON(utils.ParseHTTPErrorStatusCode(err), err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (th *TaskHTTPHandler) FetchTasks(c echo.Context) error {
	queryParams := new(model.GetTasksQueryParams)

	if err := c.Bind(queryParams); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	tasks, count, err := th.TaskUsecase.FindAll(c.Request().Context(), *queryParams)
	if err != nil {
		logrus.Error(err)
		return c.JSON(utils.ParseHTTPErrorStatusCode(err), err.Error())
	}

	return c.JSON(http.StatusOK, model.NewPaginationResponse(
		tasks,
		queryParams.Page,
		queryParams.Size,
		count,
	))
}

func (th *TaskHTTPHandler) FetchTaskByID(c echo.Context) error {
	ID, err := strconv.ParseInt(c.Param("ID"), 10, 64)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, "ID param is invalid")
	}

	task, err := th.TaskUsecase.FindByID(c.Request().Context(), ID)
	if err != nil {
		logrus.Error(err)
		return c.JSON(utils.ParseHTTPErrorStatusCode(err), err.Error())
	}

	return c.JSON(http.StatusOK, task)
}

func (th *TaskHTTPHandler) UpdateTask(c echo.Context) error {
	input := new(model.UpdateTaskInput)
	if err := c.Bind(input); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	task, err := th.TaskUsecase.Update(c.Request().Context(), input.ToModel())
	if err != nil {
		logrus.Error(err)
		return c.JSON(utils.ParseHTTPErrorStatusCode(err), err.Error())
	}

	return c.JSON(http.StatusOK, task)
}

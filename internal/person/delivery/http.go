package delivery

import (
	"errors"
	"fio/internal/enrich"
	"fio/internal/person"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const (
	ErrInternal    = "Internal error"
	ErrBadRequest  = "Bad request"
	ErrIncorrectID = "incorrect id"
	ErrNotFound    = "not found"
)

// PersonRepo описывает методы работы с репозиторием persons
type PersonRepo interface {
	GetPersons(page int, filters map[string]string) ([]person.Person, error)
	Get(id int) (*person.Person, error)
	Delete(id int) error
	Create(person *person.Person) error
	Update(person *person.Person) error
}

type PersonHandler struct {
	Logger           *logrus.Entry
	EnrichService    *enrich.Service
	PersonRepo       PersonRepo
	WhiteListFilters []string
}

type Response struct {
	Status bool   `json:"status"`
	Error  string `json:"error,omitempty"`
}

func Error(err string) *Response {
	return &Response{
		Error: err,
	}
}

func (h *PersonHandler) getFilters(c *gin.Context) map[string]string {

	filters := c.QueryMap("filters")
	data := make(map[string]string)

	for i := range h.WhiteListFilters {
		if v, ok := filters[h.WhiteListFilters[i]]; ok {
			data[h.WhiteListFilters[i]] = v
		}
	}

	return data

}

func (h *PersonHandler) GetListHandler(c *gin.Context) {

	var err error
	var page int

	pageString := c.DefaultQuery("page", "1")
	page, err = strconv.Atoi(pageString)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Error(ErrInternal))
		return
	}

	if page < 1 {
		page = 1
	}

	filters := h.getFilters(c)

	persons, err := h.PersonRepo.GetPersons(page, filters)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Error(ErrInternal))
		return
	}

	c.JSON(http.StatusOK, persons)

}

func (h *PersonHandler) UpdateHandler(c *gin.Context) {

	var err error
	var p person.Person

	idParam := c.Param("id")

	if idParam == "" {
		c.JSON(http.StatusBadRequest, Error(ErrIncorrectID))
	}

	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, Error(ErrIncorrectID))
		return
	}

	err = c.ShouldBindJSON(&p)

	if err != nil {
		h.Logger.Info(err)
		c.JSON(http.StatusBadRequest, Error(ErrBadRequest))
		return
	}

	p.ID = id

	err = h.PersonRepo.Update(&p)

	if err != nil {
		c.JSON(http.StatusNotFound, Error(ErrNotFound))
		return
	}

	c.JSON(http.StatusOK, p)

}

func (h *PersonHandler) DeleteHandler(c *gin.Context) {

	var err error
	idParam := c.Param("id")

	if idParam == "" {
		c.JSON(http.StatusBadRequest, Error(ErrIncorrectID))
	}

	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, Error(ErrIncorrectID))
		return
	}

	err = h.PersonRepo.Delete(id)

	if errors.Is(err, person.ErrNotFound) {
		c.JSON(http.StatusInternalServerError, Error("not found"))
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, Error(ErrInternal))
		return
	}

	c.JSON(http.StatusNoContent, nil)

}

func (h *PersonHandler) CreateHandler(c *gin.Context) {

	var p person.Person
	err := c.ShouldBindJSON(&p)

	if err != nil {
		h.Logger.Info(err)
		c.JSON(http.StatusBadRequest, Error(ErrBadRequest))
		return
	}

	h.EnrichService.EnrichPerson(&p)

	if err := h.PersonRepo.Create(&p); err != nil {
		c.JSON(http.StatusInternalServerError, Error(ErrInternal))
	}

	c.JSON(http.StatusCreated, p)

}

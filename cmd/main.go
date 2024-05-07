package main

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gookit/goutil/dump"
	"github.com/labstack/echo/v4"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	pt_BR_translations "github.com/go-playground/validator/v10/translations/pt_BR"
)

var validate *validator.Validate
var uni *ut.UniversalTranslator
var translator ut.Translator

type Category struct {
	ID   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

func createCategory(c echo.Context) error {
	category := new(Category)
	if err := c.Bind(category); err != nil {
		return err
	}

	err := validate.Struct(category)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			//fmt.Println(err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		translator, _ = c.Get("translator").(ut.Translator)
		//dump.Println(translator)
		//dump.Println(err)
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Translate(translator))
		}
		// Retorne um erro HTTP 400 para solicitações malformadas.
		//return c.JSON(http.StatusBadRequest, err.Error())
		return c.JSON(http.StatusBadRequest, errors)
	}

	// Aqui você pode adicionar o código para salvar a categoria no banco de dados.
	// Por enquanto, vamos apenas retornar a categoria que foi enviada.
	return c.JSON(http.StatusCreated, category)
}

func LanguageMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		lang := c.Request().Header.Get("Accept-Language")
		translator, found := uni.GetTranslator(lang)
		if !found {
			// Se o idioma solicitado não for suportado, use o português como padrão.
			translator, _ = uni.GetTranslator("pt_BR")
		}
		c.Set("translator", translator)
		return next(c)
	}
}

func main() {
	//validate = validator.New()
	initValidator()

	e := echo.New()
	e.Use(LanguageMiddleware)

	e.POST("/categories", createCategory)

	e.Start(":8080")

}

func initValidator() {
	pt_BR := pt_BR.New()
	english := en.New()

	uni = ut.New(pt_BR, pt_BR, english)

	validate = validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {

		//dump.Print(field)
		dump.Println(field)
		return field.Name
		//return faTranslation[field.Name]
	})

	// Register the pt_BR translations
	translator, found := uni.GetTranslator("pt_BR")
	if found {
		pt_BR_translations.RegisterDefaultTranslations(validate, translator)
	}

	// Register the en translations
	translator, found = uni.GetTranslator("en")
	if found {
		en_translations.RegisterDefaultTranslations(validate, translator)
	}
}

type Test struct {
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
}

func TestTranslation() {
	pt_BR := pt_BR.New()
	english := en.New()

	uni := ut.New(pt_BR, pt_BR, english)

	validate := validator.New()

	// Register the pt_BR translations
	translator, found := uni.GetTranslator("pt_BR")
	if found {
		pt_BR_translations.RegisterDefaultTranslations(validate, translator)
	}

	// Register the en translations
	translator, found = uni.GetTranslator("en")
	if found {
		en_translations.RegisterDefaultTranslations(validate, translator)
	}

	t := &Test{}

	err := validate.Struct(t)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}
		translator, _ = uni.GetTranslator("pt_BR")
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err.Translate(translator))
		}
	}
}

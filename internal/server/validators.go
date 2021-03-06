package server

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/sergeychur/avito_auto/internal/models"
	"net/http"
	"time"
)

type Validator struct {
	TimeOut int
}

func NewValidator(timeOut int) *Validator {
	validator := new(Validator)
	validator.TimeOut = timeOut
	return validator
}

func (v *Validator)ValidateLink(/*timeOut int,*/ link models.Link) error {
	err := ValidateFormat(link.RealURL)
	if err != nil {
		return err
	}
	err = ValidateURLExists(time.Duration(v.TimeOut), link.RealURL)
	if err != nil {
		return err
	}
	return nil
}

func ValidateURLExists(timeOut time.Duration, url string) error {
	timeout := time.Duration(timeOut * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	_, err := client.Head(url)
	if err != nil {
		return err
	}
	return nil
}

func ValidateFormat(link string) error {
	if govalidator.IsURL(link) {
		return nil
	}
	return errors.New("not url")
}
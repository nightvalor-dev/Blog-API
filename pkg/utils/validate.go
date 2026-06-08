package utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func DecodeAndValidate(r *http.Request, dst any) error {

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(dst)
	if err != nil {
		return err
	}

	return Validate.Struct(dst)
}

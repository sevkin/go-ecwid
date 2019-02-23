package ecwid

import (
	"encoding/json"
	"errors"

	resty "gopkg.in/resty.v1"
)

func responseUnmarshal(response *resty.Response, err error, result interface{}) error {
	if err != nil {
		return err
	}

	if response.StatusCode() != 200 {
		return errorResponse(response)
	}

	err = json.Unmarshal(response.Body(), &result)
	if err != nil {
		return err
	}

	return nil
}

func responseAdd(response *resty.Response, err error) (uint64, error) {
	var result struct {
		ID uint64 `json:"id"`
	}

	if err := responseUnmarshal(response, err, &result); err != nil {
		return 0, err
	}
	return result.ID, nil
}

func responseUpdate(response *resty.Response, err error) error {
	var result struct {
		UpdateCount uint `json:"updateCount"`
	}

	if err := responseUnmarshal(response, err, &result); err != nil {
		return err
	}

	if result.UpdateCount != 1 {
		return errors.New("no updated")
	}
	return nil
}

func responseDelete(response *resty.Response, err error) error {
	var result struct {
		DeleteCount uint `json:"deleteCount"`
	}

	if err := responseUnmarshal(response, err, &result); err != nil {
		return err
	}

	if result.DeleteCount != 1 {
		return errors.New("no deleted")
	}
	return nil
}

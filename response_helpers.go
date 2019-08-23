package ecwid

import (
	"encoding/json"
	"errors"

	"github.com/go-resty/resty/v2"
)

func errorResponse(response *resty.Response) error {
	var result struct {
		ErrorMessage string `json:"errorMessage"`
	}
	if err := json.Unmarshal(response.Body(), &result); err == nil && len(result.ErrorMessage) > 0 {
		if result.ErrorMessage != "" {
			return errors.New(result.ErrorMessage)
		}
	}
	return errors.New(response.Status())
}

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

func responseUpdateCount(response *resty.Response, err error) (int, error) {
	var result struct {
		UpdateCount int `json:"updateCount"`
	}

	return result.UpdateCount, responseUnmarshal(response, err, &result)
}

func responseUpdate(response *resty.Response, err error) error {
	count, err := responseUpdateCount(response, err)
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("no updated")
	}
	return nil
}

func responseDelete(response *resty.Response, err error) (uint, error) {
	var result struct {
		DeleteCount uint `json:"deleteCount"`
	}

	if err := responseUnmarshal(response, err, &result); err != nil {
		return 0, err
	}

	if result.DeleteCount == 0 {
		return 0, errors.New("no deleted")
	}
	return result.DeleteCount, nil
}

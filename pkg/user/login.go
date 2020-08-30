package user

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func Login(phonenumber string) error {
	resp, err := http.Post("http://localhost:8000/v1/login/"+phonenumber, "application/json", nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Failed to call login: " + string(body))
	}
	return nil
}

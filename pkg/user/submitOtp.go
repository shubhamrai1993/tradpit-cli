package user

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func SubmitOtp(phonenumber string, otp string) (string, error) {
	resp, err := http.Post("http://localhost:8000/v1/login/"+phonenumber+"/otp/"+otp, "application/json", nil)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("Failed to call login: " + string(body))
	} else {
		fmt.Println("Login successful")
		return string(body), nil
	}
}

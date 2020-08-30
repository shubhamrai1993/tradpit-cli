package user

import (
	"fmt"
	kiteconnect "github.com/zerodhatech/gokiteconnect"
	"io/ioutil"
	"net/http"
	"os"
)

var apiKey string = os.Getenv("API_KEY")

func CheckKiteLogin(phonenumber string, accessToken string, kiteUserID string) error {
	request, err := http.NewRequest("POST", "http://localhost:8000/v1/init/"+phonenumber+"/kite/"+kiteUserID, nil)
	if err != nil {
		return err
	}
	request.Header.Set("x-auth-token", "abcdef")

	client := http.DefaultClient
	resp, err := client.Do(request)

	if err != nil {
		return err
	}

	_, err = ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 401 {
		kc := kiteconnect.New(apiKey)
		loginUrl := kc.GetLoginURL()
		fmt.Println("Your kite account has not been connected to your tradpit account. Please login using this URL: " + loginUrl)
	} else if resp.StatusCode == 200 {
		// Create the code skeleton to be used for pushing the trading code
		fmt.Println("Your tradpit codebase is ready")
	} else {
		return fmt.Errorf("Failed to login")
	}
	return nil
}

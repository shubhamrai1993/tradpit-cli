package cmd

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"tradpit.com/tradpit-cli/pkg/user"
)

var (
	rootCmd = &cobra.Command{
		Use:   "tp",
		Short: "Tradpit cli for seamless building and deployment of trading strategies",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Tradpit cli for seamless building and deployment of trading strategies")
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints the version of the tradpit cli",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("v0.1")
		},
	}

	phoneNumber string

	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Log into the tradpit platform with your phone number",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := user.Login(phoneNumber)
			if err != nil {
				return err
			}
			prompt := promptui.Prompt{
				Label: "OTP",
				Validate: func(otp string) error {
					if len(otp) != 6 {
						return errors.New("Invalid otp format")
					}
					return nil
				},
			}
			otp, err := prompt.Run()
			if err != nil {
				return err
			}
			accessToken, err := user.SubmitOtp(phoneNumber, otp)
			if err != nil {
				return err
			}
			home, err := homedir.Dir()
			if err != nil {
				return err
			}

			f, err := os.Create(home + "/.tradpit.yml")
			defer f.Close()
			if err != nil {
				return err
			}
			f.WriteString("PHONE_NUMBER:" + phoneNumber + "\n")
			f.WriteString("ACCESS_TOKEN:" + accessToken + "\n")

			return nil
		},
	}

	kiteUserID string

	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialise the codebase for tradpit bot",
		RunE: func(cmd *cobra.Command, args []string) error {

			phoneNumber := viper.GetString("PHONE_NUMBER")
			accessToken := viper.GetString("ACCESS_TOKEN")

			err := user.CheckKiteLogin(phoneNumber, accessToken, kiteUserID)
			if err != nil {
				return err
			}
			return nil
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVar(&phoneNumber, "phonenumber", "", "Phone number of your account")
	loginCmd.MarkFlagRequired("phonenumber")

	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&kiteUserID, "kite-user-id", "", "User ID of your kite connect account")
	initCmd.MarkFlagRequired("kite-user-id")
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
	}

	viper.AddConfigPath(home)
	viper.SetConfigName(".tradpit")
	viper.SetConfigType("yml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}
}

package main

import (
	"fmt"
	"os"

	"time"

	"git.containerum.net/ch/grpc-proto-files/auth"
	"git.containerum.net/ch/user-manager/clients"
	"git.containerum.net/ch/user-manager/models"
	"git.containerum.net/ch/user-manager/routes"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func exitOnErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	viper.SetEnvPrefix("ch_user")
	viper.AutomaticEnv()
	exitOnErr(setupLogger())

	app := gin.New()
	app.Use(gin.RecoveryWithWriter(logrus.StandardLogger().WithField("component", "gin_recovery").WriterLevel(logrus.ErrorLevel)))
	app.Use(ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, true))

	db, err := models.DBConnect(viper.GetString("pg_url"))
	exitOnErr(err)
	defer db.Close()

	mailClient := clients.NewMailClient(viper.GetString("mail_url"))

	reCaptchaClient := clients.NewReCaptchaClient(viper.GetString("recaptcha_key"))

	clients.RegisterOAuthClient(clients.NewGithubOAuthClient(viper.GetString("github_app_id"), viper.GetString("github_secret")))
	clients.RegisterOAuthClient(clients.NewGoogleOAuthClient(viper.GetString("google_app_id"), viper.GetString("google_secret")))
	clients.RegisterOAuthClient(clients.NewFacebookOAuthClient(viper.GetString("facebook_app_id"), viper.GetString("facebook_secret")))

	authConn, err := grpc.Dial(viper.GetString("auth_grpc_addr"), grpc.WithInsecure())
	exitOnErr(err)
	defer authConn.Close()
	authClient := auth.NewAuthClient(authConn)

	webAPIClient := clients.NewWebAPIClient(viper.GetString("web_api_url"))

	routes.SetupRoutes(app, routes.Services{
		MailClient:      mailClient,
		DB:              db,
		AuthClient:      authClient,
		ReCaptchaClient: reCaptchaClient,
		WebAPIClient:    webAPIClient,
	})

	exitOnErr(app.Run(getListenAddr()))
}

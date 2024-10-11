/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"pelucio/driver"
	"pelucio/x/xerrors"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// @title          Pelucio backend API
// @version        1.0
// @description    Pelucio backend http endpoints
// @contact.name   4 Engine Inc
// @contact.url    https://fourengine.com/
// @contact.email  contato@fourengine.com.br
// @host      localhost:8091
// @BasePath  /api
// @securityDefinitions.basic  BasicAuth
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Runs v1 and v2 apis",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		d := &driver.Driver{}
		r := gin.Default()

		corsConfig := cors.Config{
			AllowAllOrigins:  false,
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
			AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
			ExposeHeaders:    []string{"*"},
			AllowCredentials: true,
		}
		c := cors.New(corsConfig)

		r.Use(c)

		openv1 := r.Group("api/v1/open", xerrors.HandleErrorMiddleware)
		{
			d.WalletHandler().RegisterOpenRoutes(openv1)
		}

		adminv1 := r.Group("api/v1/admin", xerrors.HandleErrorMiddleware)
		{
			d.WalletHandler().RegisterAdminRoutes(adminv1)
		}

		r.Run(d.Config().HttpPort())
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
}

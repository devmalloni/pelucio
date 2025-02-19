/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"pelucio/driver"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Runs migrate",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		d := &driver.Driver{}

		db, err := bolt.Open(d.Config().PATH(), 0600, nil)
		if err != nil {
			panic(err)
		}

		tx, err := db.Begin(true)
		if err != nil {
			panic(err)
		}
		defer tx.Rollback()

		_, err = tx.CreateBucket([]byte("Wallets"))
		if err != nil {
			panic(err)
		}

		_, err = tx.CreateBucket([]byte("WalletRecords"))
		if err != nil {
			panic(err)
		}

		_, err = tx.CreateBucket([]byte("WalletTransactions"))
		if err != nil {
			panic(err)
		}

		if err := tx.Commit(); err != nil {
			panic(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

package main

import (
	"os"
	"github.com/yoavl/coucher/database"
	"github.com/spf13/cobra"
	"github.com/yoavl/coucher/revclean"
)

func main() {
	db := &database.Database{
		BaseUrl:  os.Getenv("COUCH_URL"),
		Name:     os.Getenv("COUCH_DBNAME"),
		Username: os.Getenv("COUCH_USER"),
		Password: os.Getenv("COUCH_PASS"),
	}
	var rootCmd = &cobra.Command{Use: "coucher"}
	rootCmd.PersistentFlags().StringVarP(&db.BaseUrl, "couch-url", "c", db.BaseUrl,
		"The URL of the couchdb, e.g.: http://localhost:5984")
	rootCmd.PersistentFlags().StringVarP(&db.Name, "database", "d", db.Name, "Database name, e.g.: db1")
	rootCmd.PersistentFlags().StringVarP(&db.Username, "username", "u", db.Username, "User name")
	rootCmd.PersistentFlags().StringVarP(&db.Password, "password", "p", db.Password, "Password")

	var revclean = revclean.NewCmd(db)
	rootCmd.AddCommand(revclean)
	rootCmd.Execute()
}

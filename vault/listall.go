package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listallCmd represents the listall command.
var listallCmd = &cobra.Command{
	Use:   "listall",
	Short: "A way to list all secrets",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		encodingKey := viper.GetString("encoding_key")
		filePath := viper.GetString("filepath")

		fv, err := NewFileVault(encodingKey, filePath)
		if err != nil {
			log.Fatalf("failed to create new file vault with encoding key '%s' and filepath '%s', err: %s", encodingKey, filePath, err)
		}

		if err := fv.ListAll(); err != nil {
			log.Fatalf("failed to list all secrets, err: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listallCmd)
}

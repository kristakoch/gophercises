package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command.
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A way to get a secret",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		key, err := cmd.Flags().GetString("key")
		if err != nil {
			log.Fatal(err)
		}

		encodingKey := viper.GetString("encoding_key")
		filePath := viper.GetString("filepath")

		fv, err := NewFileVault(encodingKey, filePath)
		if err != nil {
			log.Fatalf("failed to create new file vault with encoding key '%s' and filepath '%s', err: %s", encodingKey, filePath, err)
		}

		value, err := fv.Get(key)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("===> %s", value)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	var key string
	getCmd.PersistentFlags().StringVarP(&key, "key", "k", "", "key value to look up secret by")
}

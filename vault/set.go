package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setCmd represents the set command.
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "A way to set a secret",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		key, err := cmd.Flags().GetString("key")
		if err != nil {
			log.Fatal(err)
		}

		val, err := cmd.Flags().GetString("val")
		if err != nil {
			log.Fatal(err)
		}

		encodingKey := viper.GetString("encoding_key")
		filePath := viper.GetString("filepath")

		fv, err := NewFileVault(encodingKey, filePath)
		if err != nil {
			log.Fatalf("failed to create new file vault with encoding key '%s' and filepath '%s', err: %s", encodingKey, filePath, err)
		}

		if err := fv.Set(key, val); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	var key string
	setCmd.PersistentFlags().StringVarP(&key, "key", "k", "", "key value to store secret with")

	var val string
	setCmd.PersistentFlags().StringVarP(&val, "val", "v", "", "value for secret")
}

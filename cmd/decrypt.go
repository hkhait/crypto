/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/hkhait/crypto/encryption/pbe"
	"github.com/spf13/cobra"
)

var ciphertext string

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt data locally",
	Long:  `Decrypt CMS On-Ramp data. Read ciphertext from stdin or CLI flag and output plaintext to stdout.`,
	Run: func(cmd *cobra.Command, args []string) {
		w := bufio.NewWriter(os.Stdout)
		defer w.Flush()

		if ciphertext == "-" {
			r := bufio.NewReader(os.Stdin)
			s := bufio.NewScanner(r)
			for s.Scan() {
				decryptor := pbe.NewStandardPBEStringEncryptor()
				decryptor.SetPassword(password)
				plaintext, err := decryptor.Decrypt(s.Text())
				if err != nil {
					panic(err)
				}
				fmt.Fprintln(w, plaintext)
				if err != nil {
					panic(err)
				}
			}
		} else {
			encryptor := pbe.NewStandardPBEStringEncryptor()
			encryptor.SetPassword(password)
			plaintext, err := encryptor.Decrypt(ciphertext)
			if err != nil {
				panic(err)
			}
			fmt.Fprintln(w, plaintext)
			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// decryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	decryptCmd.Flags().StringVar(&ciphertext, "ciphertext", "-", "the ciphertext to be decrypted, read from stdin if not specified")
}

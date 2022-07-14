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

var plaintext string

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt data locally",
	Long:  `Encrypt CMS On-Ramp data. Read plaintext from stdin or CLI flag and output ciphertext to stdout.`,
	Run: func(cmd *cobra.Command, args []string) {
		w := bufio.NewWriter(os.Stdout)
		defer w.Flush()

		if plaintext == "-" {
			r := bufio.NewReader(os.Stdin)
			s := bufio.NewScanner(r)
			for s.Scan() {
				encryptor := pbe.NewStandardPBEStringEncryptor()
				encryptor.SetPassword(password)
				ciphertext, err := encryptor.Encrypt(s.Text())
				if err != nil {
					panic(err)
				}
				fmt.Fprintln(w, ciphertext)
				if err != nil {
					panic(err)
				}
			}
		} else {
			encryptor := pbe.NewStandardPBEStringEncryptor()
			encryptor.SetPassword(password)
			ciphertext, err := encryptor.Encrypt(plaintext)
			if err != nil {
				panic(err)
			}
			fmt.Fprintln(w, ciphertext)
			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	encryptCmd.Flags().StringVar(&plaintext, "plaintext", "-", "the plaintext to be encrypted, read from stdin if not specified")
}

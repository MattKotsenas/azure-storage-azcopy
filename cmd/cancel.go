package cmd

import (
	"github.com/spf13/cobra"
	"github.com/Azure/azure-storage-azcopy/handlers"
	"github.com/Azure/azure-storage-azcopy/common"
	"errors"
)

func init() {
	var commandLineInput common.JobID = ""

	// cancelCmd represents the pause command
	cancelCmd := &cobra.Command{
		Use:        "cancel",
		SuggestFor: []string{"cancl", "ancl", "cacl"},
		Short:      "cancel cancels the existing job for given JobId",
		Long: `cancel cancels the existing job for given JobId`,
		Args: func(cmd *cobra.Command, args []string) error {
			// the cancel command requires necessarily to have an argument
			// cancel jobId -- cancel all the parts of an existing job for given jobId

			// If no argument is passed then it is not valid
			if len(args) != 1 {
				return errors.New("this command only requires jobId")
			}
			commandLineInput = common.JobID(args[0])
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			handlers.HandleCancelCommand(commandLineInput)
		},
	}
	rootCmd.AddCommand(cancelCmd)
}
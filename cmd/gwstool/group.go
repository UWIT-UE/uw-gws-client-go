package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/uwit-ue/uw-gws-client-go/gws"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Group operations",
	Long:  "Manage groups including get, create, update, and delete operations",
}

var groupGetCmd = &cobra.Command{
	Use:   "get <group-id>",
	Short: "Get group information",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		group, err := gwsClient.GetGroup(args[0])
		if err != nil {
			return err
		}
		outputResult(group)
		return nil
	},
}

var groupCreateCmd = &cobra.Command{
	Use:   "create <group-id>",
	Short: "Create a new group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID := args[0]

		group := &gws.Group{
			ID: groupID,
		}

		if interactive {
			if displayName := promptForInput("Display Name"); displayName != "" {
				group.DisplayName = displayName
			}
			if description := promptForInput("Description"); description != "" {
				group.Description = description
			}
			if contact := promptForInput("Contact (UWNetID)"); contact != "" {
				group.Contact = gws.UWNetID(contact)
			}

			// Collect entity lists (required admin)
			fmt.Println("\nEntity permissions (at least one admin is required):")
			if adminStr := promptForInput("Admins (comma-separated, required)"); adminStr != "" {
				admins := strings.Split(adminStr, ",")
				for i := range admins {
					admins[i] = strings.TrimSpace(admins[i])
				}
				group.AddAdmin(admins...)
			}

			if updaterStr := promptForInput("Updaters (comma-separated, optional)"); updaterStr != "" {
				updaters := strings.Split(updaterStr, ",")
				for i := range updaters {
					updaters[i] = strings.TrimSpace(updaters[i])
				}
				group.AddUpdater(updaters...)
			}

			if creatorStr := promptForInput("Creators (comma-separated, optional)"); creatorStr != "" {
				creators := strings.Split(creatorStr, ",")
				for i := range creators {
					creators[i] = strings.TrimSpace(creators[i])
				}
				group.AddCreator(creators...)
			}

			if readerStr := promptForInput("Readers (comma-separated, optional)"); readerStr != "" {
				readers := strings.Split(readerStr, ",")
				for i := range readers {
					readers[i] = strings.TrimSpace(readers[i])
				}
				group.AddReader(readers...)
			}
		} else {
			// Get from flags
			displayName, _ := cmd.Flags().GetString("display-name")
			description, _ := cmd.Flags().GetString("description")
			contact, _ := cmd.Flags().GetString("contact")
			admins, _ := cmd.Flags().GetStringSlice("admin")
			updaters, _ := cmd.Flags().GetStringSlice("updater")
			creators, _ := cmd.Flags().GetStringSlice("creator")
			readers, _ := cmd.Flags().GetStringSlice("reader")

			if displayName != "" {
				group.DisplayName = displayName
			}
			if description != "" {
				group.Description = description
			}
			if contact != "" {
				group.Contact = gws.UWNetID(contact)
			}

			// Add entity lists
			if len(admins) > 0 {
				group.AddAdmin(admins...)
			}
			if len(updaters) > 0 {
				group.AddUpdater(updaters...)
			}
			if len(creators) > 0 {
				group.AddCreator(creators...)
			}
			if len(readers) > 0 {
				group.AddReader(readers...)
			}
		}

		// Validate that at least one admin is specified
		if len(group.Admins) == 0 {
			return fmt.Errorf("at least one admin must be specified. Use --admin flag or provide via interactive mode")
		}

		createdGroup, err := gwsClient.CreateGroup(group)
		if err != nil {
			return err
		}
		outputResult(createdGroup)
		return nil
	},
}

var groupUpdateCmd = &cobra.Command{
	Use:   "update <group-id>",
	Short: "Update an existing group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID := args[0]

		// First get the existing group
		group, err := gwsClient.GetGroup(groupID)
		if err != nil {
			return err
		}

		if interactive {
			if displayName := promptForInput(fmt.Sprintf("Display Name [%s]", group.DisplayName)); displayName != "" {
				group.DisplayName = displayName
			}
			if description := promptForInput(fmt.Sprintf("Description [%s]", group.Description)); description != "" {
				group.Description = description
			}
			if contact := promptForInput(fmt.Sprintf("Contact [%s]", group.Contact)); contact != "" {
				group.Contact = gws.UWNetID(contact)
			}
		} else {
			// Get basic field flags
			displayName, _ := cmd.Flags().GetString("display-name")
			description, _ := cmd.Flags().GetString("description")
			contact, _ := cmd.Flags().GetString("contact")

			if displayName != "" {
				group.DisplayName = displayName
			}
			if description != "" {
				group.Description = description
			}
			if contact != "" {
				group.Contact = gws.UWNetID(contact)
			}

			// Handle entity permission changes
			// Add entities
			if addAdmins, _ := cmd.Flags().GetStringSlice("add-admin"); len(addAdmins) > 0 {
				group.AddAdmin(addAdmins...)
			}
			if addUpdaters, _ := cmd.Flags().GetStringSlice("add-updater"); len(addUpdaters) > 0 {
				group.AddUpdater(addUpdaters...)
			}
			if addCreators, _ := cmd.Flags().GetStringSlice("add-creator"); len(addCreators) > 0 {
				group.AddCreator(addCreators...)
			}
			if addReaders, _ := cmd.Flags().GetStringSlice("add-reader"); len(addReaders) > 0 {
				group.AddReader(addReaders...)
			}

			// Remove entities
			if removeAdmins, _ := cmd.Flags().GetStringSlice("remove-admin"); len(removeAdmins) > 0 {
				group.RemoveAdmin(removeAdmins...)
			}
			if removeUpdaters, _ := cmd.Flags().GetStringSlice("remove-updater"); len(removeUpdaters) > 0 {
				group.RemoveUpdater(removeUpdaters...)
			}
			if removeCreators, _ := cmd.Flags().GetStringSlice("remove-creator"); len(removeCreators) > 0 {
				group.RemoveCreator(removeCreators...)
			}
			if removeReaders, _ := cmd.Flags().GetStringSlice("remove-reader"); len(removeReaders) > 0 {
				group.RemoveReader(removeReaders...)
			}

			// Remove all entities
			if removeAllAdmins, _ := cmd.Flags().GetBool("remove-all-admins"); removeAllAdmins {
				group.RemoveAllAdmins()
			}
			if removeAllUpdaters, _ := cmd.Flags().GetBool("remove-all-updaters"); removeAllUpdaters {
				group.RemoveAllUpdaters()
			}
			if removeAllCreators, _ := cmd.Flags().GetBool("remove-all-creators"); removeAllCreators {
				group.RemoveAllCreators()
			}
			if removeAllReaders, _ := cmd.Flags().GetBool("remove-all-readers"); removeAllReaders {
				group.RemoveAllReaders()
			}
		}

		// Validate that at least one admin remains
		if len(group.Admins) == 0 {
			return fmt.Errorf("at least one admin must remain. Cannot remove all admins")
		}

		updatedGroup, err := gwsClient.UpdateGroup(group)
		if err != nil {
			return err
		}
		outputResult(updatedGroup)
		return nil
	},
}

var groupDeleteCmd = &cobra.Command{
	Use:   "delete <group-id>",
	Short: "Delete a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID := args[0]

		confirm, _ := cmd.Flags().GetBool("confirm")
		if !confirm && interactive {
			response := promptForInput(fmt.Sprintf("Are you sure you want to delete group '%s'? (yes/no)", groupID))
			if strings.ToLower(response) != "yes" {
				fmt.Println("Operation cancelled")
				return nil
			}
		} else if !confirm {
			return fmt.Errorf("use --confirm flag to confirm deletion")
		}

		err := gwsClient.DeleteGroup(groupID)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			outputResult(map[string]string{"status": "deleted", "group": groupID})
		} else {
			fmt.Printf("Group '%s' deleted successfully\n", groupID)
		}
		return nil
	},
}

func init() {
	// Add subcommands to group command
	groupCmd.AddCommand(groupGetCmd)
	groupCmd.AddCommand(groupCreateCmd)
	groupCmd.AddCommand(groupUpdateCmd)
	groupCmd.AddCommand(groupDeleteCmd)

	// Flags for create command
	groupCreateCmd.Flags().String("display-name", "", "Display name for the group")
	groupCreateCmd.Flags().String("description", "", "Description for the group")
	groupCreateCmd.Flags().String("contact", "", "Contact UWNetID for the group")
	groupCreateCmd.Flags().StringSlice("admin", []string{}, "Admin entity IDs (required - at least one). Use multiple flags or comma-separated: --admin user1 --admin user2 or --admin user1,user2")
	groupCreateCmd.Flags().StringSlice("updater", []string{}, "Updater entity IDs. Use multiple flags or comma-separated: --updater user1 --updater user2 or --updater user1,user2")
	groupCreateCmd.Flags().StringSlice("creator", []string{}, "Creator entity IDs. Use multiple flags or comma-separated: --creator user1 --creator user2 or --creator user1,user2")
	groupCreateCmd.Flags().StringSlice("reader", []string{}, "Reader entity IDs. Use multiple flags or comma-separated: --reader user1 --reader user2 or --reader user1,user2")

	// Flags for update command
	groupUpdateCmd.Flags().String("display-name", "", "Display name for the group")
	groupUpdateCmd.Flags().String("description", "", "Description for the group")
	groupUpdateCmd.Flags().String("contact", "", "Contact UWNetID for the group")

	// Entity permission modification flags
	groupUpdateCmd.Flags().StringSlice("add-admin", []string{}, "Add admin entity IDs. Use multiple flags or comma-separated: --add-admin user1 --add-admin user2 or --add-admin user1,user2")
	groupUpdateCmd.Flags().StringSlice("add-updater", []string{}, "Add updater entity IDs. Use multiple flags or comma-separated: --add-updater user1 --add-updater user2 or --add-updater user1,user2")
	groupUpdateCmd.Flags().StringSlice("add-creator", []string{}, "Add creator entity IDs. Use multiple flags or comma-separated: --add-creator user1 --add-creator user2 or --add-creator user1,user2")
	groupUpdateCmd.Flags().StringSlice("add-reader", []string{}, "Add reader entity IDs. Use multiple flags or comma-separated: --add-reader user1 --add-reader user2 or --add-reader user1,user2")

	groupUpdateCmd.Flags().StringSlice("remove-admin", []string{}, "Remove admin entity IDs. Use multiple flags or comma-separated: --remove-admin user1 --remove-admin user2 or --remove-admin user1,user2")
	groupUpdateCmd.Flags().StringSlice("remove-updater", []string{}, "Remove updater entity IDs. Use multiple flags or comma-separated: --remove-updater user1 --remove-updater user2 or --remove-updater user1,user2")
	groupUpdateCmd.Flags().StringSlice("remove-creator", []string{}, "Remove creator entity IDs. Use multiple flags or comma-separated: --remove-creator user1 --remove-creator user2 or --remove-creator user1,user2")
	groupUpdateCmd.Flags().StringSlice("remove-reader", []string{}, "Remove reader entity IDs. Use multiple flags or comma-separated: --remove-reader user1 --remove-reader user2 or --remove-reader user1,user2")

	groupUpdateCmd.Flags().Bool("remove-all-admins", false, "Remove all admin entities")
	groupUpdateCmd.Flags().Bool("remove-all-updaters", false, "Remove all updater entities")
	groupUpdateCmd.Flags().Bool("remove-all-creators", false, "Remove all creator entities")
	groupUpdateCmd.Flags().Bool("remove-all-readers", false, "Remove all reader entities")

	// Flags for delete command
	groupDeleteCmd.Flags().Bool("confirm", false, "Confirm deletion without prompting")
}

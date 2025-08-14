package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var memberCmd = &cobra.Command{
	Use:   "member",
	Short: "Member operations",
	Long:  "Manage group membership including list, add, remove, and check operations",
}

var memberListCmd = &cobra.Command{
	Use:   "list <group-id>",
	Short: "List members of a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID := args[0]
		effective, _ := cmd.Flags().GetBool("effective")

		if effective {
			members, err := gwsClient.GetEffectiveMembership(groupID)
			if err != nil {
				return err
			}
			outputResult(members)
		} else {
			members, err := gwsClient.GetMembership(groupID)
			if err != nil {
				return err
			}
			outputResult(members)
		}
		return nil
	},
}

var memberGetCmd = &cobra.Command{
	Use:   "get <group-id> <member-id>",
	Short: "Get information about a specific member",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID := args[0]
		memberID := args[1]
		effective, _ := cmd.Flags().GetBool("effective")

		if effective {
			member, err := gwsClient.GetEffectiveMember(groupID, memberID)
			if err != nil {
				return err
			}
			outputResult(member)
		} else {
			member, err := gwsClient.GetMember(groupID, memberID)
			if err != nil {
				return err
			}
			outputResult(member)
		}
		return nil
	},
}

var memberCheckCmd = &cobra.Command{
	Use:   "check <group-id> <member-id>",
	Short: "Check if a member belongs to a group",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID := args[0]
		memberID := args[1]
		effective, _ := cmd.Flags().GetBool("effective")

		var isMember bool
		var err error

		if effective {
			isMember, err = gwsClient.IsEffectiveMember(groupID, memberID)
		} else {
			isMember, err = gwsClient.IsMember(groupID, memberID)
		}

		if err != nil {
			return err
		}

		if outputFormat == "json" {
			result := map[string]interface{}{
				"group":     groupID,
				"member":    memberID,
				"is_member": isMember,
				"effective": effective,
			}
			outputResult(result)
		} else {
			status := "not a member"
			if isMember {
				status = "is a member"
			}
			if effective {
				fmt.Printf("%s %s of %s (effective)\n", memberID, status, groupID)
			} else {
				fmt.Printf("%s %s of %s\n", memberID, status, groupID)
			}
		}
		return nil
	},
}

var memberCountCmd = &cobra.Command{
	Use:   "count <group-id>",
	Short: "Get the count of members in a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID := args[0]
		effective, _ := cmd.Flags().GetBool("effective")

		var count int
		var err error

		if effective {
			count, err = gwsClient.EffectiveMemberCount(groupID)
		} else {
			count, err = gwsClient.MemberCount(groupID)
		}

		if err != nil {
			return err
		}

		if outputFormat == "json" {
			result := map[string]interface{}{
				"group":     groupID,
				"count":     count,
				"effective": effective,
			}
			outputResult(result)
		} else {
			memberType := "members"
			if effective {
				memberType = "effective members"
			}
			fmt.Printf("%s has %d %s\n", groupID, count, memberType)
		}
		return nil
	},
}

var memberAddCmd = &cobra.Command{
	Use:   "add <group-id> <member-id>...",
	Short: "Add members to a group",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID := args[0]
		memberIDs := args[1:]

		if interactive {
			fmt.Printf("Adding members to group '%s': %s\n", groupID, strings.Join(memberIDs, ", "))
			response := promptForInput("Continue? (yes/no)")
			if strings.ToLower(response) != "yes" {
				fmt.Println("Operation cancelled")
				return nil
			}
		}

		added, err := gwsClient.AddMembers(groupID, memberIDs...)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			result := map[string]interface{}{
				"group": groupID,
				"added": added,
			}
			outputResult(result)
		} else {
			if len(added) > 0 {
				fmt.Printf("Added %d members to %s: %s\n", len(added), groupID, strings.Join(added, ", "))
			} else {
				fmt.Printf("No members were added to %s\n", groupID)
			}
		}
		return nil
	},
}

var memberRemoveCmd = &cobra.Command{
	Use:   "remove <group-id> <member-id>...",
	Short: "Remove members from a group",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID := args[0]
		memberIDs := args[1:]

		if interactive {
			fmt.Printf("Removing members from group '%s': %s\n", groupID, strings.Join(memberIDs, ", "))
			response := promptForInput("Continue? (yes/no)")
			if strings.ToLower(response) != "yes" {
				fmt.Println("Operation cancelled")
				return nil
			}
		}

		err := gwsClient.DeleteMembers(groupID, memberIDs...)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			result := map[string]interface{}{
				"group":   groupID,
				"removed": memberIDs,
				"status":  "success",
			}
			outputResult(result)
		} else {
			fmt.Printf("Removed %d members from %s: %s\n", len(memberIDs), groupID, strings.Join(memberIDs, ", "))
		}
		return nil
	},
}

var memberClearCmd = &cobra.Command{
	Use:   "clear <group-id>",
	Short: "Remove all members from a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID := args[0]

		confirm, _ := cmd.Flags().GetBool("confirm")
		if !confirm && interactive {
			response := promptForInput(fmt.Sprintf("Are you sure you want to remove ALL members from group '%s'? (yes/no)", groupID))
			if strings.ToLower(response) != "yes" {
				fmt.Println("Operation cancelled")
				return nil
			}
		} else if !confirm {
			return fmt.Errorf("use --confirm flag to confirm clearing all members")
		}

		err := gwsClient.DeleteAllMembers(groupID)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			result := map[string]interface{}{
				"group":  groupID,
				"status": "cleared",
			}
			outputResult(result)
		} else {
			fmt.Printf("All members removed from group '%s'\n", groupID)
		}
		return nil
	},
}

func init() {
	// Add subcommands to member command
	memberCmd.AddCommand(memberListCmd)
	memberCmd.AddCommand(memberGetCmd)
	memberCmd.AddCommand(memberCheckCmd)
	memberCmd.AddCommand(memberCountCmd)
	memberCmd.AddCommand(memberAddCmd)
	memberCmd.AddCommand(memberRemoveCmd)
	memberCmd.AddCommand(memberClearCmd)

	// Add effective flag to relevant commands
	memberListCmd.Flags().Bool("effective", false, "Get effective membership (includes inherited)")
	memberGetCmd.Flags().Bool("effective", false, "Get effective member information")
	memberCheckCmd.Flags().Bool("effective", false, "Check effective membership")
	memberCountCmd.Flags().Bool("effective", false, "Count effective members")

	// Add confirm flag to destructive operations
	memberClearCmd.Flags().Bool("confirm", false, "Confirm clearing all members without prompting")
}

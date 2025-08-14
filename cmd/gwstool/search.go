package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uwit-ue/uw-gws-client-go/gws"
)

var searchCmd = &cobra.Command{
	Use:   "search [search-term]",
	Short: "Search for groups (requires at least one search parameter)",
	Long: `Search for groups using various criteria. At least one search parameter must be provided.

If a search term is provided as an argument, it is used as a name search.
Wildcards (*) can be used for pattern matching in name searches.

Examples:
  gwstool search mygroup          # Search for groups with name "mygroup"
  gwstool search "my*"            # Search for groups starting with "my"
  gwstool search --stem uw_it     # Search by stem
  gwstool search --member johndoe # Search for groups containing member`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		searchParams := gws.NewSearch()

		// Get search parameters from flags or interactive input
		name, _ := cmd.Flags().GetString("name")
		stem, _ := cmd.Flags().GetString("stem")
		member, _ := cmd.Flags().GetString("member")
		owner, _ := cmd.Flags().GetString("owner")
		instructor, _ := cmd.Flags().GetString("instructor")
		affiliate, _ := cmd.Flags().GetString("affiliate")
		scope, _ := cmd.Flags().GetString("scope")
		effective, _ := cmd.Flags().GetBool("effective")

		if len(args) > 0 {
			// If a search term is provided as argument, use it as name search
			name = args[0]
		}

		// Check if at least one search parameter is provided
		hasSearchTerms := name != "" || stem != "" || member != "" || owner != "" || instructor != "" || affiliate != ""

		if !hasSearchTerms {
			if interactive {
				fmt.Println("Enter search criteria (press Enter to skip any field):")

				// Collect all possible search parameters
				if name == "" {
					if searchTerm := promptForInput("Group name"); searchTerm != "" {
						name = searchTerm
					}
				}

				if stem == "" {
					if stemTerm := promptForInput("Stem"); stemTerm != "" {
						stem = stemTerm
					}
				}

				if member == "" {
					if memberTerm := promptForInput("Member"); memberTerm != "" {
						member = memberTerm
					}
				}

				if owner == "" {
					if ownerTerm := promptForInput("Owner"); ownerTerm != "" {
						owner = ownerTerm
					}
				}

				if instructor == "" {
					if instructorTerm := promptForInput("Instructor"); instructorTerm != "" {
						instructor = instructorTerm
					}
				}

				if affiliate == "" {
					if affiliateTerm := promptForInput("Affiliate"); affiliateTerm != "" {
						affiliate = affiliateTerm
					}
				}

				if scope == "" {
					if scopeTerm := promptForInput("Scope (one, sub, all)"); scopeTerm != "" {
						scope = scopeTerm
					}
				}

				// Check if at least one parameter was provided
				hasSearchTermsAfterPrompt := name != "" || stem != "" || member != "" || owner != "" || instructor != "" || affiliate != ""
				if !hasSearchTermsAfterPrompt {
					return fmt.Errorf("at least one search parameter is required")
				}
			} else {
				return fmt.Errorf("at least one search parameter is required. Use --name, --stem, --member, --owner, --instructor, --affiliate, or provide a search term as an argument")
			}
		}

		// Set search parameters using the builder pattern
		if name != "" {
			searchParams.WithName(name)
		}
		if stem != "" {
			searchParams.WithStem(stem)
		}
		if member != "" {
			searchParams.WithMember(member)
		}
		if owner != "" {
			searchParams.WithOwner(owner)
		}
		if instructor != "" {
			searchParams.WithInstructor(instructor)
		}
		if affiliate != "" {
			searchParams.WithAffiliate(affiliate)
		}
		if scope != "" {
			searchParams.WithScope(scope)
		}
		if effective {
			searchParams.InEffectiveMembers()
		}

		results, err := gwsClient.DoSearch(searchParams)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			outputResult(results)
		} else {
			if len(results) == 0 {
				fmt.Println("No groups found")
			} else {
				for _, group := range results {
					fmt.Println(group.ID)
				}
			}
		}
		return nil
	},
}

func init() {
	// Add search parameter flags directly to search command
	searchCmd.Flags().String("name", "", "Search by group name")
	searchCmd.Flags().String("stem", "", "Search by group stem")
	searchCmd.Flags().String("member", "", "Search groups containing this member")
	searchCmd.Flags().String("owner", "", "Search groups owned by this user")
	searchCmd.Flags().String("instructor", "", "Search groups with this instructor")
	searchCmd.Flags().String("affiliate", "", "Search groups with this affiliate")
	searchCmd.Flags().String("scope", "", "Search scope (one, sub, all)")
	searchCmd.Flags().Bool("effective", false, "Search effective members (default is direct)")
}

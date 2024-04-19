package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccDataSourceItemLoginAttributes(t *testing.T) {
	ensureVaultwardenConfigured(t)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: tfConfigProvider() + tfConfigResourceFolder() + tfConfigResourceItemLogin(),
			},
			{
				Config: tfConfigProvider() + tfConfigResourceFolder() + tfConfigResourceItemLogin() + tfConfigDataItemLogin(),
				Check:  checkItemLogin("data.bitwarden_item_login.foo_data"),
			},
		},
	})
}

func TestAccDataSourceItemLoginFailsOnInexistentItem(t *testing.T) {
	ensureVaultwardenConfigured(t)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      tfConfigProvider() + tfConfigInexistentDataItemLogin(),
				ExpectError: regexp.MustCompile("Error: object not found"),
			},
		},
	})
}

func TestAccDataSourceItemLoginDeleted(t *testing.T) {
	var objectID string

	ensureVaultwardenConfigured(t)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: tfConfigProvider() + tfConfigResourceItemLoginSmall(),
				Check:  getObjectID("bitwarden_item_login.foo", &objectID),
			},
			{
				Config: tfConfigProvider() + tfConfigResourceItemLoginSmall() + tfConfigDataItemLoginWithId(objectID),
				PreConfig: func() {
					err := bwTestClient(t).DeleteObject("item", objectID)
					assert.NoError(t, err)
				},
				ExpectError: regexp.MustCompile("Error: object not found"),
			},
		},
	})
}

func TestAccDataSourceItemLoginBySearch(t *testing.T) {
	var objectID string
	resourceName := "bitwarden_item_login.foo"

	ensureVaultwardenConfigured(t)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: tfConfigProvider() + tfConfigResourceFolder() + tfConfigResourceItemLogin(),
				Check: resource.ComposeTestCheckFunc(
					checkItemLogin(resourceName),
					getObjectID(resourceName, &objectID),
				),
			},
			{
				Config: tfConfigProvider() + tfConfigResourceFolder() + tfConfigResourceItemLogin() + tfConfigDataItemLoginWithSearchAndOrg("test-username"),
				Check:  checkItemLogin("data.bitwarden_item_login.foo_data"),
			},
			{
				Config: tfConfigProvider() + tfConfigResourceFolder() + tfConfigResourceItemLogin() + tfConfigResourceItemLoginDuplicate() + tfConfigDataItemLoginWithSearchAndOrg("test-username"),
				Check:  checkItemLogin("data.bitwarden_item_login.foo_data"),
			},
			{
				Config:      tfConfigProvider() + tfConfigResourceFolder() + tfConfigResourceItemLogin() + tfConfigResourceItemLoginDuplicate() + tfConfigDataItemLoginWithSearchOnly("test-username"),
				ExpectError: regexp.MustCompile("Error: too many objects found"),
			},
			{
				Config:      tfConfigProvider() + tfConfigResourceFolder() + tfConfigResourceItemLogin() + tfConfigDataItemLoginWithSearchAndOrg("missing-item"),
				ExpectError: regexp.MustCompile("Error: no object found matching the filter"),
			},
		},
	})
}

func tfConfigDataItemLoginWithId(id string) string {
	return fmt.Sprintf(`
data "bitwarden_item_login" "foo_data" {
	provider	= bitwarden

	id 			= "%s"
}
`, id)
}

func tfConfigDataItemLoginWithSearchAndOrg(search string) string {
	return fmt.Sprintf(`
data "bitwarden_item_login" "foo_data" {
	provider	= bitwarden

	search = "%s"
	organization_id = "%s"
}
`, search, testOrganizationID)
}

func tfConfigDataItemLoginWithSearchOnly(search string) string {
	return fmt.Sprintf(`
data "bitwarden_item_login" "foo_data" {
	provider	= bitwarden

	search = "%s"
}
`, search)
}

func tfConfigResourceItemLoginDuplicate() string {
	return `
	resource "bitwarden_item_login" "foo_duplicate" {
		provider 			= bitwarden

		name 					= "another item with username 'test-username'"
		username 			= "test-username"
	}
	`
}
func tfConfigDataItemLogin() string {
	return `
data "bitwarden_item_login" "foo_data" {
	provider	= bitwarden

	id 			= bitwarden_item_login.foo.id
}
`
}

func tfConfigInexistentDataItemLogin() string {
	return `
data "bitwarden_item_login" "foo_data" {
	provider	= bitwarden

	id 			= 123456789
}
`
}

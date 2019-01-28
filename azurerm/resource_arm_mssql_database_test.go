package azurerm

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMsSqlDatabase_basic(t *testing.T) {
	resourceName := "azurerm_mssql_database.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlDatabase_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_mode"},
			},
		},
	})
}
func TestAccAzureRMMsSqlDatabase_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	resourceName := "azurerm_mssql_database.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlDatabase_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMMsSqlDatabase_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_mssql_database"),
			},
		},
	})
}

func TestAccAzureRMMsSqlDatabase_disappears(t *testing.T) {
	resourceName := "azurerm_mssql_database.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlDatabase_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(resourceName),
					testCheckAzureRMSqlDatabaseDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMMsSqlDatabase_elasticPool(t *testing.T) {
	resourceName := "azurerm_mssql_database.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlDatabase_elasticPool(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "elastic_pool_name", fmt.Sprintf("acctestep%d", ri)),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_mode"},
			},
		},
	})
}

func TestAccAzureRMMsSqlDatabase_withTags(t *testing.T) {
	resourceName := "azurerm_mssql_database.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMMsSqlDatabase_withTags(ri, location)
	postConfig := testAccAzureRMMsSqlDatabase_withTagsUpdate(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMMsSqlDatabase_dataWarehouse(t *testing.T) {
	resourceName := "azurerm_mssql_database.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMsSqlDatabase_dataWarehouse(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists("azurerm_mssql_database.test"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_mode"},
			},
		},
	})
}

func TestAccAzureRMMsSqlDatabase_restorePointInTime(t *testing.T) {
	resourceName := "azurerm_mssql_database.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMMsSqlDatabase_basic(ri, location)
	timeToRestore := time.Now().Add(15 * time.Minute)
	formattedTime := string(timeToRestore.UTC().Format(time.RFC3339))
	postCongif := testAccAzureRMMsSqlDatabase_restorePointInTime(ri, formattedTime, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config:                    preConfig,
				PreventPostDestroyRefresh: true,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(resourceName),
				),
			},
			{
				PreConfig: func() { time.Sleep(timeToRestore.Sub(time.Now().Add(-1 * time.Minute))) },
				Config:    postCongif,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(resourceName),
					testCheckAzureRMMsSqlDatabaseExists("azurerm_mssql_database.test_restore"),
				),
			},
		},
	})
}

func TestAccAzureRMMsSqlDatabase_collation(t *testing.T) {
	resourceName := "azurerm_mssql_database.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMMsSqlDatabase_basic(ri, location)
	postConfig := testAccAzureRMMsSqlDatabase_collationUpdate(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "collation", "SQL_Latin1_General_CP1_CI_AS"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "collation", "Japanese_Bushu_Kakusu_100_CS_AS_KS_WS"),
				),
			},
		},
	})
}

func TestAccAzureRMMsSqlDatabase_requestedServiceObjectiveName(t *testing.T) {
	resourceName := "azurerm_mssql_database.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMMsSqlDatabase_requestedServiceObjectiveName(ri, location, "S0")
	postConfig := testAccAzureRMMsSqlDatabase_requestedServiceObjectiveName(ri, location, "S1")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "requested_service_objective_name", "S0"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "requested_service_objective_name", "S1"),
				),
			},
		},
	})
}

func TestAccAzureRMMsSqlDatabase_threatDetectionPolicy(t *testing.T) {
	resourceName := "azurerm_mssql_database.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMMsSqlDatabase_threatDetectionPolicy(ri, location, "Enabled")
	postConfig := testAccAzureRMMsSqlDatabase_threatDetectionPolicy(ri, location, "Disabled")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "threat_detection_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "threat_detection_policy.0.state", "Enabled"),
					resource.TestCheckResourceAttr(resourceName, "threat_detection_policy.0.retention_days", "15"),
					resource.TestCheckResourceAttr(resourceName, "threat_detection_policy.0.disabled_alerts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "threat_detection_policy.0.email_account_admins", "Enabled"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_mode", "threat_detection_policy.0.storage_account_access_key"},
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "threat_detection_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "threat_detection_policy.0.state", "Disabled"),
				),
			},
		},
	})
}

func testCheckAzureRMMsSqlDatabaseExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		databaseName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).sqlDatabasesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, databaseName, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("SQL Database %q (server %q / resource group %q) was not found", databaseName, serverName, resourceGroup)
			}

			return err
		}

		return nil
	}
}

func testCheckAzureRMMsSqlDatabaseDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_mssql_database" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		databaseName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).sqlDatabasesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, serverName, databaseName, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("SQL Database %q (server %q / resource group %q) still exists: %+v", databaseName, serverName, resourceGroup, resp)
	}

	return nil
}

func testCheckAzureRMSqlDatabaseDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serverName := rs.Primary.Attributes["server_name"]
		databaseName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).sqlDatabasesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if _, err := client.Delete(ctx, resourceGroup, serverName, databaseName); err != nil {
			return fmt.Errorf("Bad: Delete on sqlDatabasesClient: %+v", err)
		}

		return nil
	}
}

func TestAccAzureRMMsSqlDatabase_bacpac(t *testing.T) {
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMMsSqlDatabase_bacpac(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists("azurerm_mssql_database.test"),
				),
			},
		},
	})
}

func testAccAzureRMMsSqlDatabase_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_mssql_database" "test" {
  name                             = "acctestdb%d"
  resource_group_name              = "${azurerm_resource_group.test.name}"
  server_name                      = "${azurerm_sql_server.test.name}"
  location                         = "${azurerm_resource_group.test.location}"
  edition                          = "Standard"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes                   = "1073741824"
  requested_service_objective_name = "S0"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMMsSqlDatabase_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database" "import" {
  name                             = "${azurerm_mssql_database.test.name}"
  resource_group_name              = "${azurerm_mssql_database.test.resource_group_name}"
  server_name                      = "${azurerm_mssql_database.test.server_name}"
  location                         = "${azurerm_mssql_database.test.location}"
  edition                          = "${azurerm_mssql_database.test.edition}"
  collation                        = "${azurerm_mssql_database.test.collation}"
  max_size_bytes                   = "${azurerm_mssql_database.test.max_size_bytes}"
  requested_service_objective_name = "${azurerm_mssql_database.test.requested_service_objective_name}"
}
`, testAccAzureRMMsSqlDatabase_basic(rInt, location))
}

func testAccAzureRMMsSqlDatabase_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_mssql_database" "test" {
  name                             = "acctestdb%d"
  resource_group_name              = "${azurerm_resource_group.test.name}"
  server_name                      = "${azurerm_sql_server.test.name}"
  location                         = "${azurerm_resource_group.test.location}"
  edition                          = "Standard"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes                   = "1073741824"
  requested_service_objective_name = "S0"

  tags {
    environment = "staging"
    database    = "test"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMMsSqlDatabase_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_mssql_database" "test" {
  name                             = "acctestdb%d"
  resource_group_name              = "${azurerm_resource_group.test.name}"
  server_name                      = "${azurerm_sql_server.test.name}"
  location                         = "${azurerm_resource_group.test.location}"
  edition                          = "Standard"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes                   = "1073741824"
  requested_service_objective_name = "S0"

  tags {
    environment = "production"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMMsSqlDatabase_dataWarehouse(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest_rg_%d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_mssql_database" "test" {
  name                             = "acctestdb%d"
  resource_group_name              = "${azurerm_resource_group.test.name}"
  server_name                      = "${azurerm_sql_server.test.name}"
  location                         = "${azurerm_resource_group.test.location}"
  edition                          = "DataWarehouse"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  requested_service_objective_name = "DW400"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMMsSqlDatabase_restorePointInTime(rInt int, formattedTime string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_mssql_database" "test" {
  name                             = "acctestdb%d"
  resource_group_name              = "${azurerm_resource_group.test.name}"
  server_name                      = "${azurerm_sql_server.test.name}"
  location                         = "${azurerm_resource_group.test.location}"
  edition                          = "Standard"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes                   = "1073741824"
  requested_service_objective_name = "S0"
}

resource "azurerm_mssql_database" "test_restore" {
  name                  = "acctestdb_restore%d"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  server_name           = "${azurerm_sql_server.test.name}"
  location              = "${azurerm_resource_group.test.location}"
  create_mode           = "PointInTimeRestore"
  source_database_id    = "${azurerm_mssql_database.test.id}"
  restore_point_in_time = "%s"
}
`, rInt, location, rInt, rInt, rInt, formattedTime)
}

func testAccAzureRMMsSqlDatabase_elasticPool(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_elasticpool" "test" {
  name                = "acctestep%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  server_name         = "${azurerm_sql_server.test.name}"
  edition             = "Basic"
  dtu                 = 50
  pool_size           = 5000
}

resource "azurerm_mssql_database" "test" {
  name                             = "acctestdb%d"
  resource_group_name              = "${azurerm_resource_group.test.name}"
  server_name                      = "${azurerm_sql_server.test.name}"
  location                         = "${azurerm_resource_group.test.location}"
  edition                          = "${azurerm_sql_elasticpool.test.edition}"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes                   = "1073741824"
  elastic_pool_name                = "${azurerm_sql_elasticpool.test.name}"
  requested_service_objective_name = "ElasticPool"
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMMsSqlDatabase_collationUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_mssql_database" "test" {
  name                             = "acctestdb%d"
  resource_group_name              = "${azurerm_resource_group.test.name}"
  server_name                      = "${azurerm_sql_server.test.name}"
  location                         = "${azurerm_resource_group.test.location}"
  edition                          = "Standard"
  collation                        = "Japanese_Bushu_Kakusu_100_CS_AS_KS_WS"
  max_size_bytes                   = "1073741824"
  requested_service_objective_name = "S0"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMMsSqlDatabase_bacpac(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG_%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "bacpac"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_storage_blob" "test" {
  name                   = "test.bacpac"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"
  type                   = "block"
  source                 = "testdata/sql_import.bacpac"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_sql_firewall_rule" "test" {
  name                = "allowazure"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_sql_server.test.name}"
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "0.0.0.0"
}

resource "azurerm_mssql_database" "test" {
  name                             = "acctestdb%d"
  resource_group_name              = "${azurerm_resource_group.test.name}"
  server_name                      = "${azurerm_sql_server.test.name}"
  location                         = "${azurerm_resource_group.test.location}"
  edition                          = "Standard"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes                   = "1073741824"
  requested_service_objective_name = "S0"

  import {
    storage_uri                  = "${azurerm_storage_blob.test.url}"
    storage_key                  = "${azurerm_storage_account.test.primary_access_key}"
    storage_key_type             = "StorageAccessKey"
    administrator_login          = "${azurerm_sql_server.test.administrator_login}"
    administrator_login_password = "${azurerm_sql_server.test.administrator_login_password}"
    authentication_type          = "SQL"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMMsSqlDatabase_requestedServiceObjectiveName(rInt int, location, requestedServiceObjectiveName string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_mssql_database" "test" {
  name                             = "acctestdb%d"
  resource_group_name              = "${azurerm_resource_group.test.name}"
  server_name                      = "${azurerm_sql_server.test.name}"
  location                         = "${azurerm_resource_group.test.location}"
  edition                          = "Standard"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes                   = "1073741824"
  requested_service_objective_name = %q
}
`, rInt, location, rInt, rInt, requestedServiceObjectiveName)
}

func testAccAzureRMMsSqlDatabase_threatDetectionPolicy(rInt int, location, state string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "test%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctestsqlserver%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  location                     = "${azurerm_resource_group.test.location}"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_mssql_database" "test" {
  name                = "acctestdb%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_sql_server.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  edition             = "Standard"
  collation           = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes      = "1073741824"

  threat_detection_policy {
    retention_days             = 15
    state                      = "%s"
    disabled_alerts            = ["Sql_Injection"]
    email_account_admins       = "Enabled"
    storage_account_access_key = "${azurerm_storage_account.test.primary_access_key}"
    storage_endpoint           = "${azurerm_storage_account.test.primary_blob_endpoint}"
    use_server_default         = "Disabled"
  }
}
`, rInt, location, rInt, rInt, rInt, state)
}
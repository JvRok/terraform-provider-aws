package ec2_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func TestAccEC2SerialConsoleAccess_basic(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_ec2_serial_console_access.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSerialConsoleAccessDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSerialConsoleAccessConfig_basic(false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSerialConsoleAccess(ctx, resourceName, false),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSerialConsoleAccessConfig_basic(true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSerialConsoleAccess(ctx, resourceName, true),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func testAccCheckSerialConsoleAccessDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn(ctx)

		response, err := conn.GetSerialConsoleAccessStatusWithContext(ctx, &ec2.GetSerialConsoleAccessStatusInput{})
		if err != nil {
			return err
		}

		if aws.BoolValue(response.SerialConsoleAccessEnabled) != false {
			return fmt.Errorf("Serial console access not disabled on resource removal")
		}

		return nil
	}
}

func testAccCheckSerialConsoleAccess(ctx context.Context, n string, enabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn(ctx)

		response, err := conn.GetSerialConsoleAccessStatusWithContext(ctx, &ec2.GetSerialConsoleAccessStatusInput{})
		if err != nil {
			return err
		}

		if aws.BoolValue(response.SerialConsoleAccessEnabled) != enabled {
			return fmt.Errorf("Serial console access is not in expected state (%t)", enabled)
		}

		return nil
	}
}

func testAccSerialConsoleAccessConfig_basic(enabled bool) string {
	return fmt.Sprintf(`
resource "aws_ec2_serial_console_access" "test" {
  enabled = %[1]t
}
`, enabled)
}

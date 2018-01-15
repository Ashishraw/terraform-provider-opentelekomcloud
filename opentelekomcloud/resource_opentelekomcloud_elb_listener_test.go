package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/listeners"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// PASS
func TestAccELBListener_basic(t *testing.T) {
	var listener listeners.Listener

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckELBListenerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccELBListenerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBListenerExists("opentelekomcloud_elb_listener.listener_1", &listener),
				),
			},
			resource.TestStep{
				Config: TestAccELBListenerConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"opentelekomcloud_elb_listener.listener_1", "name", "listener_1_updated"),
				),
			},
		},
	})
}

func testAccCheckELBListenerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.otcV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_lb_listener_v2" {
			continue
		}

		_, err := listeners.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Listener still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckELBListenerExists(n string, listener *listeners.Listener) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		client, err := config.otcV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
		}

		found, err := listeners.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Member not found")
		}

		*listener = *found

		return nil
	}
}

var TestAccELBListenerConfig_basic = fmt.Sprintf(`
resource "opentelekomcloud_vpc_v1" "vpc_1" {
	name = "terraform-testacc-vpc-elb"
	cidr= "192.168.0.0/16"
}

resource "opentelekomcloud_elb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vpc_id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
  type = "External"
  bandwidth = 5
}

resource "opentelekomcloud_elb_listener" "listener_1" {
  name = "listener_1"
  protocol = "TCP"
  protocol_port = 8080
  backend_protocol = "TCP"
  backend_port = 8080
  lb_algorithm = "roundrobin"
  loadbalancer_id = "${opentelekomcloud_elb_loadbalancer.loadbalancer_1.id}"

	timeouts {
		create = "5m"
		update = "5m"
		delete = "5m"
	}
}
`)

var TestAccELBListenerConfig_update = fmt.Sprintf(`
resource "opentelekomcloud_vpc_v1" "vpc_1" {
	name = "terraform-testacc-vpc-elb"
	cidr= "192.168.0.0/16"
}

resource "opentelekomcloud_elb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vpc_id = "${opentelekomcloud_vpc_v1.vpc_1.id}"
  type = "External"
  bandwidth = 5
}

resource "opentelekomcloud_elb_listener" "listener_1" {
  name = "listener_1_updated"
  protocol = "TCP"
  protocol_port = 8080
  backend_protocol = "TCP"
  backend_port = 8080
  lb_algorithm = "roundrobin"
  loadbalancer_id = "${opentelekomcloud_elb_loadbalancer.loadbalancer_1.id}"

	timeouts {
		create = "5m"
		update = "5m"
		delete = "5m"
	}
}
`)

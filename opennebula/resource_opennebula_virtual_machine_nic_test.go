package opennebula

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccVirtualMachineNICUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachineNIC,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "name", "test-virtual_machine"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "nic.#", "1"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "nic.0.ip", "172.16.100.131"),
				),
			},
			{
				Config: testAccVirtualMachineNICAttach,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "name", "test-virtual_machine"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "nic.#", "1"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "nic.0.ip", "172.16.100.131"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine_nic.test", "ip", "172.16.100.111"),
				),
			},
			{
				Config: testAccVirtualMachineNICUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "name", "test-virtual_machine"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "nic.#", "1"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "nic.0.ip", "172.16.100.131"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine_nic.test", "ip", "172.16.100.151"),
				),
			},
			{
				Config: testAccVirtualMachineNICDetach,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "name", "test-virtual_machine"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "nic.#", "1"),
					resource.TestCheckResourceAttr("opennebula_virtual_machine.test", "nic.0.ip", "172.16.100.131"),
				),
			},
		},
	})
}

var testAccVirtualMachineNIC = `

resource "opennebula_virtual_network" "net1" {
	name = "test-net1"
	type            = "dummy"
	bridge          = "onebr"
	mtu             = 1500
	ar {
	  ar_type = "IP4"
	  size    = 12
	  ip4     = "172.16.100.130"
	}
	permissions = "642"
	group = "oneadmin"
	security_groups = [0]
	clusters = [0]
}

resource "opennebula_virtual_network" "net2" {
	name = "test-net2"
	type            = "dummy"
	bridge          = "onebr"
	mtu             = 1500
	ar {
	  ar_type = "IP4"
	  size    = 16
	  ip4     = "172.16.100.110"
	}
	permissions = "642"
	group = "oneadmin"
	security_groups = [0]
	clusters = [0]
}

resource "opennebula_virtual_network" "net3" {
	name = "test-net3"
	type            = "dummy"
	bridge          = "onebr"
	mtu             = 1500
	ar {
	  ar_type = "IP4"
	  size    = 16
	  ip4     = "172.16.100.150"
	}
	permissions = "642"
	group = "oneadmin"
	security_groups = [0]
	clusters = [0]
}

resource "opennebula_virtual_machine" "test" {
	name        = "test-virtual_machine"
	group       = "oneadmin"
	permissions = "642"
	memory = 128
	cpu = 0.1
  
	context = {
	  NETWORK  = "YES"
	  SET_HOSTNAME = "$NAME"
	}
  
	graphics {
	  type   = "VNC"
	  listen = "0.0.0.0"
	  keymap = "en-us"
	}
  
	os {
	  arch = "x86_64"
	  boot = ""
	}
  
	nic {
		network_id = opennebula_virtual_network.net1.id
		ip = "172.16.100.131"
	}
  
	timeout = 5
}
`

var testAccVirtualMachineNICAttach = `

resource "opennebula_virtual_network" "net1" {
	name = "test-net1"
	type            = "dummy"
	bridge          = "onebr"
	mtu             = 1500
	ar {
	  ar_type = "IP4"
	  size    = 12
	  ip4     = "172.16.100.130"
	}
	permissions = "642"
	group = "oneadmin"
	security_groups = [0]
	clusters = [0]
}

resource "opennebula_virtual_network" "net2" {
	name = "test-net2"
	type            = "dummy"
	bridge          = "onebr"
	mtu             = 1500
	ar {
	  ar_type = "IP4"
	  size    = 16
	  ip4     = "172.16.100.110"
	}
	permissions = "642"
	group = "oneadmin"
	security_groups = [0]
	clusters = [0]
}

resource "opennebula_virtual_network" "net3" {
	name = "test-net3"
	type            = "dummy"
	bridge          = "onebr"
	mtu             = 1500
	ar {
	  ar_type = "IP4"
	  size    = 16
	  ip4     = "172.16.100.150"
	}
	permissions = "642"
	group = "oneadmin"
	security_groups = [0]
	clusters = [0]
}

resource "opennebula_virtual_machine_nic" "test" {
	vm_id      = opennebula_virtual_machine.test.id
	network_id = opennebula_virtual_network.net2.id
	ip         = "172.16.100.111"	
}

resource "opennebula_virtual_machine" "test" {
	name        = "test-virtual_machine"
	group       = "oneadmin"
	permissions = "642"
	memory = 128
	cpu = 0.1
  
	context = {
	  NETWORK  = "YES"
	  SET_HOSTNAME = "$NAME"
	}
  
	graphics {
	  type   = "VNC"
	  listen = "0.0.0.0"
	  keymap = "en-us"
	}
  
	os {
	  arch = "x86_64"
	  boot = ""
	}

	nic {
		network_id = opennebula_virtual_network.net1.id
		ip = "172.16.100.131"
	}
  
	timeout = 5
}
`

var testAccVirtualMachineNICUpdate = `

resource "opennebula_virtual_network" "net1" {
	name = "test-net1"
	type            = "dummy"
	bridge          = "onebr"
	mtu             = 1500
	ar {
	  ar_type = "IP4"
	  size    = 12
	  ip4     = "172.16.100.130"
	}
	permissions = "642"
	group = "oneadmin"
	security_groups = [0]
	clusters = [0]
}

resource "opennebula_virtual_network" "net2" {
	name = "test-net2"
	type            = "dummy"
	bridge          = "onebr"
	mtu             = 1500
	ar {
	  ar_type = "IP4"
	  size    = 16
	  ip4     = "172.16.100.110"
	}
	permissions = "642"
	group = "oneadmin"
	security_groups = [0]
	clusters = [0]

}

resource "opennebula_virtual_network" "net3" {
	name = "test-net3"
	type            = "dummy"
	bridge          = "onebr"
	mtu             = 1500
	ar {
	  ar_type = "IP4"
	  size    = 16
	  ip4     = "172.16.100.150"
	}
	permissions = "642"
	group = "oneadmin"
	security_groups = [0]
	clusters = [0]
}

resource "opennebula_virtual_machine_nic" "test" {
	vm_id      = opennebula_virtual_machine.test.id
	network_id = opennebula_virtual_network.net3.id
	ip         = "172.16.100.151"	
}

  resource "opennebula_virtual_machine" "test" {
	  name        = "test-virtual_machine"
	  group       = "oneadmin"
	  permissions = "642"
	  memory = 128
	  cpu = 0.1
	
	  context = {
		NETWORK  = "YES"
		SET_HOSTNAME = "$NAME"
	  }
	
	  graphics {
		type   = "VNC"
		listen = "0.0.0.0"
		keymap = "en-us"
	  }
	
	  os {
		arch = "x86_64"
		boot = ""
	  }

	  nic {
		network_id = opennebula_virtual_network.net1.id
		ip = "172.16.100.131"
	  }
	
	  timeout = 5
}
`

var testAccVirtualMachineNICDetach = `

resource "opennebula_virtual_network" "net1" {
	name = "test-net1"
	type            = "dummy"
	bridge          = "onebr"
	mtu             = 1500
	ar {
	  ar_type = "IP4"
	  size    = 12
	  ip4     = "172.16.100.130"
	}
	permissions = "642"
	group = "oneadmin"
	security_groups = [0]
	clusters = [0]
}

resource "opennebula_virtual_network" "net2" {
	name = "test-net2"
	type            = "dummy"
	bridge          = "onebr"
	mtu             = 1500
	ar {
	  ar_type = "IP4"
	  size    = 16
	  ip4     = "172.16.100.110"
	}
	permissions = "642"
	group = "oneadmin"
	security_groups = [0]
	clusters = [0]
	
}

resource "opennebula_virtual_network" "net3" {
	name = "test-net3"
	type            = "dummy"
	bridge          = "onebr"
	mtu             = 1500
	ar {
	  ar_type = "IP4"
	  size    = 16
	  ip4     = "172.16.100.150"
	}
	permissions = "642"
	group = "oneadmin"
	security_groups = [0]
	clusters = [0]
}

resource "opennebula_virtual_machine" "test" {
	  name        = "test-virtual_machine"
	  group       = "oneadmin"
	  permissions = "642"
	  memory = 128
	  cpu = 0.1
	
	  context = {
		NETWORK  = "YES"
		SET_HOSTNAME = "$NAME"
	  }
	
	  graphics {
		type   = "VNC"
		listen = "0.0.0.0"
		keymap = "en-us"
	  }
	
	  os {
		arch = "x86_64"
		boot = ""
	  }

	  nic {
		network_id = opennebula_virtual_network.net1.id
		ip = "172.16.100.131"
	  }
	
	  timeout = 5
}
`

package opennebula

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/OpenNebula/one/src/oca/go/src/goca"
	"github.com/OpenNebula/one/src/oca/go/src/goca/schemas/shared"
)

var (
	vmDiskUpdateReadyStates = []string{"RUNNING", "POWEROFF"}
	vmNICUpdateReadyStates  = vmDiskUpdateReadyStates
)

func resourceOpennebulaVirtualMachineNIC() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpennebulaVirtualMachineNICCreate,
		Read:   resourceOpennebulaVirtualMachineNICRead,
		Exists: resourceOpennebulaVirtualMachineNICExists,
		Delete: resourceOpennebulaVirtualMachineNICDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vm_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"network_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"network": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"mac": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"model": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"physical_device": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"security_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceOpennebulaVirtualMachineNICCreate(d *schema.ResourceData, meta interface{}) error {
	controller := meta.(*goca.Controller)

	nicTpl := shared.NewNIC()
	vmID := d.Get("vm_id").(int)
	vnetID := d.Get("network_id").(int)

	nicTpl.Add(shared.NetworkID, vnetID)

	if v, ok := d.GetOk("ip"); ok {
		nicTpl.Add(shared.IP, v.(string))
	}
	if v, ok := d.GetOk("mac"); ok {
		nicTpl.Add(shared.MAC, v.(string))
	}
	if v, ok := d.GetOk("model"); ok {
		nicTpl.Add(shared.Model, v.(string))
	}
	if v, ok := d.GetOk("physical_device"); ok {
		nicTpl.Add("PHYDEV", v.(string))
	}
	if v, ok := d.GetOk("security_groups"); ok {
		secGroups := ArrayToString(v.([]interface{}), ",")
		nicTpl.Add(shared.SecurityGroups, secGroups)
	}

	// wait before checking NIC
	// TODO: remove fixed value timeout
	vmc := controller.VM(vmID)
	nicID, err := vmNICAttach(vmc, 30, nicTpl)
	if err != nil {
		return fmt.Errorf("VM NIC attach: %s", err)
	}

	d.SetId(fmt.Sprintf("%d", nicID))

	log.Printf("[INFO] Successfully attached VM NIC\n")

	return resourceOpennebulaVirtualMachineNICRead(d, meta)
}

func resourceOpennebulaVirtualMachineNICRead(d *schema.ResourceData, meta interface{}) error {

	vmc, err := getVirtualMachineController(d, meta, -2, -1, -1)
	if err != nil {
		if NoExists(err) {
			log.Printf("[WARN] Removing VM NIC %s from state because it no longer exists in", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	// TODO: fix it after 5.10 release
	// Force the "decrypt" bool to false to keep ONE 5.8 behavior
	vm, err := vmc.Info(false)
	if err != nil {
		return err
	}

	// get the nic ID from the nic list
	var nic *shared.NIC

	nics := vm.Template.GetNICs()
	for _, n := range nics {
		nicID, _ := n.Get(shared.NICID)
		if nicID == d.Id() {
			nic = &n
			break
		}
	}

	if nic == nil {
		return fmt.Errorf("VM NIC (ID:%s) not found", d.Id())
	}

	networkID, _ := nic.Get(shared.NetworkID)
	ip, _ := nic.Get(shared.IP)
	mac, _ := nic.Get(shared.MAC)
	phyDev, _ := nic.GetStr("PHYDEV")
	network, _ := nic.Get(shared.Network)
	model, _ := nic.Get(shared.Model)

	sg := make([]int, 0)
	securityGroupsArray, _ := nic.Get(shared.SecurityGroups)
	sgString := strings.Split(securityGroupsArray, ",")
	for _, s := range sgString {
		sgInt, _ := strconv.ParseInt(s, 10, 32)
		sg = append(sg, int(sgInt))
	}

	d.Set("network_id", networkID)
	d.Set("vm_id", vm.ID)
	d.Set("ip", ip)
	d.Set("mac", mac)
	d.Set("physical_device", phyDev)
	d.Set("network", network)
	d.Set("model", model)
	d.Set("security_groups", sg)

	return nil
}

func resourceOpennebulaVirtualMachineNICExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	err := resourceOpennebulaVirtualMachineNICRead(d, meta)
	if err != nil || d.Id() == "" {
		return false, err
	}

	return true, nil
}

func resourceOpennebulaVirtualMachineNICDelete(d *schema.ResourceData, meta interface{}) error {
	err := resourceOpennebulaVirtualMachineNICRead(d, meta)
	if err != nil || d.Id() == "" {
		return err
	}

	//Get VM
	vmc, err := getVirtualMachineController(d, meta)
	if err != nil {
		return err
	}

	nicID, err := strconv.ParseInt(d.Id(), 10, 0)
	if err != nil {
		return fmt.Errorf("VM NIC can't parse the ID: %s", err)
	}

	// wait before checking NIC
	// TODO: remove fixed value timeout
	err = vmNICDetach(vmc, 30, int(nicID))
	if err != nil {
		return fmt.Errorf("VM NIC detach: %s", err)
	}

	log.Printf("[INFO] Successfully detached VM NIC\n")
	return nil
}

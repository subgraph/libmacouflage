package libmacouflage

import (
	"testing"
	"fmt"
	"os"
	"strings"
)

var testInterface = os.Getenv("TEST_INTERFACE")

func GetTestInterface() (iface string) {
	if testInterface != "" {
		iface = testInterface
	} else {
		ifaces, _ := GetInterfaces()
		iface = ifaces[0].Name
	}
	return
}

func Test_GetCurrentMac_1(t *testing.T) {
	_, err := GetCurrentMac("badinterface")
	if err == nil {
		t.Errorf("GetCurrentMac_1 error for non-existent interface name: %s\n",
		err)
	} else {
		fmt.Printf("GetCurrentMac_1 result for non-existent interface name: %s\n",
		err)
	}
}

func Test_GetCurrentMac_2(t *testing.T) {
	iface := GetTestInterface()
	_, err := GetCurrentMac(iface)
	if err != nil {
		t.Errorf("GetCurrentMac_2 for existing interface name %s failed: %s\n",
		iface, err)
	}
}

func Test_GetAllCurrentMacs_1(t *testing.T) {
	ifaces, err := GetInterfaces()
	macs, err := GetAllCurrentMacs()
	if err != nil {
		t.Errorf("GetAllCurrentMacs_1 failed: %s\n", err)
	}
	if len(macs) != len(ifaces) {
		t.Errorf("GetAllCurrentMacs_1 failed with length mismatch\n")
	}
}

func Test_GetPermanentMac_1(t *testing.T) {
	_, err := GetPermanentMac("badinterface")
	if err == nil {
		t.Errorf("GetPermanentMac_1 for non-existent interface name: %s\n",
		err)
	} else {
		fmt.Printf("GetPermanentMac_1 for non-existent interface name: %s\n", 
		err)
	}
}

func Test_GetPermanentMac_2(t *testing.T) {
	iface := GetTestInterface()
	_, err := GetPermanentMac(iface)
	if err != nil {
		t.Errorf("GetPermanentMac_2 for existing interface name %s failed: %s\n",
		iface, err)
	}
}

func Test_MacChanged_1(t *testing.T) {
	_, err := MacChanged("badinterface")
	if err == nil {
		t.Errorf("MacChanged_1 for non-existent interface failed to generate error\n")
	} else {
		fmt.Printf("MacChanged_1 for non-existent interface name: %s\n",
		err)
	}
}

func Test_MacChanged_2(t *testing.T) {
	iface := GetTestInterface()
	result, err := MacChanged(iface)
	if err != nil {
		t.Errorf("MacChanged_2 for existing interface name %s failed: %s\n",
		iface, err)
	}
	fmt.Printf("MacChanged_2 result: %s has changed = %t\n", iface,
	result)
}

func Test_IsIfUp_1(t *testing.T) {
	_, err := IsIfUp("badinterface")
	if err == nil {
		t.Errorf("IsIfUp_1 for non-existent interface failed to generate error\n")
	} else {
		fmt.Printf("IsIfUp_1 for non-existent interface name: %s\n", err)
	}
}

func Test_IsIfUp_2(t *testing.T) {
	iface := GetTestInterface()
	result, err := IsIfUp(iface)
	if err != nil {
		t.Errorf("IsIfUp_2 for existing interface name %s failed: %s\n",
		iface, err)
		return
	}
	fmt.Printf("IsIfUp_2 result: %s is up = %t\n", iface, result)
}

func Test_RandomizeMac_1(t *testing.T) {
	bytes := []byte{0,0,0,0,0,0}
	_, err := RandomizeMac(bytes, 2, false)
	if err == nil {
		t.Error("RandomizeMac_1 failed to generate error with bad start index")
	}
}

func Test_RandomizeMac_2(t *testing.T) {
	bytes := []byte{0,0,0,0,0,0,0}
	_, err := RandomizeMac(bytes, 0, false)
	if err == nil {
		t.Error("RandomizeMac_2 failed to generate error with macbytes byte array size")
	}
}

func Test_RunningAsRoot_1(t *testing.T) {
	result, err := RunningAsRoot()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("RunningAsRoot_1 result: user is root = %t\n", result)
}

func Test_SetMac_1(t *testing.T) {
	iface := GetTestInterface()
	oldMac, err := GetCurrentMac(iface)
	if err != nil {
		t.Error(err)
		return
	}
	bytes := []byte{0, 0, 0, 0, 0, 0}
	newMac, err := RandomizeMac(bytes, 0, false)
	if err != nil {
		t.Error(err)
		return
	}
	err = SetMac(iface, newMac.String())
	if err != nil {
		t.Error(err)
		return
	}
	same := CompareMacs(oldMac, newMac)
	if !same {
		fmt.Printf("SetMac_1 result: MAC address successfully changed for %s\n", iface)
	} else {
		t.Errorf("SetMac_1 error: MAC address not changed for %s\n", iface)
	}
}


func Test_SetMac_2(t *testing.T) {
	iface := GetTestInterface()
	oldMac, err := GetCurrentMac(iface)
	if err != nil {
		t.Error(err)
		return
	}
	bytes := []byte{0, 0, 0, 0, 0, 0}
	newMac, err := RandomizeMac(bytes, 0, true)
	if err != nil {
		t.Error(err)
		return
	}
	err = SetMac(iface, newMac.String())
	if err != nil {
		t.Error(err)
		return
	}
	same := CompareMacs(oldMac, newMac)
	if !same {
		fmt.Printf("SetMac_2 result: MAC address successfully changed for %s\n", iface)
	} else {
		t.Errorf("SetMac_2 error: MAC address not changed for %s\n", iface)
	}
}

func Test_SpoofMacRandom_1(t *testing.T) {
	iface := GetTestInterface()
	oldMac, err := GetCurrentMac(iface)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = SpoofMacRandom(iface, true)
	if err != nil {
		t.Error(err)
		return
	}
	newMac, err := GetCurrentMac(iface)
	if err != nil {
		t.Error(err)
		return
	}
	same := CompareMacs(oldMac, newMac)
	if !same {
		fmt.Printf("SpoofMacRandom_1 result: MAC address successfully changed for %s\n", iface)
	} else {
		t.Errorf("SpoofMacRandom_1 error: MAC address not changed for %s\n", iface)
	}
}

func Test_SpoofMacSameVendor_1(t *testing.T) {
	iface := GetTestInterface()
	oldMac, err := GetCurrentMac(iface)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = SpoofMacSameVendor(iface, true)
	if err != nil {
		t.Error(err)
		return
	}
	newMac, err := GetCurrentMac(iface)
	if err != nil {
		t.Error(err)
		return
	}
	same := CompareMacs(oldMac, newMac)
	if !same {
		fmt.Printf("SpoofMacSameVendor_1 result: MAC address successfully changed for %s\n", iface)
	} else {
		t.Errorf("SpoofMacSameVendor_1 error: MAC address not changed for %s\n", iface)
	}
}

func Test_SpoofMacSameDeviceType_1(t *testing.T) {
	iface := GetTestInterface()
	oldMac, err := GetCurrentMac(iface)
	if err != nil {
		t.Error(err)
		return
	}
	oldDeviceType, err := FindDeviceTypeByMac(oldMac.String())
	if err != nil {
		t.Error(err)
		return
	}
	_, err = SpoofMacSameDeviceType(iface)
	if err != nil {
		t.Error(err)
		return
	}
	newMac, err := GetCurrentMac(iface)
	if err != nil {
		t.Error(err)
		return
	}
	newDeviceType, err := FindDeviceTypeByMac(newMac.String())
	if err != nil {
		t.Error(err)
		return
	}
	if !(strings.EqualFold(oldDeviceType, newDeviceType)) {
		t.Errorf("SpoofMacSameDeviceType_1 error - device mismatch between old (%s) and new (%s)",
		oldDeviceType, newDeviceType)
	}
}


func Test_SpoofMacAnyDeviceType_1(t *testing.T) {
	iface := GetTestInterface()
	_, err := SpoofMacAnyDeviceType(iface)
	if err != nil {
		t.Error(err)
		return
	}
	newMac, err := GetCurrentMac(iface)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = FindVendorByMac(newMac.String())
	if err != nil {
		t.Errorf("SpoofMacAnyDeviceType_1 error - cannot find vendor for new MAC: %s",
			newMac)
	}

}

func Test_SpoofMacPopular_1(t *testing.T) {
	iface := GetTestInterface()
	_, err := SpoofMacPopular(iface)
	if err != nil {
		t.Error(err)
		return
	}
	newMac, err := GetCurrentMac(iface)
	if err != nil {
		t.Error(err)
		return
	}
	vendor, err := FindVendorByMac(newMac.String())
	if err != nil {
		t.Error(err)
		return
	}
	if !vendor.Popular {
		t.Errorf("SpoofMacPopular_1 error - vendor selected for new MAC is not in popular list: %s",
		newMac)
	}
}

func Test_RevertMac_1(t *testing.T) {
	iface := GetTestInterface()
	changed, _ := MacChanged(iface)
	if !changed {
		bytes := []byte{0, 0, 0, 0, 0, 0}
		mac, err := RandomizeMac(bytes, 0, true)
		if err != nil {
			t.Error(err)
			return
		}
		err = SetMac(iface, mac.String())
		if err != nil {
			t.Error(err)
			return
		}
	}
	err := RevertMac(iface)
	if err != nil {
		t.Error(err)
		return
	}
	changed, err = MacChanged(iface)
	if !changed {
		fmt.Printf("RevertMac_1 result: MAC address successfully reverted for %s\n", iface)
	} else {
		t.Errorf("RevertMac_1 error: MAC address not changed for %s\n", iface)
	}
}

func Test_FindAllPopularOuis_1(t *testing.T) {
	popular, err := FindAllPopularOuis()
	if err != nil {
		t.Error(err)
		return
	}
	for _, oui := range popular {
		if(!oui.Popular) {
			t.Errorf("FindAllPopularOuis_1 error, found erroneous oui: ", oui)
			return
		}
	}
}

func Test_FindDeviceTypeByMac_1(t *testing.T) {
	iface := GetTestInterface()
	mac, err := GetCurrentMac(iface)
	if err != nil {
		t.Error(err)
		return
	}
	deviceType, err := FindDeviceTypeByMac(mac.String())
	if err != nil {
		t.Errorf("FindDeviceTypeByMac_1 error for device type %s: %s",
			deviceType, err)
		return
	}
}

func Test_FindAllVendorsByDeviceType_1(t *testing.T) {
	iface := GetTestInterface()
	mac, err := GetCurrentMac(iface)
	if err != nil {
		t.Error(err)
		return
	}
	deviceType, err := FindDeviceTypeByMac(mac.String())
	if err != nil {
		t.Error(err)
		return
	}
	vendors, err := FindAllVendorsByDeviceType(deviceType)
	for _, vendor := range vendors {
		if !(strings.EqualFold(vendor.Devices[0].DeviceType, deviceType)) {
			t.Errorf("FindAllVendorByDeviceType_1 error with %s and device type: %s \n",
			vendor.Vendor, deviceType)
			return
		}
	}
}

func Test_FindVendorByMac_1(t *testing.T) {
	// Test against locally administered address, will not appear in OuiDb
	mac := "06:00:00:00:00:00"
	_, err := FindVendorByMac(mac)
	if err == nil {
		t.Errorf("FindVendorByMac_1 failed for locally administered address: ",
		mac)
	}
}

func Test_FindVendorByMac_2(t *testing.T) {
	mac := "00:00:00:00:00:00"
	_, err := FindVendorByMac(mac)
	if err != nil {
		t.Errorf("FindVendorByMac_2 failed for valid vendor MAC: %s", mac)
	}
}

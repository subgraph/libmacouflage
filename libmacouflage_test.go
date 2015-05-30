package libmacouflage

import (
	"testing"
	"fmt"
)

func Test_GetCurrentMac_1(t *testing.T) {
	_, err := GetCurrentMac("badinterface")
	if err == nil {
		t.Errorf("GetCurrentMac_1 for non-existent interface name: %s\n",
		err)
	} else {
		fmt.Printf("GetCurrentMac_1 for non-existent interface name: %s\n",
		err)
	}
}

func Test_GetCurrentMac_2(t *testing.T) {
	ifaces, err := GetInterfaces()
	_, err = GetCurrentMac(ifaces[0].Name)
	if err != nil {
		t.Errorf("GetCurrentMac_2 for existing interface name %s failed: %s\n",
		ifaces[0].Name, err)
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
	ifaces, err := GetInterfaces()
	_, err = GetPermanentMac(ifaces[0].Name)
	if err != nil {
		t.Errorf("GetPermanentMac_2 for existing interface name %s failed: %s\n",
		ifaces[0].Name, err)
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
	ifaces, err := GetInterfaces()
	result, err := MacChanged(ifaces[0].Name)
	if err != nil {
		t.Errorf("MacChanged_2 for existing interface name %s failed: %s\n",
		ifaces[0].Name, err)
	}
	fmt.Printf("MacChanged_2 result: %s has changed = %t\n", ifaces[0].Name,
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
	ifaces, err := GetInterfaces()
	result, err := IsIfUp(ifaces[0].Name)
	if err != nil {
		t.Errorf("IsIfUp_2 for existing interface name %s failed: %s\n",
		ifaces[0].Name, err)
	}
	fmt.Printf("IsIfUp_2 result: %s is up = %t\n", ifaces[0].Name, result)
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
	bytes := []byte{0, 0, 0, 0, 0, 0}
	mac, err := RandomizeMac(bytes, 0, false)
	if err != nil {
		t.Error(err)
	}
	ifaces, err := GetInterfaces()
	err = SetMac(ifaces[0].Name, mac.String())
	if err != nil {
		t.Error(err)
	}
	changed, err := MacChanged(ifaces[0].Name)
	if err != nil {
		t.Error(err)
	}
	if changed {
		fmt.Printf("SetMac_1 result: MAC address successfully changed for %s\n", ifaces[0].Name)
	} else {
		t.Errorf("SetMac_1 error: MAC address not changed for %s\n", ifaces[0].Name)
	}

}

func Test_SetMac_2(t *testing.T) {
	bytes := []byte{0, 0, 0, 0, 0, 0}
	mac, err := RandomizeMac(bytes, 0, true)
	if err != nil {
		t.Error(err)
	}
	ifaces, err := GetInterfaces()
	err = SetMac(ifaces[0].Name, mac.String())
	if err != nil {
		t.Error(err)
	}
	changed, err := MacChanged(ifaces[0].Name)
	if changed {
		fmt.Printf("SetMac_2 result: MAC address successfully changed for %s\n", ifaces[0].Name)
	} else {
		t.Errorf("SetMac_2 error: MAC address not changed for %s\n", ifaces[0].Name)
	}
}

func Test_SpoofMacRandom_1(t *testing.T) {
	ifaces, err := GetInterfaces()
	changed, err := SpoofMacRandom(ifaces[0].Name, true)
	if err != nil {
		t.Error(err)
	}
	if changed {
		fmt.Printf("SpoofMacRandom_1 result: MAC address successfully changed for %s\n", ifaces[0].Name)
	} else {
		t.Errorf("SpoofMacRandom_1 error: MAC address not changed for %s\n", ifaces[0].Name)
	}
}
func Test_RevertMac_1(t *testing.T) {
	ifaces, err := GetInterfaces()
	changed, _ := MacChanged(ifaces[0].Name)
	if !changed {
		bytes := []byte{0, 0, 0, 0, 0, 0}
		mac, err := RandomizeMac(bytes, 0, true)
		if err != nil {
			t.Error(err)
		}
		err = SetMac(ifaces[0].Name, mac.String())
		if err != nil {
			t.Error(err)
		}
	}
	err = RevertMac(ifaces[0].Name)
	if err != nil {
		t.Error(err)
	}
	changed, err = MacChanged(ifaces[0].Name)
	if !changed {
		fmt.Printf("RevertMac_1 result: MAC address successfully reverted for %s\n", ifaces[0].Name)
	} else {
		t.Errorf("RevertMac_1 error: MAC address not changed for %s\n", ifaces[0].Name)
	}
}


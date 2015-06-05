package libmacouflage

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
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
	assert.Error(t, err, "Function failed to generate error for bad interface")
}

func Test_GetCurrentMac_2(t *testing.T) {
	iface := GetTestInterface()
	_, err := GetCurrentMac(iface)
	assert.NoError(t, err, iface)
}

func Test_GetAllCurrentMacs_1(t *testing.T) {
	ifaces, err := GetInterfaces()
	macs, err := GetAllCurrentMacs()
	assert.NoError(t, err)
	assert.Equal(t, len(macs), len(ifaces))
}

func Test_GetPermanentMac_1(t *testing.T) {
	_, err := GetPermanentMac("badinterface")
	assert.Error(t, err, "Function failed to generate error for bad interface")
}

func Test_GetPermanentMac_2(t *testing.T) {
	iface := GetTestInterface()
	_, err := GetPermanentMac(iface)
	assert.NoError(t, err, iface)
}

func Test_MacChanged_1(t *testing.T) {
	_, err := MacChanged("badinterface")
	assert.Error(t, err, "Function failed to generate error for bad interface")
}

func Test_MacChanged_2(t *testing.T) {
	iface := GetTestInterface()
	_, err := MacChanged(iface)
	assert.NoError(t, err, iface)
}

func Test_IsIfUp_1(t *testing.T) {
	_, err := IsIfUp("badinterface")
	assert.Error(t, err, "Function failed to generate error for bad interface")
}

func Test_IsIfUp_2(t *testing.T) {
	iface := GetTestInterface()
	_, err := IsIfUp(iface)
	assert.NoError(t, err, iface)
}

func Test_RandomizeMac_1(t *testing.T) {
	bytes := []byte{0,0,0,0,0,0}
	_, err := RandomizeMac(bytes, 2, false)
	assert.Error(t, err, "Function failed to generate error for bad start index")
}

func Test_RandomizeMac_2(t *testing.T) {
	bytes := []byte{0,0,0,0,0,0,0}
	_, err := RandomizeMac(bytes, 0, false)
	assert.Error(t, err, "Function failed to generate error for long byte array")
}

func Test_RunningAsRoot_1(t *testing.T) {
	_, err := RunningAsRoot()
	assert.NoError(t, err)
}

func Test_SetMac_1(t *testing.T) {
	iface := GetTestInterface()
	bytes := []byte{0, 0, 0, 0, 0, 0}
	newMac, err := RandomizeMac(bytes, 0, false)
	assert.NoError(t, err)
	err = SetMac(iface, newMac.String())
	assert.NoError(t, err)
}


func Test_SetMac_2(t *testing.T) {
	iface := GetTestInterface()
	bytes := []byte{0, 0, 0, 0, 0, 0}
	newMac, err := RandomizeMac(bytes, 0, true)
	assert.NoError(t, err)
	err = SetMac(iface, newMac.String())
	assert.NoError(t, err)
}

func Test_SpoofMacRandom_1(t *testing.T) {
	iface := GetTestInterface()
	_, err := SpoofMacRandom(iface, true)
	assert.NoError(t, err)
}

func Test_SpoofMacSameVendor_1(t *testing.T) {
	iface := GetTestInterface()
	_, err := SpoofMacSameVendor(iface, true)
	assert.NoError(t, err)
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
	assert.NoError(t, err)
	newMac, err := GetCurrentMac(iface)
	assert.NoError(t, err)
	newDeviceType, err := FindDeviceTypeByMac(newMac.String())
	assert.NoError(t, err)
	assert.Equal(t, oldDeviceType, newDeviceType)
}


func Test_SpoofMacAnyDeviceType_1(t *testing.T) {
	iface := GetTestInterface()
	_, err := SpoofMacAnyDeviceType(iface)
	assert.NoError(t, err)
	newMac, err := GetCurrentMac(iface)
	assert.NoError(t, err)
	_, err = FindVendorByMac(newMac.String())
	assert.NoError(t, err)
}

func Test_SpoofMacPopular_1(t *testing.T) {
	iface := GetTestInterface()
	_, err := SpoofMacPopular(iface)
	assert.NoError(t, err)
	newMac, err := GetCurrentMac(iface)
	assert.NoError(t, err)
	vendor, err := FindVendorByMac(newMac.String())
	assert.NoError(t, err)
	assert.True(t, vendor.Popular, "newMac is not flagged as popular: %s", newMac.String())
}

func Test_RevertMac_1(t *testing.T) {
	iface := GetTestInterface()
	changed, _ := MacChanged(iface)
	if !changed {
		bytes := []byte{0, 0, 0, 0, 0, 0}
		mac, err := RandomizeMac(bytes, 0, true)
		assert.NoError(t, err)
		err = SetMac(iface, mac.String())
		assert.NoError(t, err, "Spoofing dummy MAC prior to reverting")
	}
	err := RevertMac(iface)
	assert.NoError(t, err, "Reverting MAC")
	changed, err = MacChanged(iface)
	assert.False(t, changed)
}

func Test_FindAllPopularOuis_1(t *testing.T) {
	popular, err := FindAllPopularOuis()
	assert.NoError(t, err)
	for _, oui := range popular {
		if(!oui.Popular) {
			assert.True(t, oui.Popular)
		}
	}
}

func Test_FindDeviceTypeByMac_1(t *testing.T) {
	iface := GetTestInterface()
	mac, err := GetCurrentMac(iface)
	assert.NoError(t, err)
	_, err = FindDeviceTypeByMac(mac.String())
	assert.NoError(t, err, mac.String())
}

func Test_FindAllVendorsByDeviceType_1(t *testing.T) {
	iface := GetTestInterface()
	mac, err := GetCurrentMac(iface)
	assert.NoError(t, err)
	deviceType, err := FindDeviceTypeByMac(mac.String())
	assert.NoError(t, err)
	vendors, err := FindAllVendorsByDeviceType(deviceType)
	for _, vendor := range vendors {
		assert.Equal(t, vendor.Devices[0].DeviceType, deviceType)
	}
}

func Test_FindVendorByMac_1(t *testing.T) {
	// Test against locally administered address, will not appear in OuiDb
	mac := "06:00:00:00:00:00"
	_, err := FindVendorByMac(mac)
	assert.Error(t, err, "Functioned failed to generate error for locally administered address")
}

func Test_FindVendorByMac_2(t *testing.T) {
	mac := "00:00:00:00:00:00"
	_, err := FindVendorByMac(mac)
	assert.NoError(t, err)
}

func Test_FindVendorsByKeyword_1(t *testing.T) {
	results, err := FindVendorsByKeyword("intel")
	assert.NoError(t,  err)
	assert.NotEqual(t, 0, len(results))
}

func Test_FindVendorsByKeyword_2(t *testing.T) {
	results, err := FindVendorsByKeyword("")
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(results))
}

func Test_FindVendorsByKeyword_3(t *testing.T) {
	results, err := FindVendorsByKeyword("salkslfkdlfkdf8dfurewkjfiew8f9ewf")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(results))
}

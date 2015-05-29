package libmacouflage

import "C"
import (
	"syscall"
	"fmt"
	"net"
	"unsafe"
	"crypto/rand"
	"os/user"
)

const SIOCSIFHWADDR = 0x8924
const SIOCETHTOOL = 0x8946
const ETHTOOL_GPERMADDR = 0x00000020
const IFHWADDRLEN = 6

// TODO: Ad-hoc structs that work, fix
type NetInfo struct {
	name [16]byte
	family uint16
	data [6]byte
}

type ifreq struct {
	name [16]byte
        epa *EthtoolPermAddr
}

type EthtoolPermAddr struct {
	cmd uint32
	size uint32
	data [6]byte
}

func GetCurrentMac(name string) (mac net.HardwareAddr, err error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return
	}
	mac = iface.HardwareAddr
	return
}

func GetAllCurrentMacs() (macs map[string]string, err error) {
	ifaces, err := GetInterfaces()
	macs = make(map[string]string)
	for _, iface := range ifaces {
		macs[iface.Name] = iface.HardwareAddr.String()
	}
	return
}

func GetInterfaces() (ifaces []net.Interface, err error) {
	allIfaces, err := net.Interfaces()
	for _, iface := range allIfaces {
		// Skip Loopback interfaces
		if iface.Flags&net.FlagLoopback == 0 {
			ifaces = append(ifaces, iface)
		}
	}
	return
}

func GetPermanentMac(name string) (mac net.HardwareAddr, err error) {
	sockfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	defer syscall.Close(sockfd)
	var ifr ifreq
	copy(ifr.name[:], []byte(name))
	var epa EthtoolPermAddr
	epa.cmd = ETHTOOL_GPERMADDR
	epa.size = IFHWADDRLEN
	ifr.epa = &epa
	_, _, errno  := syscall.Syscall(syscall.SYS_IOCTL, uintptr(sockfd), SIOCETHTOOL, uintptr(unsafe.Pointer(&ifr)))
	if errno != 0 {
		err = syscall.Errno(errno)
		return
	}
	mac = net.HardwareAddr(C.GoBytes(unsafe.Pointer(&ifr.epa.data), 6))
	return
}

func GetAllPermanentMacs() (macs map[string]string, err error) {
	ifaces, err := GetInterfaces()
	for _, iface := range ifaces {
		name, err := GetPermanentMac(iface.Name)
		if err != nil {
			fmt.Println(err)
		}
		macs[name.String()] = iface.HardwareAddr.String()
	}
	return
}

func SetMac(name string, mac string) (err error) {
	result, err := RunningAsRoot() 
	if err != nil {
		return
	}
	if !result {
		err = fmt.Errorf("Not running as root, insufficient privileges to set MAC on %s",
		name)
		return
	}
	result, err = IsIfUp(name)
	if err != nil {
		return
	}
	if result {
		err = fmt.Errorf("%s interface is still up, cannot set MAC", name)
		return
	}
	sockfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	defer syscall.Close(sockfd)
	iface, err := net.ParseMAC(mac)
	if err != nil {
		return
	}
	var netinfo NetInfo
	copy(netinfo.name[:], []byte(name))
	netinfo.family = syscall.AF_UNIX
	copy(netinfo.data[:], []byte(iface))
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(sockfd), SIOCSIFHWADDR, uintptr(unsafe.Pointer(&netinfo)))
	if errno != 0 {
		err = syscall.Errno(errno)
		return
	}
	return
}

func CompareMacs(first net.HardwareAddr, second net.HardwareAddr) (same bool) {
	same = first.String() == second.String()
	return
}

func MacChanged(iface string) (changed bool, err error) {
	current, err := GetCurrentMac(iface)
	if err != nil {
		return
	}
	permanent, err := GetPermanentMac(iface)
	if err != nil {
		return
	}
	if !CompareMacs(current, permanent) {
		changed = true
	} 
	return
}

func IsIfUp(name string) (result bool, err error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return
	}
        if iface.Flags&net.FlagUp != 0 {
		result = true
	}
	return
}

func RevertMac(name string) (err error) {
	_, err = net.InterfaceByName(name)
	if err != nil {
		return
	}
	mac, err := GetPermanentMac(name)
	if err != nil {
		return
	}
	err = SetMac(name, mac.String())
	return
}

func RandomizeMac(macbytes net.HardwareAddr, start int, bia bool) (mac net.HardwareAddr, err error) {
	if len(macbytes) != 6 {
		err = fmt.Errorf("Invalid size for macbytes byte array: %d", 
		len(macbytes))
		return
	}
	if (start != 0 && start != 3) {
		err = fmt.Errorf("Invalid start index: %d", start) 
		return
	}
	for i := start; i < 6; i++ {
		buf := make([]byte, 1)
		rand.Read(buf)
		if i == 0 {
			macbytes[i] = buf[0] & 0xf
		} else {
			macbytes[i] = buf[0]
		}
	}
	if bia {
		macbytes[0] &= 2
	} else {
		macbytes[0] |= 2
	}
	mac = macbytes
	return
}

func RunningAsRoot() (result bool, err error) {
	current, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}
	if current.Uid == "0" && current.Gid == "0" && current.Username == "root" {
		result = true
	}
	return 
}

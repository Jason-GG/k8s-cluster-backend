package util

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type Listen struct {
	IPV4 []uint64
	IPV6 []uint64
}

type FSInfo struct {
	MountPoint string
	Used       uint64
	Free       uint64
}

type NetIntfInfo struct {
	IPv4 string
	IPv6 string
	Rx   uint64
	Tx   uint64
}

type cpuRaw struct {
	User    uint64 // time spent in user mode
	Nice    uint64 // time spent in user mode with low priority (nice)
	System  uint64 // time spent in system mode
	Idle    uint64 // time spent in the idle task
	Iowait  uint64 // time spent waiting for I/O to complete (since Linux 2.5.41)
	Irq     uint64 // time spent servicing  interrupts  (since  2.6.0-test4)
	SoftIrq uint64 // time spent servicing softirqs (since 2.6.0-test4)
	Steal   uint64 // time spent in other OSes when running in a virtualized environment
	Guest   uint64 // time spent running a virtual CPU for guest operating systems under the control of the Linux kernel.
	Total   uint64 // total of all time fields
}

type CPUInfo struct {
	User    float64
	Nice    float64
	System  float64
	Idle    float64
	Iowait  float64
	Irq     float64
	SoftIrq float64
	Steal   float64
	Guest   float64
}

type Stats struct {
	Uptime       time.Duration
	Hostname     string
	Load1        string
	Load5        string
	Load10       string
	RunningProcs string
	TotalProcs   string
	MemTotal     uint64
	MemFree      uint64
	MemBuffers   uint64
	MemCached    uint64
	SwapTotal    uint64
	SwapFree     uint64
	FSInfos      []FSInfo
	NetIntf      map[string]NetIntfInfo
	Listens      Listen
	CPU          CPUInfo // or []CPUInfo to get all the cpu-core's stats?
}

func shellCmd(client *ssh.Client, command string) (stdout []byte, err error) {
	session, err := client.NewSession()
	if err != nil {
		return
	}

	defer session.Close()
	var buf bytes.Buffer
	session.Stdout = &buf
	err = session.Run(command)
	if err != nil {
		return
	}
	stdout = buf.Bytes()
	return
}

func shellCmdString(client *ssh.Client, command string) (stdout string, err error) {
	session, err := client.NewSession()
	if err != nil {
		//log.Print(err)
		return
	}
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf
	err = session.Run(command)
	if err != nil {
		//log.Print(err)
		return
	}
	stdout = string(buf.Bytes())

	return
}

func GetAllStates(client *ssh.Client, stats *Stats) {
	getUptime(client, stats)
	getHostname(client, stats)
	getLoad(client, stats)
	getMemInfo(client, stats)
	getFSInfo(client, stats)
	getInterfaces(client, stats)
	getInterfaceInfo(client, stats)
	getCPU(client, stats)
	getListenV4(client, stats)
	getListenV6(client, stats)
}

func getUptime(client *ssh.Client, stats *Stats) (err error) {
	uptime, err := shellCmdString(client, "/bin/cat /proc/uptime")
	if err != nil {
		return
	}

	parts := strings.Fields(uptime)
	if len(parts) == 2 {
		var upsecs float64
		upsecs, err = strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return
		}
		stats.Uptime = time.Duration(upsecs * 1e9)
	}

	return
}

func getHostname(client *ssh.Client, stats *Stats) (err error) {
	hostname, err := shellCmdString(client, "/bin/hostname -f")
	if err != nil {
		return
	}

	stats.Hostname = strings.TrimSpace(hostname)
	return
}

func getLoad(client *ssh.Client, stats *Stats) (err error) {
	line, err := shellCmdString(client, "/bin/cat /proc/loadavg")
	if err != nil {
		return
	}

	parts := strings.Fields(line)
	if len(parts) == 5 {
		stats.Load1 = parts[0]
		stats.Load5 = parts[1]
		stats.Load10 = parts[2]
		if i := strings.Index(parts[3], "/"); i != -1 {
			stats.RunningProcs = parts[3][0:i]
			if i+1 < len(parts[3]) {
				stats.TotalProcs = parts[3][i+1:]
			}
		}
	}

	return
}

const (
	v4_rem_address = "00000000:0000"
	v6_rem_address = "00000000000000000000000000000000:0000"
)

func getListenV4(client *ssh.Client, stats *Stats) (err error) {
	lines, err := shellCmdString(client, "/bin/cat /proc/net/tcp")
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(lines))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		n := len(parts)
		sl := n > 0 && parts[0] == "sl"
		if sl {
			continue
		}
		if n > 0 && parts[2] == v4_rem_address {
			hex := strings.Split(parts[1], ":")
			port, err := strconv.ParseUint(hex[1], 16, 32)
			if err != nil {
				return err
			}
			stats.Listens.IPV4 = append(stats.Listens.IPV4, port)
		}
	}
	return nil
}

func getListenV6(client *ssh.Client, stats *Stats) (err error) {
	lines, err := shellCmdString(client, "/bin/cat /proc/net/tcp6")
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(lines))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		n := len(parts)
		sl := n > 0 && parts[0] == "sl"
		if sl {
			continue
		}
		if n > 0 && parts[2] == v6_rem_address {
			hex := strings.Split(parts[1], ":")
			port, err := strconv.ParseUint(hex[1], 16, 32)
			if err != nil {
				return err
			}
			stats.Listens.IPV6 = append(stats.Listens.IPV6, port)
		}
	}
	return nil
}

func getMemInfo(client *ssh.Client, stats *Stats) (err error) {
	lines, err := shellCmdString(client, "/bin/cat /proc/meminfo")
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(lines))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 3 {
			val, err := strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				continue
			}
			val *= 1024
			switch parts[0] {
			case "MemTotal:":
				stats.MemTotal = val
			case "MemFree:":
				stats.MemFree = val
			case "Buffers:":
				stats.MemBuffers = val
			case "Cached:":
				stats.MemCached = val
			case "SwapTotal:":
				stats.SwapTotal = val
			case "SwapFree:":
				stats.SwapFree = val
			}
		}
	}

	return
}

func getFSInfo(client *ssh.Client, stats *Stats) (err error) {
	lines, err := shellCmdString(client, "/bin/df -B1")
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(lines))
	flag := 0
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		n := len(parts)
		dev := n > 0 && strings.Index(parts[0], "/dev/") == 0
		if n == 1 && dev {
			flag = 1
		} else if (n == 5 && flag == 1) || (n == 6 && dev) {
			i := flag
			flag = 0
			used, err := strconv.ParseUint(parts[2-i], 10, 64)
			if err != nil {
				continue
			}
			free, err := strconv.ParseUint(parts[3-i], 10, 64)
			if err != nil {
				continue
			}
			stats.FSInfos = append(stats.FSInfos, FSInfo{
				parts[5-i], used, free,
			})
		}
	}

	return
}

func getInterfaces(client *ssh.Client, stats *Stats) (err error) {
	var lines string
	lines, err = shellCmdString(client, "/bin/ip -o addr")
	if err != nil {
		// try /sbin/ip
		lines, err = shellCmdString(client, "/sbin/ip -o addr")
		if err != nil {
			return
		}
	}

	if stats.NetIntf == nil {
		stats.NetIntf = make(map[string]NetIntfInfo)
	}

	scanner := bufio.NewScanner(strings.NewReader(lines))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 4 && (parts[2] == "inet" || parts[2] == "inet6") {
			ipv4 := parts[2] == "inet"
			intfname := parts[1]
			if info, ok := stats.NetIntf[intfname]; ok {
				if ipv4 {
					info.IPv4 = parts[3]
				} else {
					info.IPv6 = parts[3]
				}
				stats.NetIntf[intfname] = info
			} else {
				info := NetIntfInfo{}
				if ipv4 {
					info.IPv4 = parts[3]
				} else {
					info.IPv6 = parts[3]
				}
				stats.NetIntf[intfname] = info
			}
		}
	}

	return
}

func getInterfaceInfo(client *ssh.Client, stats *Stats) (err error) {
	lines, err := shellCmdString(client, "/bin/cat /proc/net/dev")
	if err != nil {
		return
	}

	if stats.NetIntf == nil {
		return
	} // should have been here already

	scanner := bufio.NewScanner(strings.NewReader(lines))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 17 {
			intf := strings.TrimSpace(parts[0])
			intf = strings.TrimSuffix(intf, ":")
			if info, ok := stats.NetIntf[intf]; ok {
				rx, err := strconv.ParseUint(parts[1], 10, 64)
				if err != nil {
					continue
				}
				tx, err := strconv.ParseUint(parts[9], 10, 64)
				if err != nil {
					continue
				}
				info.Rx = rx
				info.Tx = tx
				stats.NetIntf[intf] = info
			}
		}
	}

	return
}

func parseCPUFields(fields []string, stat *cpuRaw) {
	numFields := len(fields)
	for i := 1; i < numFields; i++ {
		val, err := strconv.ParseUint(fields[i], 10, 64)
		if err != nil {
			continue
		}

		stat.Total += val
		switch i {
		case 1:
			stat.User = val
		case 2:
			stat.Nice = val
		case 3:
			stat.System = val
		case 4:
			stat.Idle = val
		case 5:
			stat.Iowait = val
		case 6:
			stat.Irq = val
		case 7:
			stat.SoftIrq = val
		case 8:
			stat.Steal = val
		case 9:
			stat.Guest = val
		}
	}
}

// the CPU stats that were fetched last time round
var preCPU cpuRaw

func getCPU(client *ssh.Client, stats *Stats) (err error) {
	lines, err := shellCmdString(client, "/bin/cat /proc/stat")
	if err != nil {
		return
	}

	var (
		nowCPU cpuRaw
		total  float64
	)

	scanner := bufio.NewScanner(strings.NewReader(lines))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == "cpu" { // changing here if want to get every cpu-core's stats
			parseCPUFields(fields, &nowCPU)
			break
		}
	}
	if preCPU.Total == 0 { // having no pre raw cpu data
		goto END
	}

	total = float64(nowCPU.Total - preCPU.Total)
	stats.CPU.User = float64(nowCPU.User-preCPU.User) / total * 100
	stats.CPU.Nice = float64(nowCPU.Nice-preCPU.Nice) / total * 100
	stats.CPU.System = float64(nowCPU.System-preCPU.System) / total * 100
	stats.CPU.Idle = float64(nowCPU.Idle-preCPU.Idle) / total * 100
	stats.CPU.Iowait = float64(nowCPU.Iowait-preCPU.Iowait) / total * 100
	stats.CPU.Irq = float64(nowCPU.Irq-preCPU.Irq) / total * 100
	stats.CPU.SoftIrq = float64(nowCPU.SoftIrq-preCPU.SoftIrq) / total * 100
	stats.CPU.Guest = float64(nowCPU.Guest-preCPU.Guest) / total * 100
END:
	preCPU = nowCPU
	return
}

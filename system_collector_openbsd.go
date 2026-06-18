//go:build openbsd

// OpenBSD SystemCollectorImpl for cloudflared's diagnostic collector
// (cloudflared tunnel diag). Reads memory and file-descriptor counts via
// sysctl: hw.physmem, kern.maxfiles, kern.nfiles. Based on the FreeBSD ports
// patch. Apache-2.0, same as cloudflared.
package diagnostic

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type SystemCollectorImpl struct {
	version string
}

func NewSystemCollectorImpl(version string) *SystemCollectorImpl {
	return &SystemCollectorImpl{version}
}

func (collector *SystemCollectorImpl) Collect(ctx context.Context) (*SystemInformation, error) {
	memoryInfo, memoryInfoRaw, memoryInfoErr := collectMemoryInformation(ctx)
	fdInfo, fdInfoRaw, fdInfoErr := collectFileDescriptorInformation(ctx)
	disks, disksRaw, diskErr := collectDiskVolumeInformationUnix(ctx)
	osInfo, osInfoRaw, osInfoErr := collectOSInformationUnix(ctx)

	var memoryMaximum, memoryCurrent, fileDescriptorMaximum, fileDescriptorCurrent uint64
	var osSystem, name, osVersion, osRelease, architecture string
	gerror := SystemInformationGeneralError{}

	if memoryInfoErr != nil {
		gerror.MemoryInformationError = SystemInformationError{Err: memoryInfoErr, RawInfo: memoryInfoRaw}
	} else {
		memoryMaximum = memoryInfo.MemoryMaximum
		memoryCurrent = memoryInfo.MemoryCurrent
	}

	if fdInfoErr != nil {
		gerror.FileDescriptorsInformationError = SystemInformationError{Err: fdInfoErr, RawInfo: fdInfoRaw}
	} else {
		fileDescriptorMaximum = fdInfo.FileDescriptorMaximum
		fileDescriptorCurrent = fdInfo.FileDescriptorCurrent
	}

	if diskErr != nil {
		gerror.DiskVolumeInformationError = SystemInformationError{Err: diskErr, RawInfo: disksRaw}
	}

	if osInfoErr != nil {
		gerror.OperatingSystemInformationError = SystemInformationError{Err: osInfoErr, RawInfo: osInfoRaw}
	} else {
		osSystem = osInfo.OsSystem
		name = osInfo.Name
		osVersion = osInfo.OsVersion
		osRelease = osInfo.OsRelease
		architecture = osInfo.Architecture
	}

	info := NewSystemInformation(
		memoryMaximum, memoryCurrent,
		fileDescriptorMaximum, fileDescriptorCurrent,
		osSystem, name, osVersion, osRelease, architecture,
		collector.version, runtime.Version(), runtime.GOARCH, disks,
	)

	return info, gerror
}

func collectMemoryInformation(ctx context.Context) (*MemoryInformation, string, error) {
	command := exec.CommandContext(ctx, "sysctl", "-n", "hw.physmem", "hw.pagesize")
	stdout, err := command.Output()
	if err != nil {
		return nil, "", fmt.Errorf("error retrieving output from command '%s': %w", command.String(), err)
	}
	output := string(stdout)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 2 {
		return nil, output, fmt.Errorf("unexpected sysctl output format")
	}
	physmem, err := strconv.ParseUint(lines[0], 10, 64)
	if err != nil {
		return nil, output, fmt.Errorf("error parsing physmem: %w", err)
	}
	// OpenBSD has no simple free-memory sysctl scalar; report total only.
	memoryInfo := &MemoryInformation{
		MemoryMaximum: physmem / 1024,
		MemoryCurrent: 0,
	}
	return memoryInfo, output, nil
}

func collectFileDescriptorInformation(ctx context.Context) (*FileDescriptorInformation, string, error) {
	command := exec.CommandContext(ctx, "sysctl", "-n", "kern.maxfiles", "kern.nfiles")
	stdout, err := command.Output()
	if err != nil {
		return nil, "", fmt.Errorf("error retrieving output from command '%s': %w", command.String(), err)
	}
	output := string(stdout)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 2 {
		return nil, output, fmt.Errorf("unexpected sysctl output format")
	}
	maxFiles, err := strconv.ParseUint(lines[0], 10, 64)
	if err != nil {
		return nil, output, fmt.Errorf("error parsing maxfiles: %w", err)
	}
	openFiles, err := strconv.ParseUint(lines[1], 10, 64)
	if err != nil {
		return nil, output, fmt.Errorf("error parsing nfiles: %w", err)
	}
	fdInfo := &FileDescriptorInformation{
		FileDescriptorMaximum: maxFiles,
		FileDescriptorCurrent: openFiles,
	}
	return fdInfo, output, nil
}

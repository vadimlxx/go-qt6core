package qdata

import (
	"strings"
	"testing"
)

func TestQSysInfoWordSize(t *testing.T) {
	if WordSize != 32 && WordSize != 64 {
		t.Fatalf("unexpected WordSize: %d", WordSize)
	}
}

func TestQSysInfoByteOrder(t *testing.T) {
	if ByteOrder != BigEndian && ByteOrder != LittleEndian {
		t.Fatalf("unexpected ByteOrder: %d", ByteOrder)
	}
}

func TestQSysInfoBuildCpuArchitecture(t *testing.T) {
	s := QSysInfoBuildCpuArchitecture()
	if s == "" {
		t.Fatal("buildCpuArchitecture returned empty string")
	}
	if strings.ContainsRune(s, '\x00') {
		t.Fatalf("buildCpuArchitecture contains NUL: %q", s)
	}
}

func TestQSysInfoCurrentCpuArchitecture(t *testing.T) {
	s := QSysInfoCurrentCpuArchitecture()
	if s == "" {
		t.Fatal("currentCpuArchitecture returned empty string")
	}
}

func TestQSysInfoBuildAbi(t *testing.T) {
	s := QSysInfoBuildAbi()
	if s == "" {
		t.Fatal("buildAbi returned empty string")
	}
}

func TestQSysInfoKernelType(t *testing.T) {
	s := QSysInfoKernelType()
	if s == "" {
		t.Fatal("kernelType returned empty string")
	}
}

func TestQSysInfoKernelVersion(t *testing.T) {
	s := QSysInfoKernelVersion()
	if s == "" {
		t.Fatal("kernelVersion returned empty string")
	}
}

func TestQSysInfoProductType(t *testing.T) {
	s := QSysInfoProductType()
	if s == "" {
		t.Fatal("productType returned empty string")
	}
}

func TestQSysInfoProductVersion(t *testing.T) {
	s := QSysInfoProductVersion()
	if s == "" {
		t.Fatal("productVersion returned empty string")
	}
}

func TestQSysInfoPrettyProductName(t *testing.T) {
	s := QSysInfoPrettyProductName()
	if s == "" {
		t.Fatal("prettyProductName returned empty string")
	}
}

func TestQSysInfoMachineHostName(t *testing.T) {
	s := QSysInfoMachineHostName()
	// Hostname can be empty in some container/CI environments — that is
	// acceptable. Only check it does not contain an embedded NUL.
	if strings.ContainsRune(s, '\x00') {
		t.Fatalf("machineHostName returned string with NUL: %q", s)
	}
}

func TestQSysInfoMachineUniqueId(t *testing.T) {
	b := QSysInfoMachineUniqueId()
	// May be nil on some platforms — that is expected. When non-nil it must
	// be a non-empty byte slice. The value is opaque; do not decode it as
	// text, since it may contain embedded NUL bytes.
	if b != nil && len(b) == 0 {
		t.Fatal("machineUniqueId returned non-nil empty slice")
	}
}

func TestQSysInfoBootUniqueId(t *testing.T) {
	b := QSysInfoBootUniqueId()
	// May be nil on some platforms — that is expected. When non-nil it must
	// be a non-empty byte slice.
	if b != nil && len(b) == 0 {
		t.Fatal("bootUniqueId returned non-nil empty slice")
	}
}
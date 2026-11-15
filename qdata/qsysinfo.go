package qdata

/*
   QSysInfo в Qt — чисто статический класс с двумя enum-ами (Sizes, Endian)
   и одиннадцатью статическими getter-ами. В Go он представлен:
     - пакетными функциями QSysInfoBuildCpuArchitecture() и т.п.;
     - переменными WordSize и ByteOrder в этом файле;
     - типами Sizes и Endian в qsysinfo_sizes.go / qsysinfo_endian.go.

   Возвращаемые значения:
     - QString-методы отдают Go string (UTF-8, NUL-терминированный буфер
       аллоцируется в C++, копируется в Go, затем освобождается через
       qsysinfo_free_string). Так избегаем циклической зависимости на qtext.
     - QByteArray-методы отдают []byte с явной длиной — на случай
       embedded NUL в уникальных идентификаторах (machine/boot unique id).
       nil означает null QByteArray (например, platform unique id
       недоступен); пустой не-nil []byte для null-случая не возвращается,
       т.к. разница null/empty через []byte не выразительна — оба дают
       len(x) == 0.

   WordSize и ByteOrder — переменные, а не const, потому что их значения
   зависят от целевой платформы (sizeof(void*)<<3 и порядок байтов сборки).
   Значения запрашиваются у Qt один раз при инициализации пакета.

 Что сознательно не вошло:

   - Всё, что не входит в публичный API QSysInfo из qsysinfo.h —
     у QSysInfo нет ни protected-секций, ни членов-данных, ни свободных
     функций, ни операторов.
   - Обёртки над QString / QByteArray на стороне Go — из-за явного
     требования не тянуть qtext в qdata.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stddef.h>
#include "qsysinfo_c.h"
*/
import "C"

import "unsafe"

// WordSize mirrors QSysInfo::WordSize (== sizeof(void*) << 3).
// Qt: QSysInfo::Sizes::WordSize.
var WordSize = Sizes(C.qsysinfo_wordsize())

// ByteOrder mirrors QSysInfo::ByteOrder — equal to BigEndian or LittleEndian
// depending on the byte order of the build.
// Qt: QSysInfo::Endian::ByteOrder.
var ByteOrder = Endian(C.qsysinfo_byte_order())

// QSysInfoBuildCpuArchitecture returns the CPU architecture the application
// binary was compiled for.
// Qt: QSysInfo::buildCpuArchitecture() → QString (returned here as string).
func QSysInfoBuildCpuArchitecture() string {
	return goString(C.qsysinfo_build_cpu_architecture())
}

// QSysInfoCurrentCpuArchitecture returns the CPU architecture the application
// is currently running on (may differ from BuildCpuArchitecture under
// emulation or multi-arch).
// Qt: QSysInfo::currentCpuArchitecture() → QString (returned here as string).
func QSysInfoCurrentCpuArchitecture() string {
	return goString(C.qsysinfo_current_cpu_architecture())
}

// QSysInfoBuildAbi returns the full ABI (architecture, endianness, word size)
// the application binary was compiled for.
// Qt: QSysInfo::buildAbi() → QString (returned here as string).
func QSysInfoBuildAbi() string {
	return goString(C.qsysinfo_build_abi())
}

// QSysInfoKernelType returns the kernel type of the operating system
// (e.g. "linux", "darwin", "winnt").
// Qt: QSysInfo::kernelType() → QString (returned here as string).
func QSysInfoKernelType() string {
	return goString(C.qsysinfo_kernel_type())
}

// QSysInfoKernelVersion returns the kernel version of the operating system.
// Qt: QSysInfo::kernelVersion() → QString (returned here as string).
func QSysInfoKernelVersion() string {
	return goString(C.qsysinfo_kernel_version())
}

// QSysInfoProductType returns the product name of the OS as reported by Qt
// (e.g. "debian", "windows", "macos").
// Qt: QSysInfo::productType() → QString (returned here as string).
func QSysInfoProductType() string {
	return goString(C.qsysinfo_product_type())
}

// QSysInfoProductVersion returns the product version of the OS.
// Qt: QSysInfo::productVersion() → QString (returned here as string).
func QSysInfoProductVersion() string {
	return goString(C.qsysinfo_product_version())
}

// QSysInfoPrettyProductName returns a human-readable product name of the OS.
// Qt: QSysInfo::prettyProductName() → QString (returned here as string).
func QSysInfoPrettyProductName() string {
	return goString(C.qsysinfo_pretty_product_name())
}

// QSysInfoMachineHostName returns the machine's hostname.
// Qt: QSysInfo::machineHostName() → QString (returned here as string).
func QSysInfoMachineHostName() string {
	return goString(C.qsysinfo_machine_host_name())
}

// QSysInfoMachineUniqueId returns a unique identifier for this machine,
// or nil if the platform does not expose one.
// Qt: QSysInfo::machineUniqueId() → QByteArray (returned here as []byte).
func QSysInfoMachineUniqueId() []byte {
	var n C.size_t
	p := C.qsysinfo_machine_unique_id(&n)
	return goBytes(p, n)
}

// QSysInfoBootUniqueId returns a unique identifier for the current boot,
// or nil if the platform does not expose one.
// Qt: QSysInfo::bootUniqueId() → QByteArray (returned here as []byte).
func QSysInfoBootUniqueId() []byte {
	var n C.size_t
	p := C.qsysinfo_boot_unique_id(&n)
	return goBytes(p, n)
}

// goString converts a heap-allocated NUL-terminated C string to a Go string
// and releases the C allocation. Returns "" for a nil pointer.
func goString(p *C.char) string {
	if p == nil {
		return ""
	}
	defer C.qsysinfo_free_string(p)
	return C.GoString(p)
}

// goBytes copies n bytes from a heap-allocated C buffer into a Go []byte
// (preserving embedded NUL bytes) and releases the C allocation.
// Returns nil for a nil pointer or zero length.
func goBytes(p *C.char, n C.size_t) []byte {
	if p == nil {
		return nil
	}
	defer C.qsysinfo_free_string(p)
	if n == 0 {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(p), C.int(n))
}

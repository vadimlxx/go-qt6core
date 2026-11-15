#include "qsysinfo_c.h"

#include <QtCore/QSysInfo>
#include <QtCore/QString>
#include <QtCore/QByteArray>

#include <cstdlib>
#include <cstring>

namespace {

char* copyUtf8(const QString& s) {
    const QByteArray ba = s.toUtf8();
    const size_t n = static_cast<size_t>(ba.size());
    char* out = static_cast<char*>(std::malloc(n + 1));
    if (!out)
        return nullptr;
    if (n)
        std::memcpy(out, ba.constData(), n);
    out[n] = '\0';
    return out;
}

char* copyBytes(const QByteArray& ba, size_t* out_size) {
    const size_t n = static_cast<size_t>(ba.size());
    char* out = static_cast<char*>(std::malloc(n ? n : 1));
    if (!out) {
        if (out_size)
            *out_size = 0;
        return nullptr;
    }
    if (n)
        std::memcpy(out, ba.constData(), n);
    if (out_size)
        *out_size = n;
    return out;
}

} // namespace

extern "C" {

int qsysinfo_wordsize(void) {
    return QSysInfo::WordSize;
}

int qsysinfo_byte_order(void) {
    return static_cast<int>(QSysInfo::ByteOrder);
}

char* qsysinfo_build_cpu_architecture(void) {
    return copyUtf8(QSysInfo::buildCpuArchitecture());
}
char* qsysinfo_current_cpu_architecture(void) {
    return copyUtf8(QSysInfo::currentCpuArchitecture());
}
char* qsysinfo_build_abi(void) {
    return copyUtf8(QSysInfo::buildAbi());
}
char* qsysinfo_kernel_type(void) {
    return copyUtf8(QSysInfo::kernelType());
}
char* qsysinfo_kernel_version(void) {
    return copyUtf8(QSysInfo::kernelVersion());
}
char* qsysinfo_product_type(void) {
    return copyUtf8(QSysInfo::productType());
}
char* qsysinfo_product_version(void) {
    return copyUtf8(QSysInfo::productVersion());
}
char* qsysinfo_pretty_product_name(void) {
    return copyUtf8(QSysInfo::prettyProductName());
}
char* qsysinfo_machine_host_name(void) {
    return copyUtf8(QSysInfo::machineHostName());
}

char* qsysinfo_machine_unique_id(size_t* out_size) {
    return copyBytes(QSysInfo::machineUniqueId(), out_size);
}
char* qsysinfo_boot_unique_id(size_t* out_size) {
    return copyBytes(QSysInfo::bootUniqueId(), out_size);
}

void qsysinfo_free_string(char* s) {
    std::free(s);
}

} // extern "C"
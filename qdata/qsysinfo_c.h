/* C ABI for Qt 6.11 QtCore/QSysInfo. */
#ifndef GO_QT6CORE_QDATA_QSYSINFO_C_H
#define GO_QT6CORE_QDATA_QSYSINFO_C_H

#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

/* Enum helpers: values depend on the target platform, so they are queried
   rather than hardcoded in Go. */
int qsysinfo_wordsize(void);
int qsysinfo_byte_order(void);

/* QString-returning getters. Each returns a heap-allocated NUL-terminated
   UTF-8 buffer; the caller must free it with qsysinfo_free_string(). */
char* qsysinfo_build_cpu_architecture(void);
char* qsysinfo_current_cpu_architecture(void);
char* qsysinfo_build_abi(void);
char* qsysinfo_kernel_type(void);
char* qsysinfo_kernel_version(void);
char* qsysinfo_product_type(void);
char* qsysinfo_product_version(void);
char* qsysinfo_pretty_product_name(void);
char* qsysinfo_machine_host_name(void);

/* QByteArray-returning getters. Each returns a heap-allocated raw byte
   buffer whose length is written to *out_size; the caller must free it with
   qsysinfo_free_string(). Size is preserved so embedded NUL bytes (if any)
   survive the round-trip. */
char* qsysinfo_machine_unique_id(size_t* out_size);
char* qsysinfo_boot_unique_id(size_t* out_size);

/* Release a buffer returned by any of the getters above. */
void qsysinfo_free_string(char* s);

#ifdef __cplusplus
}
#endif

#endif
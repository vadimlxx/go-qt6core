/* Single source of truth for the minimum Qt version. */
#ifndef GO_QT6CORE_CHECK_QT_VERSION_H
#define GO_QT6CORE_CHECK_QT_VERSION_H

#include <QtCore/qglobal.h>

#if !defined(QT_VERSION) || QT_VERSION < QT_VERSION_CHECK(6, 11, 0)
#  error "go-qt6core requires Qt >= 6.11"
#endif

#endif

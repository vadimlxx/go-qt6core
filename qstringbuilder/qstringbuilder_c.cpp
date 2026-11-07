/* Order: version check first, then real Qt headers. */
#include "../internal/qt6core/check_qt_version.h"

#include <QString>
#include <QByteArray>
#include <QStringView>
#include <QLatin1StringView>
#include <QLatin1Char>
#include <QByteArrayView>
#include <QChar>
#include <QUtf8StringView>

#include "qstringbuilder_c.h"

extern "C" {

void qstringbuilder_append_to_qstring(void *qstr, const qsb_part_t *parts, size_t n) {
    QString *q = static_cast<QString *>(qstr);
    for (size_t i = 0; i < n; ++i) {
        const qsb_part_t &p = parts[i];
        switch (p.kind) {
        case QSB_KIND_UTF8_STRING:
            q->append(QUtf8StringView(reinterpret_cast<const char *>(p.data), p.size));
            break;
        case QSB_KIND_QLATIN1CHAR:
            q->append(QLatin1Char(static_cast<char>(p.tag)));
            break;
        case QSB_KIND_QCHAR:
            q->append(QChar(static_cast<char16_t>(p.tag)));
            break;
        case QSB_KIND_QSTRING:
            q->append(*reinterpret_cast<const QString *>(p.data));
            break;
        case QSB_KIND_QSTRINGVIEW:
            q->append(QStringView(reinterpret_cast<const char16_t *>(p.data), p.size));
            break;
        case QSB_KIND_QLATIN1STRINGVIEW:
            q->append(QLatin1StringView(reinterpret_cast<const char *>(p.data), p.size));
            break;
        case QSB_KIND_QBYTEARRAY: {
            const QByteArray *ba = reinterpret_cast<const QByteArray *>(p.data);
            q->append(QUtf8StringView(ba->constData(), ba->size()));
            break;
        }
        case QSB_KIND_QBYTEARRAYVIEW:
            q->append(QUtf8StringView(reinterpret_cast<const char *>(p.data), p.size));
            break;
        }
    }
}

void qstringbuilder_append_to_qbytearray(void *qba, const qsb_part_t *parts, size_t n) {
    QByteArray *b = static_cast<QByteArray *>(qba);
    for (size_t i = 0; i < n; ++i) {
        const qsb_part_t &p = parts[i];
        switch (p.kind) {
        case QSB_KIND_UTF8_STRING:
        case QSB_KIND_QLATIN1STRINGVIEW:
        case QSB_KIND_QBYTEARRAYVIEW:
            b->append(reinterpret_cast<const char *>(p.data), p.size);
            break;
        case QSB_KIND_QLATIN1CHAR:
            b->append(static_cast<char>(p.tag));
            break;
        case QSB_KIND_QCHAR:
            b->append(QString(QChar(static_cast<char16_t>(p.tag))).toUtf8());
            break;
        case QSB_KIND_QSTRING:
            b->append(reinterpret_cast<const QString *>(p.data)->toUtf8());
            break;
        case QSB_KIND_QSTRINGVIEW:
            b->append(QString(QStringView(reinterpret_cast<const char16_t *>(p.data), p.size)).toUtf8());
            break;
        case QSB_KIND_QBYTEARRAY:
            b->append(*reinterpret_cast<const QByteArray *>(p.data));
            break;
        }
    }
}

} // extern "C"
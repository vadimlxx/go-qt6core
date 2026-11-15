#include "qtextboundaryfinder_c.h"

#include <QtCore/QTextBoundaryFinder>
#include <QtCore/QString>
#include <QtCore/QStringView>

extern "C" {

/* --- lifecycle --- */
QTextBoundaryFinder* qtextboundaryfinder_new_empty() {
    return new QTextBoundaryFinder();
}

QTextBoundaryFinder* qtextboundaryfinder_new_copy(const QTextBoundaryFinder* other) {
    return new QTextBoundaryFinder(*other);
}

QTextBoundaryFinder* qtextboundaryfinder_new_from_qstring(int type, const QString* s) {
    return new QTextBoundaryFinder(static_cast<QTextBoundaryFinder::BoundaryType>(type), *s);
}

QTextBoundaryFinder* qtextboundaryfinder_new_from_view(int type, const uint16_t* data, ptrdiff_t size) {
    QStringView sv(reinterpret_cast<const char16_t*>(data), size);
    return new QTextBoundaryFinder(static_cast<QTextBoundaryFinder::BoundaryType>(type), sv);
}

void qtextboundaryfinder_delete(QTextBoundaryFinder* self) { delete self; }

void qtextboundaryfinder_swap(QTextBoundaryFinder* self, QTextBoundaryFinder* other) {
    self->swap(*other);
}

/* --- accessors --- */
bool qtextboundaryfinder_is_valid(const QTextBoundaryFinder* self) { return self->isValid(); }

int qtextboundaryfinder_type(const QTextBoundaryFinder* self) {
    return static_cast<int>(self->type());
}

QString* qtextboundaryfinder_string(const QTextBoundaryFinder* self) {
    return new QString(self->string());
}

/* --- navigation --- */
void qtextboundaryfinder_to_start(QTextBoundaryFinder* self) { self->toStart(); }
void qtextboundaryfinder_to_end(QTextBoundaryFinder* self) { self->toEnd(); }

ptrdiff_t qtextboundaryfinder_position(const QTextBoundaryFinder* self) {
    return static_cast<ptrdiff_t>(self->position());
}

void qtextboundaryfinder_set_position(QTextBoundaryFinder* self, ptrdiff_t pos) {
    self->setPosition(static_cast<qsizetype>(pos));
}

ptrdiff_t qtextboundaryfinder_to_next_boundary(QTextBoundaryFinder* self) {
    return static_cast<ptrdiff_t>(self->toNextBoundary());
}

ptrdiff_t qtextboundaryfinder_to_previous_boundary(QTextBoundaryFinder* self) {
    return static_cast<ptrdiff_t>(self->toPreviousBoundary());
}

/* --- predicates --- */
bool qtextboundaryfinder_is_at_boundary(const QTextBoundaryFinder* self) {
    return self->isAtBoundary();
}

int qtextboundaryfinder_boundary_reasons(const QTextBoundaryFinder* self) {
    return static_cast<int>(self->boundaryReasons().toInt());
}

} // extern "C"
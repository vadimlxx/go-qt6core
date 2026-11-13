# Qt6 QtCore binding for Golang, under construction

Qt 6.11

#### All string classes are assigned to QText:

- QChar
- QLatin1Char
- QByteArray

Go not supported operator[], use At/SetAt

No FromRawData, use FromBytes

- QByteArrayList

QList::begin/end -- Go not supported, replaced with ToSlice()

operator+ -- use ToSlice + FromQByteArrays.

- QByteArrayView
- QStringView

Use Front()/Back() for bytes, First(n)/Last(n) for subviews — as in qbytearray.

STL-iterators (begin/end/rbegin/rend) replaced with Bytes()

- QLatin1StringView ( = QLatin1String)

toWCharArray() — wchar_t platform size mismatch

- QString

Go-panic on out-of-range for At/First/Last/Sliced/SlicedN/Chopped

Left/Right/Mid without panic

- QUtf8StringView

fromArray()  — trivial; use FromBytes instead

No reverse_iterator methods (rbegin/rend/crbegin/crend)

QString overloads of compare/equal (side) — would create an import cycle with qstringview; use qstring.QString → .Utf8() → FromQByteArray, or convert to QStringView / QLatin1StringView first

- QAnyStringView

fromArray() — use FromBytes / FromString

Comparison operators (==, !=, <, <=, >, >=) — use Compare / Equal

operator std::string_view / operator std::u8string_view — use Bytes()

- QStringBuilder

operator% (a % b): Go has no user-defined operators. Use Concat or Append.

operator+ : same; use Concat / AppendPart.

QStringBuilder<A,B>::size() and ::isNull(): not exposed (semantics are "approximate" in Qt — upper bound rather than exact), use ToQString().Size() / ToQByteArray().Size() instead.

QStringBuilderCommon::toUpper() / toLower(): not exposed as builder methods (return type depends on flavor). Use ToQString().ToUpper() / ToLower() or ToQByteArray().ToUpper() / ToLower().

- QStringConverter

Invoking Interface::{toUtf16, toUtf16Len, fromUtf16, fromUtf16Len} from Go — requires cgo function-pointer support; the pointers are exposed as unsafe.Pointer, but calling them is not implemented.

operator|, operator& — Go has no operator overloading; use Flags.Has(flag).

- QStringEncoder

QStringEncoder::operator()  Go has no operator(); use Encode / EncodeView.

QStringEncoder::encode(const QString&) vs encode(QStringView): Go has one method per input type; no overload set.

- QStringDecoder

QStringDecoder::operator() Go has no operator(); use Decode / DecodeView.

QStringDecoder::decode(const QByteArray&) vs decode(QByteArrayView): Go has one method per input type; no overload set.

EncodedData<T> as a public type: not exposed; Decode / DecodeView return the materialized QString directly.

- QStringMatcher
- QBitArray
- QBitRef

Operator overloads & / | / ^ / ~ / &= / |= / ^=: Go has no operator overloading. Use package functions And(a, b) / Or(a, b) / Xor(a, b) / Not(a) and methods (*QBitArray).AndAssign(other) / OrAssign / XorAssign. They are semantically equivalent.

QBitRef::operator bool() / operator!() / operator=(bool) / operator=(const QBitRef&): exposed as methods Bool() / Not() / Set(bool) / Assign(*QBitRef).

QBitArray::count() (no-arg deprecated alias for size()): not exposed; use Size().

- QRegularExpression
- QRegularExpressionMatch
- QRegularExpressionMatchIterator

#### All geometry classes are assigned to qgeometry:

- QMargins & QMarginsF

Because Go has no operator overloading, Qt operators map to methods -- see margins_doc

- QSize & QSizeF

```
Qt-constructor QSize() returns QSize(-1, -1), NOT (0, 0).
   - NewQSizeDefault()  → (-1, -1)  — default QSize()
   - NewQSizeZero()     → ( 0,  0)  — «zero» QSize
   - NewQSize(w, h)     → ( w,  h)  — manual
```

- QPoint & QPointF

Operator mapping (QPoint / QPointF) -- see qpoint_doc

qHash(QPointF) is deleted in Qt itself.

rx() / ry() — covered by SetX / SetY.

std::tuple protocol / structured bindings — no Go equivalent; use X()/Y().

get<I>(QPoint) / get<I>(QPointF) — same; use X()/Y().

- QLine & QLineF

operator== / operator!= — use Equal / NotEqual.

- QRect & QRectF

std::tuple protocol / structured bindings — no Go equivalent; use Left()/Top()/Width()/Height().

Free operator+(QRect, QMargins) / operator+(QMargins, QRect) — covered by MarginsAdded (commutative in effect).

Free operator-(QRect, QMargins) — covered by MarginsRemoved.

#### All date, time & locate classes are assigned to qdatetime:

- QLocate

see qlocate_doc

- QCalendar

The Go zero value QCalendar{} is NOT a valid calendar:

QCalendar{}       → d_ptr == nil → IsValid() == false
NewQCalendar()    → Qt's QCalendar() → valid, Gregorian

- QTime

QTime{}          → msecs = 0       → valid time 00:00:00.000 (midnight)
NewQTime()       → null QTime      → mirrors Qt's QTime() (msecs = -1)
NewQTimeFromHMS(h,m,s,ms)          → explicit time

- QDate

QDate{}          → Julian day 0    → valid date (1 January 4713 BCE)
NewQDate()       → null QDate      → mirrors Qt's QDate() (Julian day INT64_MIN)
NewQDateFromYMD(y,m,d)             → explicit date

- QDateTime

The Go zero value QDateTime{} is invalid (nil pointer). Use NewQDateTime() for a null QDateTime that mirrors Qt's QDateTime().

- QTimeZone

The Go zero value QTimeZone{} is invalid (nil pointer). Use NewQTimeZone() for a null QTimeZone that mirrors Qt's QTimeZone(). Use Clone() for an independent copy, Close() to free immediately; otherwise a finalizer handles it.


#### All classes that store and process various data without dependencies are located in qdata.

- QInternal

- QKeyCombination

All details see in doc.go files

### License

All Qt6 Licenses

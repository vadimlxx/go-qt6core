# Qt6 QtCore binding for Golang, under construction

Qt 6.11

#### Order:

- QChar
- QLatin1Char
- QByteArray

Go not supported operator[], use At/SetAt

No FromRawData, use FromBytes

- QByteArrayList

QList::begin/end -- Go not supported, replaced with ToSlice()

operator+ -- use ToSlice + FromQByteArrays.

- QByteArrayView
- QLatin1Char
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




All details see in doc.go files

### License

All Qt6 Licenses

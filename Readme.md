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

toWCharArray()   — wchar_t platform size mismatch

- QString

Go-panic on out-of-range for At/First/Last/Sliced/SlicedN/Chopped

Left/Right/Mid without panic

- QUtf8StringView

fromArray()  — trivial; use FromBytes instead

No reverse_iterator methods (rbegin/rend/crbegin/crend)

QString overloads of compare/equal (side) — would create an import cycle with qstringview; use qstring.QString → .Utf8() → FromQByteArray, or convert to QStringView / QLatin1StringView first


### License

All Qt6 Licenses

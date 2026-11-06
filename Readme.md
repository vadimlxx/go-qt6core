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


### License

All Qt6 Licenses

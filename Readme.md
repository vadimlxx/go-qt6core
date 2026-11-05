# Qt6 QtCore binding for Golang, in development

Qt 6.11

#### Order:

- QChar
- QByteArray

Go not supported operator[], use At/SetAt

No FromRawData, use FromBytes

- QByteArrayList

QList::begin/end -- Go not supported, replaced with ToSlice()

operator+ -- use ToSlice + FromQByteArrays.

- QByteArrayView

Use Front()/Back() for bytes, First(n)/Last(n) for subviews — as in qbytearray.

STL-iterators (begin/end/rbegin/rend) replaced with Bytes()

### License

All Qt6 Licenses

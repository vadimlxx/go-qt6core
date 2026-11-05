# Qt6 QtCore binding for Golang, in development

Qt 6.11

#### Order:

- QChar
- QByteArray

Go not supported operator[], use At/SetAt

no FromRawData, use FromBytes

- QByteArrayList

QList::begin/end -- Go not supported, replaced with ToSlice()

operator+ -- use ToSlice + FromQByteArrays.

### License

All Qt6 Licenses

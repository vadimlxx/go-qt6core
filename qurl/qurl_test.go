package qurl

import (
	"bytes"
	"testing"
)

// ---- construction & basic parsing ---------------------------------------

func TestQUrlEmpty(t *testing.T) {
	u := NewQUrl()
	defer u.Close()
	if !u.IsEmpty() {
		t.Fatal("NewQUrl should be empty")
	}
	if u.IsValid() {
		t.Fatal("empty QUrl should be invalid")
	}
}

func TestQUrlFromString(t *testing.T) {
	u := QUrlFromString("https://user:pass@example.com:8080/path/file?q=1#frag", TolerantMode)
	defer u.Close()
	if !u.IsValid() {
		t.Fatalf("parse failed: %s", u.ErrorString())
	}
	if got := u.SchemeString(); got != "https" {
		t.Fatalf("scheme: %q", got)
	}
	if got := u.Host(FullyDecoded); got != "example.com" {
		t.Fatalf("host: %q", got)
	}
	if got := u.Port(-1); got != 8080 {
		t.Fatalf("port: %d", got)
	}
	if got := u.UserName(FullyDecoded); got != "user" {
		t.Fatalf("username: %q", got)
	}
	if got := u.Password(FullyDecoded); got != "pass" {
		t.Fatalf("password: %q", got)
	}
	if got := u.Path(FullyDecoded); got != "/path/file" {
		t.Fatalf("path: %q", got)
	}
	if got := u.FileName(FullyDecoded); got != "file" {
		t.Fatalf("filename: %q", got)
	}
	if !u.HasQuery() {
		t.Fatal("expected has query")
	}
	if got := u.Query(FullyDecoded); got != "q=1" {
		t.Fatalf("query: %q", got)
	}
	if !u.HasFragment() {
		t.Fatal("expected has fragment")
	}
	if got := u.Fragment(FullyDecoded); got != "frag" {
		t.Fatalf("fragment: %q", got)
	}
}

func TestQUrlInvalid(t *testing.T) {
	u := QUrlFromString("://bad url", StrictMode)
	defer u.Close()
	if u.IsValid() {
		t.Fatal("expected invalid")
	}
	if u.ErrorString() == "" {
		t.Fatal("expected non-empty error string")
	}
}

func TestQUrlCloneSwap(t *testing.T) {
	a := QUrlFromString("https://a.example/", TolerantMode)
	b := QUrlFromString("https://b.example/", TolerantMode)
	defer a.Close()
	defer b.Close()

	c := a.Clone()
	defer c.Close()
	if c.Host(FullyDecoded) != "a.example" {
		t.Fatalf("clone host: %q", c.Host(FullyDecoded))
	}

	a.Swap(b)
	if a.Host(FullyDecoded) != "b.example" {
		t.Fatalf("after swap a: %q", a.Host(FullyDecoded))
	}
	if b.Host(FullyDecoded) != "a.example" {
		t.Fatalf("after swap b: %q", b.Host(FullyDecoded))
	}
}

// ---- string representation ---------------------------------------------

func TestQUrlToString(t *testing.T) {
	u := QUrlFromString("https://example.com/path?q=1#frag", TolerantMode)
	defer u.Close()

	got := u.URLString(FormattingOptions(RemoveQuery | RemoveFragment))
	if got != "https://example.com/path" {
		t.Fatalf("URLString(RemoveQuery|RemoveFragment): %q", got)
	}
}

func TestQUrlToEncoded(t *testing.T) {
	u := QUrlFromString("https://example.com/a b?x=y", TolerantMode)
	defer u.Close()
	b := u.ToEncodedBytes(FullyEncoded)
	if !bytes.Contains(b, []byte("%20")) {
		t.Fatalf("expected %%20 in encoded URL, got %q", b)
	}
}

func TestQUrlAdjusted(t *testing.T) {
	u := QUrlFromString("https://user:pass@example.com:8443/path", TolerantMode)
	defer u.Close()
	adj := u.Adjusted(FormattingOptions(RemoveUserInfo | RemovePort))
	defer adj.Close()
	got := adj.ToString(FormattingOptions(PrettyDecoded))
	if got != "https://example.com/path" {
		t.Fatalf("Adjusted: %q", got)
	}
}

// ---- setters ------------------------------------------------------------

func TestQUrlSetComponents(t *testing.T) {
	u := NewQUrl()
	defer u.Close()

	u.SetSchemeString("https")
	u.SetHostString("example.com", DecodedMode)
	u.SetPort(8443)
	u.SetPathString("/api/v1", DecodedMode)
	u.SetQueryString("k=v", TolerantMode)
	u.SetFragmentString("frag", TolerantMode)

	if got := u.ToString(FormattingOptions(PrettyDecoded)); got != "https://example.com:8443/api/v1?k=v#frag" {
		t.Fatalf("assembled URL: %q", got)
	}
}

func TestQUrlSetUserInfo(t *testing.T) {
	u := QUrlFromString("https://example.com/", TolerantMode)
	defer u.Close()
	// Qt forbids DecodedMode in setUserInfo (only Tolerant/Strict are
	// permitted); use TolerantMode as documented.
	u.SetUserInfoString("bob:secret", TolerantMode)
	if got := u.UserName(FullyDecoded); got != "bob" {
		t.Fatalf("username: %q", got)
	}
	if got := u.Password(FullyDecoded); got != "secret" {
		t.Fatalf("password: %q", got)
	}
}

// ---- relations ----------------------------------------------------------

func TestQUrlResolved(t *testing.T) {
	base := QUrlFromString("https://example.com/a/b/c", TolerantMode)
	rel := QUrlFromString("../d", TolerantMode)
	defer base.Close()
	defer rel.Close()

	r := base.Resolved(rel)
	defer r.Close()
	if got := r.Path(FullyDecoded); got != "/a/d" {
		t.Fatalf("resolved path: %q", got)
	}
}

func TestQUrlIsRelative(t *testing.T) {
	abs := QUrlFromString("https://example.com/", TolerantMode)
	rel := QUrlFromString("/path/to", TolerantMode)
	defer abs.Close()
	defer rel.Close()
	if abs.IsRelative() {
		t.Fatal("absolute URL reported relative")
	}
	if !rel.IsRelative() {
		t.Fatal("relative URL reported absolute")
	}
}

func TestQUrlIsParentOf(t *testing.T) {
	parent := QUrlFromString("https://example.com/a", TolerantMode)
	child := QUrlFromString("https://example.com/a/b", TolerantMode)
	defer parent.Close()
	defer child.Close()
	if !parent.IsParentOf(child) {
		t.Fatal("expected parent to be parent of child")
	}
	if child.IsParentOf(parent) {
		t.Fatal("expected child not to be parent of parent")
	}
}

func TestQUrlMatches(t *testing.T) {
	a := QUrlFromString("https://example.com/path", TolerantMode)
	b := QUrlFromString("https://example.com/other", TolerantMode)
	defer a.Close()
	defer b.Close()
	if !a.Matches(b, FormattingOptions(RemovePath)) {
		t.Fatal("expected match with path removed")
	}
}

// ---- local files --------------------------------------------------------

func TestQUrlLocalFile(t *testing.T) {
	u := FromLocalFile("/tmp/file.txt")
	defer u.Close()
	if !u.IsLocalFile() {
		t.Fatal("expected local file")
	}
	if got := u.ToLocalFile(); got != "/tmp/file.txt" {
		t.Fatalf("toLocalFile: %q", got)
	}
}

// ---- static helpers -----------------------------------------------------

func TestQUrlFromEncoded(t *testing.T) {
	u := FromEncoded([]byte("https://example.com/a%20b"), TolerantMode)
	defer u.Close()
	if got := u.Path(FullyDecoded); got != "/a b" {
		t.Fatalf("path: %q", got)
	}
}

func TestQUrlFromUserInput(t *testing.T) {
	u := FromUserInputString("https://example.com/x", "", DefaultResolution)
	defer u.Close()
	if !u.IsValid() {
		t.Fatalf("fromUserInput invalid: %s", u.ErrorString())
	}
	if got := u.Host(FullyDecoded); got != "example.com" {
		t.Fatalf("host: %q", got)
	}
}

func TestQUrlPercentEncoding(t *testing.T) {
	enc := ToPercentEncodingBytes("héllo world", nil, nil)
	if !bytes.Contains(enc, []byte("%20")) {
		t.Fatalf("expected %%20 in %q", enc)
	}
	dec := FromPercentEncodingString(enc)
	if dec != "héllo world" {
		t.Fatalf("round-trip: %q", dec)
	}
}

func TestQUrlAce(t *testing.T) {
	enc := ToAceBytes("münchen.de")
	if len(enc) == 0 || !bytes.HasPrefix(enc, []byte("xn--")) {
		t.Fatalf("ToAce: %q", enc)
	}
	dec := FromAceString(enc)
	if dec != "münchen.de" {
		t.Fatalf("FromAce: %q", dec)
	}
}

func TestQUrlStringList(t *testing.T) {
	urls := []*QUrl{
		QUrlFromString("https://a.example/x", TolerantMode),
		QUrlFromString("https://b.example/y", TolerantMode),
	}
	defer func() {
		for _, u := range urls {
			u.Close()
		}
	}()

	strs := ToStringListStrings(urls, FormattingOptions(PrettyDecoded))
	if len(strs) != 2 {
		t.Fatalf("ToStringListStrings len: %d", len(strs))
	}
	if strs[0] != "https://a.example/x" {
		t.Fatalf("strs[0]: %q", strs[0])
	}

	back := FromStringListStrings(strs, TolerantMode)
	if len(back) != 2 {
		t.Fatalf("FromStringListStrings len: %d", len(back))
	}
	defer func() {
		for _, u := range back {
			u.Close()
		}
	}()
	if back[1].Host(FullyDecoded) != "b.example" {
		t.Fatalf("back[1] host: %q", back[1].Host(FullyDecoded))
	}
}

// ---- detach -------------------------------------------------------------

func TestQUrlDetach(t *testing.T) {
	u := QUrlFromString("https://example.com/", TolerantMode)
	defer u.Close()
	c := u.Clone()
	defer c.Close()
	// Both share the same QUrlPrivate initially; after detach on c, mutating
	// c must not affect u.
	c.Detach()
	c.SetPathString("/other", DecodedMode)
	if got := u.Path(FullyDecoded); got != "/" {
		t.Fatalf("original path mutated: %q", got)
	}
	if got := c.Path(FullyDecoded); got != "/other" {
		t.Fatalf("clone path: %q", got)
	}
}

// ---- QUrlQuery interop --------------------------------------------------

func TestQUrlSetQueryQUrlQuery(t *testing.T) {
	u := QUrlFromString("https://example.com/path", TolerantMode)
	defer u.Close()

	src := QUrlQueryFromString("a=1&b=2")
	defer src.Close()

	u.SetQueryQUrlQuery(src)
	if !u.HasQuery() {
		t.Fatal("expected has query")
	}

	// Round-trip through QUrlQuery(url) to check both keys.
	got := QUrlQueryFromUrl(u)
	defer got.Close()
	items := got.QueryItems(PrettyDecoded)
	if items["a"] != "1" || items["b"] != "2" {
		t.Fatalf("items: %v", items)
	}
}

func TestQUrlSetQueryQUrlQueryReplaces(t *testing.T) {
	u := QUrlFromString("https://example.com/path?old=1&keep=0", TolerantMode)
	defer u.Close()

	src := QUrlQueryFromString("new=2")
	defer src.Close()

	u.SetQueryQUrlQuery(src)

	got := QUrlQueryFromUrl(u)
	defer got.Close()
	items := got.QueryItems(PrettyDecoded)
	if _, ok := items["old"]; ok {
		t.Fatalf("old key still present: %v", items)
	}
	if items["new"] != "2" {
		t.Fatalf("new key: %v", items)
	}
	if _, ok := items["keep"]; ok {
		t.Fatalf("keep key still present: %v", items)
	}
}

func TestQUrlSetQueryQUrlQueryNil(t *testing.T) {
	u := QUrlFromString("https://example.com/path?a=1", TolerantMode)
	defer u.Close()
	if !u.HasQuery() {
		t.Fatal("precondition: expected query")
	}

	u.SetQueryQUrlQuery(nil)
	if u.HasQuery() {
		t.Fatalf("expected query cleared, got %q", u.Query(FullyDecoded))
	}
}

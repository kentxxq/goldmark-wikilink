package wikilink

import "path/filepath"

// DefaultResolver is a minimal wiklink resolver that resolves wikilinks
// relative to the source page.
//
// It adds ".html" to the end of the target
// if the target does not have an extension.
//
// For example,
//
//	[[Foo]]      // => "Foo.html"
//	[[Foo bar]]  // => "Foo bar.html"
//	[[foo/Bar]]  // => "foo/Bar.html"
//	[[foo.pdf]]  // => "foo.pdf"
//	[[foo.png]]  // => "foo.png"
var DefaultResolver Resolver = defaultResolver{}

// pretty url
// https://gohugo.io/content-management/urls/#appearance
//  [[Foo]]      // => "Foo/"
var PrettyResolver Resolver = prettyResolver{}

// when i use pretty url with [[relative path]]
//  /root/Foo.md                 url: /root/Foo/
//  /root/a.md include [[Foo]] . url: /root/a/    wikilink: /root/a/Foo/ not found!
//  so...
//  [[Foo]]      // => "../Foo/"    worked!
var RelResolver Resolver = relResolver{}

// when i use pretty url with [[absolute path]]
//  /Foo.md                      url: /posts/Foo/
//  /a.md include [[root/Foo]] . url: /posts/a/    wikilink: /posts/a/posts/Foo/ not found!
//  so...
//  [[Foo]]      // => "/root/Foo/" worked!
var RootResolver = func(b string) Resolver {
	return &rootResolver{
		base: b,
	}
}

// Resolver resolves pages referenced by wikilinks to their destinations.
type Resolver interface {
	// ResolveWikilink returns the address of the page that the provided
	// wikilink points to. The destination will be URL-escaped before
	// being placed into a link.
	//
	// If ResolveWikilink returns a non-nil error, rendering will be
	// halted.
	//
	// If ResolveWikilink returns a nil destination and error, the
	// Renderer will omit the link and render its contents as a regular
	// string.
	ResolveWikilink(*Node) (destination []byte, err error)
}

var _html = []byte(".html")

type defaultResolver struct{}

func (defaultResolver) ResolveWikilink(n *Node) ([]byte, error) {
	dest := make([]byte, len(n.Target)+len(_html)+len(_hash)+len(n.Fragment))
	var i int
	if len(n.Target) > 0 {
		i += copy(dest, n.Target)
		if filepath.Ext(string(n.Target)) == "" {
			i += copy(dest[i:], _html)
		}
	}
	if len(n.Fragment) > 0 {
		i += copy(dest[i:], _hash)
		i += copy(dest[i:], n.Fragment)
	}
	return dest[:i], nil
}

var pretty_html = []byte("/")

type prettyResolver struct{}

func (prettyResolver) ResolveWikilink(n *Node) ([]byte, error) {
	dest := make([]byte, len(n.Target)+len(pretty_html)+len(_hash)+len(n.Fragment))
	var i int
	if len(n.Target) > 0 {
		i += copy(dest, n.Target)
		if filepath.Ext(string(n.Target)) == "" {
			i += copy(dest[i:], pretty_html)
		}
	}
	if len(n.Fragment) > 0 {
		i += copy(dest[i:], _hash)
		i += copy(dest[i:], n.Fragment)
	}
	return dest[:i], nil
}

var rel_head = []byte("../")

type relResolver struct{}

func (relResolver) ResolveWikilink(n *Node) ([]byte, error) {
	dest := make([]byte, len(rel_head)+len(n.Target)+len(pretty_html)+len(_hash)+len(n.Fragment))
	var i int
	if len(n.Target) > 0 {
		i += copy(dest, rel_head)
		i += copy(dest[i:], n.Target)
		if filepath.Ext(string(n.Target)) == "" {
			i += copy(dest[i:], pretty_html)
		}
	}
	if len(n.Fragment) > 0 {
		i += copy(dest[i:], _hash)
		i += copy(dest[i:], n.Fragment)
	}
	return dest[:i], nil
}

type rootResolver struct {
	base string
}

func (r rootResolver) ResolveWikilink(n *Node) ([]byte, error) {
	dest := make([]byte, len(r.base)+len(n.Target)+len(pretty_html)+len(_hash)+len(n.Fragment))
	var i int
	if len(n.Target) > 0 {
		i += copy(dest, []byte(r.base))
		i += copy(dest[i:], n.Target)
		if filepath.Ext(string(n.Target)) == "" {
			i += copy(dest[i:], pretty_html)
		}
	}
	if len(n.Fragment) > 0 {
		i += copy(dest[i:], _hash)
		i += copy(dest[i:], n.Fragment)
	}
	return dest[:i], nil
}

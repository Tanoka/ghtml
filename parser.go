package ghtml

import (
	"golang.org/x/net/html"
	"strings"
)

// Busca un tag que tenga el attributo en el nodo indicado
func GetElement(tag string, attr string, id string, n *html.Node) (element *html.Node, ok bool) {
	if n.Type == html.ElementNode && n.Data == tag {
		//to get attr need to for all attr
		for _, a := range n.Attr {
			if a.Key == attr && strings.TrimSpace(a.Val) == strings.TrimSpace(id) {
				return n, true
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if element, ok = GetElement(tag, attr, id, c); ok {
			return
		}
	}
	return
}

// Obtiene todos los nodos que tengan el tag (+atributo y texto) indicado
func GetAllElement(tag string, attr string, id string, n *html.Node) (element []*html.Node, ok bool) {
	result := make([]*html.Node, 0)

	if n.Type == html.ElementNode && n.Data == tag {
		//to get attr need to for all attr
		for _, a := range n.Attr {
			if a.Key == attr && strings.TrimSpace(a.Val) == strings.TrimSpace(id) {
				result = append(result, n)
				return result, true
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if element, ok = GetAllElement(tag, attr, id, c); ok {
			result = append(result, element...)
		}
	}
	return result, len(result) > 0
}

//Busca un texto en un tag, o alguno de sus subnodos, que tenga un atributo con un valor indicado
func GetText(tag string, search AttrVal, n *html.Node) (dataret string, ok bool) {
	m, ok := GetElement(tag, search.attr, search.val, n)
	if ok {
		var getT func(n *html.Node) (dataret string, ok bool)

		getT = func(n *html.Node) (dataret string, ok bool) { //recursive function in a closure!, can't have name (can't be nested named functions), but can be asigned to a variable and be executed
			if n.Type == html.TextNode && len(strings.TrimSpace(n.Data)) > 0 { //Discart \n, spaces, etc in tags
				return n.Data, true
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				dataret, ok = getT(c)
				if ok { //only return when string found, otherwise continue with next tag
					return
				}
			}
			return
		}
		return getT(m)
	}
	return
}

// Busca el valor de un atributo en un tag que tenga otro atributo con un valor indicado
func GetAttr(tag string, search AttrVal, id string, n *html.Node) (dataret string, ok bool) {
	m, ok := GetElement(tag, search.attr, search.val, n)
	if ok {
		var getT func(tag string, attr string, n *html.Node) (dataret string, ok bool)

		getT = func(tag string, attr string, n *html.Node) (dataret string, ok bool) {
			if n.Type == html.ElementNode && n.Data == tag {
				for _, a := range n.Attr {
					if a.Key == attr {
						return a.Val, true
					}
				}

			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				dataret, ok = getT(tag, attr, c)
				return
			}
			return
		}

		return getT(tag, id, m)
	}
	return
}

// Busca una cadena de texto y retorna el texto desde el final de esa cadena hasta el inicio de otra cadena final indicada.
func GetMidValue(text string, keyWord string, endChar string) string {
	var pos, fin int
	if pos = strings.Index(text, keyWord); pos == -1 {
		return ""
	}
	st := pos + len(keyWord)
	if fin = strings.Index(text[st:], endChar); fin == -1 {
		return ""
	}
	en := st + fin
	return text[st:en]
}

//Tipo para tuplas {atributo, valor}
type AttrVal struct {
	attr string
	val  string
}

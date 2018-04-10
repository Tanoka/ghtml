package ghtml

import (
	"bytes"
	"golang.org/x/net/html"
	"testing"
)

func getNode(st string) *html.Node {
	nodeHtml := bytes.NewReader([]byte(st))
	res, _ := html.Parse(nodeHtml)
	return res
}

func TestGetText(t *testing.T) {
	shtml := "<html><div class='classdiv'><span class='classspan'>uno</span> </div> </html>"

	node := getNode(shtml)

	var flagtests = []struct {
		tag string
		att AttrVal
		nod *html.Node
		res string
		ok  bool
	}{
		{"div", AttrVal{"class", "classdiv"}, node, "uno", true},
		{"div", AttrVal{"clas", "classdiv"}, node, "", false},
		{"div", AttrVal{"class", "clas"}, node, "", false},
		{"span", AttrVal{"class", "classspan"}, node, "uno", true},
		{"span", AttrVal{"class", "classdiv"}, node, "", false},
		{"tr", AttrVal{"class", "classdiv"}, node, "", false},
	}

	for _, res := range flagtests {
		val, isok := GetText(res.tag, res.att, res.nod)
		if val != res.res || isok != res.ok {
			t.Errorf("Not equals. Params:(%s, %s, %s).  %s != %s or %v != %v", res.tag, res.att.attr, res.att.val, res.res, val, res.ok, isok)
		}
	}

}

func TestGetAttr(t *testing.T) {
	shtml := "<html><div class='classdiv'><span class='classspan' data='azul' >uno</span> </div> </html>"

	node := getNode(shtml)

	var flagtests = []struct {
		tag string
		att AttrVal
		id  string
		nod *html.Node
		res string
		ok  bool
	}{
		{"div", AttrVal{"class", "classdiv"}, "data", node, "", false},
		{"div", AttrVal{"clas", "classdiv"}, "i", node, "", false},
		{"div", AttrVal{"class", "clas"}, "ii", node, "", false},
		{"span", AttrVal{"class", "classspan"}, "data", node, "azul", true},
		{"span", AttrVal{"class", "classspan"}, "info", node, "", false},
		{"span", AttrVal{"class", "classdiv"}, "data", node, "", false},
		{"tr", AttrVal{"class", "classdiv"}, "data", node, "", false},
	}

	for _, res := range flagtests {
		val, isok := GetAttr(res.tag, res.att, res.id, res.nod)
		if val != res.res || isok != res.ok {
			t.Errorf("Not equals. Params:(%s, %s, %s, %s).  %s != %s or %v != %v", res.tag, res.att.attr, res.att.val, res.id, res.res, val, res.ok, isok)
		}
	}

}

func TestGetMidValue(t *testing.T) {

	var shtml string
	shtml = "<html><div class='classdiv'><span class='classspan' data='azul' >uÑo</span> </div> \n"
	shtml += "<div class='classdiv'>Texto: \"Valor\"</div>Fácil texto fin</html>"

	var flagtest = []struct {
		word string
		end  string
		res  string
	}{
		{"nada", "lor", ""},
		{"Texto: \"", "\"", "Valor"},
		{"Rexto;", "fin", ""},
		{"Fácil", "fin", " texto "},
		{"Fá", " fin", "cil texto"},
	}

	for _, restest := range flagtest {
		val := GetMidValue(shtml, restest.word, restest.end)
		if val != restest.res {
			t.Errorf("Not equals. Params:(%s, %s).  Expected:%s != Actual:%s ", restest.word, restest.end, restest.res, val)
		}

	}
}


func TestGetAllElement(t *testing.T) {

	html := "<html><div class='classdiv'><span class='classspan'>uno</span> </div> "
	html += "<div class='classdiv'><span class='classspan'>uno</span> </div>"
	html += "<div class='classdiv'><span class='classspan'>dos</span> </div>"
	html += "<table><tr><td>Hola</td></tr></table>"
	html += "<div class='classdiv2'><span class='classspan'>tres</span> </div> \n"
	html += "<div class='classdiv'><span class='classspan'>cuatro</span> </div>"
	html += "</html>"

	var flagtests = []struct{
		tag string
		att string
		val string
		ok  bool
		cua int
	}{
		{"div","class","no",false, 0},
		{"tr","class","classdiv",false, 0},
		{"div","class","classdiv",true, 4},
		{"div","","",false, 0},
		{"div","class","classdiv2",true, 1},
	}

	root := getNode(html)

	for _, flagte := range flagtests {
		respo, oki := GetAllElement(flagte.tag, flagte.att, flagte.val, root)
		if (oki != flagte.ok) {
			t.Errorf("Not equals. Params:(%s, %s).  Expected:%t != Actual:%t ",flagte.tag, flagte.att, flagte.ok, oki)
		}
		if (len(respo) != flagte.cua) {
			t.Errorf("Not equals. Params:(%s, %s).  Expected:%d != Actual:%d ",flagte.tag, flagte.att, flagte.cua, len(respo))
		}
	}
}


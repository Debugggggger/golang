// Copyright 2015 Light Code Labs, LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metadata

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

var TOML = [5]string{`
title = "A title"
template = "default"
name = "value"
positive = true
negative = false
number = 1410
float = 1410.07
`,
	`+++
title = "A title"
template = "default"
name = "value"
positive = true
negative = false
number = 1410
float = 1410.07
+++
Page content
	`,
	`+++
title = "A title"
template = "default"
name = "value"
positive = true
negative = false
number = 1410
float = 1410.07
	`,
	`title = "A title" template = "default" [variables] name = "value"`,
	`+++
title = "A title"
template = "default"
name = "value"
positive = true
negative = false
number = 1410
float = 1410.07
+++
`,
}

var YAML = [5]string{`
title : A title
template : default
name : value
positive : true
negative : false
number : 1410
float : 1410.07
`,
	`---
title : A title
template : default
name : value
positive : true
negative : false
number : 1410
float : 1410.07
---
	Page content
	`,
	`---
title : A title
template : default
name : value
number : 1410
float : 1410.07
	`,
	`title : A title template : default variables : name : value : positive : true : negative : false`,
	`---
title : A title
template : default
name : value
positive : true
negative : false
number : 1410
float : 1410.07
---
`,
}

var JSON = [5]string{`
	"title" : "A title",
	"template" : "default",
	"name" : "value",
	"positive" : true,
	"negative" : false,
	"number": 1410,
	"float": 1410.07
`,
	`{
	"title" : "A title",
	"template" : "default",
	"name" : "value",
	"positive" : true,
	"negative" : false,
	"number" : 1410,
	"float": 1410.07
}
Page content
	`,
	`
{
	"title" : "A title",
	"template" : "default",
	"name" : "value",
	"positive" : true,
	"negative" : false,
	"number" : 1410,
	"float": 1410.07
	`,
	`
{
	"title" :: "A title",
	"template" : "default",
	"name" : "value",
	"positive" : true,
	"negative" : false,
	"number" : 1410,
	"float": 1410.07
}
	`,
	`{
	"title" : "A title",
	"template" : "default",
	"name" : "value",
	"positive" : true,
	"negative" : false,
	"number" : 1410,
	"float": 1410.07
}
`,
}

func TestParsers(t *testing.T) {
	expected := Metadata{
		Title:    "A title",
		Template: "default",
		Variables: map[string]interface{}{
			"name":     "value",
			"title":    "A title",
			"template": "default",
			"number":   1410,
			"float":    1410.07,
			"positive": true,
			"negative": false,
		},
	}
	compare := func(m Metadata) bool {
		if m.Title != expected.Title {
			return false
		}
		if m.Template != expected.Template {
			return false
		}
		for k, v := range m.Variables {
			if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", expected.Variables[k]) {
				return false
			}
		}

		varLenOK := len(m.Variables) == len(expected.Variables)
		return varLenOK
	}

	data := []struct {
		parser   Parser
		testData [5]string
		name     string
	}{
		{&JSONParser{}, JSON, "JSON"},
		{&YAMLParser{}, YAML, "YAML"},
		{&TOMLParser{}, TOML, "TOML"},
	}

	for _, v := range data {
		// metadata without identifiers
		if v.parser.Init(bytes.NewBufferString(v.testData[0])) {
			t.Fatalf("Expected error for invalid metadata for %v", v.name)
		}

		// metadata with identifiers
		if !v.parser.Init(bytes.NewBufferString(v.testData[1])) {
			t.Fatalf("Metadata failed to initialize, type %v", v.parser.Type())
		}
		md := v.parser.Markdown()
		if !compare(v.parser.Metadata()) {
			t.Fatalf("Expected %v, found %v for %v", expected, v.parser.Metadata(), v.name)
		}
		if "Page content" != strings.TrimSpace(string(md)) {
			t.Fatalf("Expected %v, found %v for %v", "Page content", string(md), v.name)
		}
		// Check that we find the correct metadata parser type
		if p := GetParser([]byte(v.testData[1])); p.Type() != v.name {
			t.Fatalf("Wrong parser found, expected %v, found %v", v.name, p.Type())
		}

		// metadata without closing identifier
		if v.parser.Init(bytes.NewBufferString(v.testData[2])) {
			t.Fatalf("Expected error for missing closing identifier for %v parser", v.name)
		}

		// invalid metadata
		if v.parser.Init(bytes.NewBufferString(v.testData[3])) {
			t.Fatalf("Expected error for invalid metadata for %v", v.name)
		}

		// front matter but no body
		if !v.parser.Init(bytes.NewBufferString(v.testData[4])) {
			t.Fatalf("Unexpected error for valid metadata but no body for %v", v.name)
		}
	}
}

func TestLargeBody(t *testing.T) {

	var JSON = `{
"template": "chapter"
}

Mycket olika byggnader har man i de nordiska rikena: pyramidformiga, kilformiga, v??lvda, runda och fyrkantiga. De pyramidformiga best??r helt enkelt av tr??ribbor, som upptill l??per samman och nedtill bildar en vidare krets; de ??r avsedda att anv??ndas av hantverkarna under sommaren, f??r att de inte ska pl??gas av solen, p?? samma g??ng som de besv??ras av r??k och eld. De kilformiga husen ??r i regel f??rsedda med h??ga tak, f??r att de t??ta och tunga sn??massorna fortare ska kunna bl??sa av och inte tynga ned taken. Dessa ??r t??ckta av bj??rkn??ver, tegel eller kluvet sp??n av furu - f??r k??dans skull -, gran, ek eller bok; taken p?? de f??rm??gnas hus d??remot med pl??tar av koppar eller bly, i likhet med kyrktaken. Valvbyggnaderna uppf??rs ganska konstn??rligt till skydd mot v??ldsamma vindar och sn??fall, g??rs av sten eller tr??, och ??r avsedda f??r olika alldagliga viktiga ??ndam??l. Liknande byggnader kan finnas i storm??nnens g??rdar d??r de anv??nds som f??rvaringsrum f??r husger??d och jordbruksredskap. De runda byggnaderna - som f??r ??vrigt ??r de h??gst s??llsynta - anv??nds av konstn??rer, som vid sitt arbete beh??ver ett j??mnt f??rdelat ljus fr??n taket. Vanligast ??r de fyrkantiga husen, vars grova bj??lkar ??r synnerligen v??l hopfogade i h??rnen - ett sant m??sterverk av byggnadskonst; ??ven dessa har f??nster h??gt uppe i taken, f??r att dagsljuset skall kunna str??mma in och ge alla d??rinne full belysning. Stenhusen har d??rr??ppningar i f??rh??llande till byggnadens storlek, men smala f??nstergluggar, som skydd mot den str??nga k??lden, frosten och sn??n. Vore de st??rre och vidare, s??som f??nstren i Italien, skulle husen i f??ljd av den fint yrande sn??n, som r??res upp av den starka bl??sten, precis som dammet av virvelvinden, snart nog fyllas med massor av sn?? och inte kunna st?? emot dess tryck, utan st??rta samman.

	`
	var TOML = `+++
template = "chapter"
+++

Mycket olika byggnader har man i de nordiska rikena: pyramidformiga, kilformiga, v??lvda, runda och fyrkantiga. De pyramidformiga best??r helt enkelt av tr??ribbor, som upptill l??per samman och nedtill bildar en vidare krets; de ??r avsedda att anv??ndas av hantverkarna under sommaren, f??r att de inte ska pl??gas av solen, p?? samma g??ng som de besv??ras av r??k och eld. De kilformiga husen ??r i regel f??rsedda med h??ga tak, f??r att de t??ta och tunga sn??massorna fortare ska kunna bl??sa av och inte tynga ned taken. Dessa ??r t??ckta av bj??rkn??ver, tegel eller kluvet sp??n av furu - f??r k??dans skull -, gran, ek eller bok; taken p?? de f??rm??gnas hus d??remot med pl??tar av koppar eller bly, i likhet med kyrktaken. Valvbyggnaderna uppf??rs ganska konstn??rligt till skydd mot v??ldsamma vindar och sn??fall, g??rs av sten eller tr??, och ??r avsedda f??r olika alldagliga viktiga ??ndam??l. Liknande byggnader kan finnas i storm??nnens g??rdar d??r de anv??nds som f??rvaringsrum f??r husger??d och jordbruksredskap. De runda byggnaderna - som f??r ??vrigt ??r de h??gst s??llsynta - anv??nds av konstn??rer, som vid sitt arbete beh??ver ett j??mnt f??rdelat ljus fr??n taket. Vanligast ??r de fyrkantiga husen, vars grova bj??lkar ??r synnerligen v??l hopfogade i h??rnen - ett sant m??sterverk av byggnadskonst; ??ven dessa har f??nster h??gt uppe i taken, f??r att dagsljuset skall kunna str??mma in och ge alla d??rinne full belysning. Stenhusen har d??rr??ppningar i f??rh??llande till byggnadens storlek, men smala f??nstergluggar, som skydd mot den str??nga k??lden, frosten och sn??n. Vore de st??rre och vidare, s??som f??nstren i Italien, skulle husen i f??ljd av den fint yrande sn??n, som r??res upp av den starka bl??sten, precis som dammet av virvelvinden, snart nog fyllas med massor av sn?? och inte kunna st?? emot dess tryck, utan st??rta samman.

	`
	var YAML = `---
template : chapter
---

Mycket olika byggnader har man i de nordiska rikena: pyramidformiga, kilformiga, v??lvda, runda och fyrkantiga. De pyramidformiga best??r helt enkelt av tr??ribbor, som upptill l??per samman och nedtill bildar en vidare krets; de ??r avsedda att anv??ndas av hantverkarna under sommaren, f??r att de inte ska pl??gas av solen, p?? samma g??ng som de besv??ras av r??k och eld. De kilformiga husen ??r i regel f??rsedda med h??ga tak, f??r att de t??ta och tunga sn??massorna fortare ska kunna bl??sa av och inte tynga ned taken. Dessa ??r t??ckta av bj??rkn??ver, tegel eller kluvet sp??n av furu - f??r k??dans skull -, gran, ek eller bok; taken p?? de f??rm??gnas hus d??remot med pl??tar av koppar eller bly, i likhet med kyrktaken. Valvbyggnaderna uppf??rs ganska konstn??rligt till skydd mot v??ldsamma vindar och sn??fall, g??rs av sten eller tr??, och ??r avsedda f??r olika alldagliga viktiga ??ndam??l. Liknande byggnader kan finnas i storm??nnens g??rdar d??r de anv??nds som f??rvaringsrum f??r husger??d och jordbruksredskap. De runda byggnaderna - som f??r ??vrigt ??r de h??gst s??llsynta - anv??nds av konstn??rer, som vid sitt arbete beh??ver ett j??mnt f??rdelat ljus fr??n taket. Vanligast ??r de fyrkantiga husen, vars grova bj??lkar ??r synnerligen v??l hopfogade i h??rnen - ett sant m??sterverk av byggnadskonst; ??ven dessa har f??nster h??gt uppe i taken, f??r att dagsljuset skall kunna str??mma in och ge alla d??rinne full belysning. Stenhusen har d??rr??ppningar i f??rh??llande till byggnadens storlek, men smala f??nstergluggar, som skydd mot den str??nga k??lden, frosten och sn??n. Vore de st??rre och vidare, s??som f??nstren i Italien, skulle husen i f??ljd av den fint yrande sn??n, som r??res upp av den starka bl??sten, precis som dammet av virvelvinden, snart nog fyllas med massor av sn?? och inte kunna st?? emot dess tryck, utan st??rta samman.

	`
	var NONE = `

Mycket olika byggnader har man i de nordiska rikena: pyramidformiga, kilformiga, v??lvda, runda och fyrkantiga. De pyramidformiga best??r helt enkelt av tr??ribbor, som upptill l??per samman och nedtill bildar en vidare krets; de ??r avsedda att anv??ndas av hantverkarna under sommaren, f??r att de inte ska pl??gas av solen, p?? samma g??ng som de besv??ras av r??k och eld. De kilformiga husen ??r i regel f??rsedda med h??ga tak, f??r att de t??ta och tunga sn??massorna fortare ska kunna bl??sa av och inte tynga ned taken. Dessa ??r t??ckta av bj??rkn??ver, tegel eller kluvet sp??n av furu - f??r k??dans skull -, gran, ek eller bok; taken p?? de f??rm??gnas hus d??remot med pl??tar av koppar eller bly, i likhet med kyrktaken. Valvbyggnaderna uppf??rs ganska konstn??rligt till skydd mot v??ldsamma vindar och sn??fall, g??rs av sten eller tr??, och ??r avsedda f??r olika alldagliga viktiga ??ndam??l. Liknande byggnader kan finnas i storm??nnens g??rdar d??r de anv??nds som f??rvaringsrum f??r husger??d och jordbruksredskap. De runda byggnaderna - som f??r ??vrigt ??r de h??gst s??llsynta - anv??nds av konstn??rer, som vid sitt arbete beh??ver ett j??mnt f??rdelat ljus fr??n taket. Vanligast ??r de fyrkantiga husen, vars grova bj??lkar ??r synnerligen v??l hopfogade i h??rnen - ett sant m??sterverk av byggnadskonst; ??ven dessa har f??nster h??gt uppe i taken, f??r att dagsljuset skall kunna str??mma in och ge alla d??rinne full belysning. Stenhusen har d??rr??ppningar i f??rh??llande till byggnadens storlek, men smala f??nstergluggar, som skydd mot den str??nga k??lden, frosten och sn??n. Vore de st??rre och vidare, s??som f??nstren i Italien, skulle husen i f??ljd av den fint yrande sn??n, som r??res upp av den starka bl??sten, precis som dammet av virvelvinden, snart nog fyllas med massor av sn?? och inte kunna st?? emot dess tryck, utan st??rta samman.

	`
	var expectedBody = `Mycket olika byggnader har man i de nordiska rikena: pyramidformiga, kilformiga, v??lvda, runda och fyrkantiga. De pyramidformiga best??r helt enkelt av tr??ribbor, som upptill l??per samman och nedtill bildar en vidare krets; de ??r avsedda att anv??ndas av hantverkarna under sommaren, f??r att de inte ska pl??gas av solen, p?? samma g??ng som de besv??ras av r??k och eld. De kilformiga husen ??r i regel f??rsedda med h??ga tak, f??r att de t??ta och tunga sn??massorna fortare ska kunna bl??sa av och inte tynga ned taken. Dessa ??r t??ckta av bj??rkn??ver, tegel eller kluvet sp??n av furu - f??r k??dans skull -, gran, ek eller bok; taken p?? de f??rm??gnas hus d??remot med pl??tar av koppar eller bly, i likhet med kyrktaken. Valvbyggnaderna uppf??rs ganska konstn??rligt till skydd mot v??ldsamma vindar och sn??fall, g??rs av sten eller tr??, och ??r avsedda f??r olika alldagliga viktiga ??ndam??l. Liknande byggnader kan finnas i storm??nnens g??rdar d??r de anv??nds som f??rvaringsrum f??r husger??d och jordbruksredskap. De runda byggnaderna - som f??r ??vrigt ??r de h??gst s??llsynta - anv??nds av konstn??rer, som vid sitt arbete beh??ver ett j??mnt f??rdelat ljus fr??n taket. Vanligast ??r de fyrkantiga husen, vars grova bj??lkar ??r synnerligen v??l hopfogade i h??rnen - ett sant m??sterverk av byggnadskonst; ??ven dessa har f??nster h??gt uppe i taken, f??r att dagsljuset skall kunna str??mma in och ge alla d??rinne full belysning. Stenhusen har d??rr??ppningar i f??rh??llande till byggnadens storlek, men smala f??nstergluggar, som skydd mot den str??nga k??lden, frosten och sn??n. Vore de st??rre och vidare, s??som f??nstren i Italien, skulle husen i f??ljd av den fint yrande sn??n, som r??res upp av den starka bl??sten, precis som dammet av virvelvinden, snart nog fyllas med massor av sn?? och inte kunna st?? emot dess tryck, utan st??rta samman.
`

	data := []struct {
		pType    string
		testData string
	}{
		{"JSON", JSON},
		{"TOML", TOML},
		{"YAML", YAML},
		{"None", NONE},
	}
	for _, v := range data {
		p := GetParser([]byte(v.testData))
		if v.pType != p.Type() {
			t.Fatalf("Wrong parser type, expected %v, got %v", v.pType, p.Type())
		}
		md := p.Markdown()
		if strings.TrimSpace(string(md)) != strings.TrimSpace(expectedBody) {
			t.Log("Provided:", v.testData)
			t.Log("Returned:", p.Markdown())
			t.Fatalf("Error, mismatched body in expected type %v, matched type %v", v.pType, p.Type())
		}
	}
}

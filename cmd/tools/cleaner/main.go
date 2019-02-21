package main

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

type alias int

var fset = token.NewFileSet()

func (
	a *alias,
) String() string {
	return ""
	// fmt.Sprintln(*a)
}

func liner(bb []byte) []byte {

	// first find the beginnings of each function in the source
	re := regexp.MustCompile("\nfunc ")
	fi := re.FindAllIndex(bb, -1)
	for _, x := range fi {
	again:
		if bb[x[1]] == '(' {
			if bb[x[1]+1] == '\n' &&
				bb[x[1]+2] == '\t' {
				rer := regexp.MustCompile(
					"([_a-zA-Z0-9]*[ ]?[*_a-zA-Z0-9]+)")
				rest := bb[x[1]+3:]
				fr := rer.FindIndex(rest)
				if rest[fr[1]] == ',' {
					if rest[fr[1]+1] == '\n' &&
						rest[fr[1]+2] == ')' {
						goto step1
					}
				}
				if rest[fr[1]] == ')' {
					bb = append(bb[:x[1]+3+fr[1]], append([]byte{',', '\n'}, bb[x[1]+3+fr[1]:]...)...)
					goto step1
				}
			} else {
				bb = append(bb[:x[1]+1],
					append([]byte{'\n', '\t'}, bb[x[1]+1:]...)...)
				goto again
			}
		step1:
		} else {
			rest := bb[x[1]:]
			re := regexp.MustCompile("[_a-zA-Z][._a-zA-Z0-9]*")
			ff := re.FindIndex(rest)
			fmt.Println(string(rest[ff[1] : ff[1]+10]))
		}
	}

	return bb
}

func main() {

	bb, e := ioutil.ReadFile(os.Args[1])
	if e != nil {

		panic(e)
	}
	ss := string(bb)
	bb = sorter(ss)
	bb = liner(bb)
	if e := ioutil.WriteFile(os.Args[1], bb, 0644); e != nil {
		panic(e)
	}
}

func sorter(ss string) []byte {

	splittedraw := strings.Split(ss, "\n")
	imports := []string{}
	for i, x := range splittedraw {

		if x == "import (" {

			impfound := false
			for j := i; !impfound; j++ {

				imports = append(imports, splittedraw[j])
				if splittedraw[j] == ")" {

					goto imported
				}
			}
		}
	}

imported:
	file, err := decorator.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	constcounter := 0
	unsortedDecls := make(map[string]dst.Decl)
	for _, decl := range file.Decls {

		switch gen := decl.(type) {

		case *dst.GenDecl:

			switch gen.Tok {

			case token.CONST:
				already := false
				for _, x := range unsortedDecls {

					if decl == x {
						already = true
					}
				}

				if already {

					continue
				}
				unsortedDecls["2const"+fmt.Sprint(constcounter)] = decl
				constcounter++

			case token.VAR:
				for _, y := range gen.Specs {

					if is, ok := y.(*dst.ValueSpec); ok {

						var declName string
						key := "3:"
						for i, z := range is.Names {

							key += fmt.Sprint(z)
							_ = i
							declName = z.Name
						}

						d := dst.Clone(decl).(*dst.GenDecl)
						for i, x := range decl.(*dst.GenDecl).Specs {

							dd := d.Specs[i].(*dst.ValueSpec)
							if dd.Names[0].Name == declName {

								d.Specs = nil
								d.Specs = append(d.Specs, x)
								break
							}
						}

						unsortedDecls[key] = d
					}
				}

			case token.TYPE:
				for _, y := range gen.Specs {

					if is, ok := y.(*dst.TypeSpec); ok {

						var declName string
						key := "1:"
						declName = is.Name.Name
						d := dst.Clone(decl).(*dst.GenDecl)
						for i, x := range decl.(*dst.GenDecl).Specs {

							dd := d.Specs[i].(*dst.TypeSpec)
							if dd.Name.Name == declName {

								key += dd.Name.Name
								d.Specs = nil
								d.Specs = append(d.Specs, x)
								break
							}
						}
						unsortedDecls[key] = d
					}
				}
			}

		case *dst.FuncDecl:
			if fun, ok := decl.(*dst.FuncDecl); ok {

				if fun.Recv != nil {

					for _, x := range fun.Recv.List {

						var typename, star string
						if rt, ok := x.Type.(*dst.StarExpr); ok {

							star = "*"
							typename = fmt.Sprint(rt.X)
						} else {

							typename = fmt.Sprint(x.Type)
						}
						for range x.Names {

							unsortedDecls[fmt.Sprintf("4:(%s%s) %s", star, typename, fun.Name.Name)] = decl
						}
					}
				} else {

					unsortedDecls[fmt.Sprint("4:", fun.Name.Name)] = decl
				}
			}
		}
	}
	var sortedDecls []string
	for i := range unsortedDecls {

		sortedDecls = append(sortedDecls, i)
	}

	var decls []dst.Decl
	sort.Strings(sortedDecls)
	for _, x := range sortedDecls {

		decls = append(decls, unsortedDecls[x])
	}

	file.Decls = decls
	var buf bytes.Buffer
	decorator.Fprint(&buf, file)
	output := string(buf.Bytes())
	var splitout []string
	splitted := strings.Split(output, "\n")
	packagefound := false
	if len(imports) < 1 {

		packagefound = true
	}
	for _, x := range splitted {

		splitout = append(splitout, x)
		if !packagefound {

			sss := strings.Split(x, " ")
			if sss[0] == "package" {

				splitout = append(splitout, "")

				splitout = append(splitout, imports...)
			}
		}
	}
	joined := []byte(strings.Join(splitout, "\n"))
	return joined
}

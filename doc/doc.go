package doc

import (
	"fmt"
	baseStrings "strings"

	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/HomayoonAlimohammadi/protoc-gen-doc/comment"
	"github.com/HomayoonAlimohammadi/protoc-gen-doc/strings"
)

type Doc struct {
	widgetDocs map[string]*WidgetDoc
}

type WidgetDoc struct {
	name      string
	comm      *comment.WidgetStructure
	fieldsDoc []*WidgetFieldDoc
}

type WidgetFieldDoc struct {
	name string
	typ  string
	comm *comment.WidgetFieldStructure
}

func Generate(gen *protogen.Plugin, exclude string) (*Doc, error) {
	gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	var doc Doc
	doc.widgetDocs = make(map[string]*WidgetDoc)

	for _, f := range gen.Files {

		// extract widget type from widgets.proto
		// widget types are available in a nested structure like:
		// widgets.proto -> message Widget -> enum Type
		if f.Proto.GetName() == "divar_interface/widgets/widgets.proto" {
			for _, m := range f.Messages {
				if m.Desc.Name() == "Widget" {
					for _, e := range m.Enums {
						if e.Desc.Name() == "Type" {
							for _, v := range e.Values {
								name := strings.NormalizeWidgetType(string(v.Desc.Name()))

								comm, err := comment.ParseWidget(string(v.Comments.Leading))
								if err != nil {
									return nil, fmt.Errorf("%s: %w", v.Desc.Name(), err)
								}

								// add to previous WidgetDoc if available, otherwise create new
								d := doc.widgetDocs[name]
								if d != nil {
									d.name = string(v.Desc.Name())
									d.comm = comm
								} else {
									d := &WidgetDoc{
										name: string(v.Desc.Name()),
										comm: comm,
									}
									doc.widgetDocs[name] = d
								}
							}
						}
					}
				}
			}
		} else if f.Proto.GetName() == "divar_interface/widgets/widgets_data.proto" {
			// extract widgets data messages from widget_data.proto
			for _, m := range f.Messages {
				name := strings.NormalizeWidgetDataMessageName(string(m.Desc.Name()))

				fieldsDoc := []*WidgetFieldDoc{}
				for _, f := range m.Fields {
					comm, err := comment.ParseWidgetField(string(f.Comments.Leading))
					if err != nil {
						return nil, fmt.Errorf("%s: %w", f.Desc.Name(), err)
					}
					fd := &WidgetFieldDoc{
						name: string(f.Desc.Name()),
						typ:  fmt.Sprint(f.Desc.Kind()),
						comm: comm,
					}
					fieldsDoc = append(fieldsDoc, fd)
				}

				// add to previous WidgetDoc if available, otherwise create new
				d := doc.widgetDocs[name]
				if d != nil {
					d.fieldsDoc = fieldsDoc
				} else {
					d := &WidgetDoc{
						fieldsDoc: fieldsDoc,
					}
					doc.widgetDocs[name] = d
				}

			}
		}
	}

	// All widgets that are excluded from doc validation should be listed in a
	// dash-separated fashion, like: SELECTOR_ROW-MY_RANDOM_WIDGET-SOMETHING_ELSE
	// This will make 'SELECTOR_ROW' , 'MY_RANDOM_WIDGET' and 'SOMETHING_ELSE' to
	// bypass validation. Note that these names SHOULD be formatted like the widget
	// 'Type' enum inside widgets.proto
	//
	// If __ALL__ is available in the values, all widgets will be excluded from verification.
	excludedWidgets := baseStrings.Split(exclude, "-")

	// __ALL__ excludes all fieldsDocs to not get checked
	if slices.Contains(excludedWidgets, "__ALL__") {
		return &doc, nil
	}

	// all widgets should have fields documentation except some specific ones
	for name, d := range doc.widgetDocs {
		if d.fieldsDoc == nil && !slices.Contains(excludedWidgets, d.name) {
			return nil, fmt.Errorf("%w: %s", ErrNoFields, name)
		}
	}

	return &doc, nil
}

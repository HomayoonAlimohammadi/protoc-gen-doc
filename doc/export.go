package doc

import (
	"strconv"

	"google.golang.org/protobuf/compiler/protogen"
)

func Export(doc *Doc, gen *protogen.Plugin) error {
	g := gen.NewGeneratedFile("widget-documentation.md", "")

	for _, wd := range doc.widgetDocs {

		// name
		g.P("# ", wd.name)

		// support
		g.P("## Client Support")
		g.P("| IOS | Android | Web |")
		g.P("|-----|---------|-----|")
		var iosMaxVersion, androidMaxVersion, webSupportSymbol, androidSupportSymbol, iosSupportSymbol string

		// web
		if wd.comm.PlatformSupport.IsWebSupported {
			webSupportSymbol = ":white_check_mark:"
		} else {
			webSupportSymbol = ":x:"
		}

		// ios
		if wd.comm.PlatformSupport.IOSMaxVersion == 0 {
			iosMaxVersion = "Now"
		} else {
			iosMaxVersion = strconv.Itoa(wd.comm.PlatformSupport.IOSMaxVersion)
		}
		if wd.comm.PlatformSupport.IOSMaxVersion == 0 && wd.comm.PlatformSupport.IOSMinVersion == 0 {
			iosSupportSymbol = ":x:"
		} else {
			iosSupportSymbol = strconv.Itoa(wd.comm.PlatformSupport.IOSMinVersion) + " &rarr; " + iosMaxVersion
		}

		// android
		if wd.comm.PlatformSupport.AndroidMaxVersion == 0 {
			androidMaxVersion = "Now"
		} else {
			androidMaxVersion = strconv.Itoa(wd.comm.PlatformSupport.AndroidMaxVersion)
		}
		if wd.comm.PlatformSupport.AndroidMaxVersion == 0 && wd.comm.PlatformSupport.AndroidMinVersion == 0 {
			androidSupportSymbol = ":x:"
		} else {
			androidSupportSymbol = strconv.Itoa(wd.comm.PlatformSupport.AndroidMinVersion) + " &rarr; " + androidMaxVersion
		}

		g.P("|", iosSupportSymbol, "|",
			androidSupportSymbol, "|",
			webSupportSymbol, "|")

		// design
		g.P("## Design URL")
		g.P("<li><a>", wd.comm.DesignURL, "</a></li>")
		g.P("<br>")
		g.P()

		// widgetify
		g.P("## Widgetify URL")
		g.P("<li><a>", wd.comm.WidgetifyUrl, "</a></li>")
		g.P("<br>")
		g.P()

		// fields
		g.P("## Fields")
		g.P("| Name | Type | IOS | Android | Web |")
		g.P("|------|------|-----|---------|-----|")
		for _, fd := range wd.fieldsDoc {

			var iosMaxVersion, androidMaxVersion, webSupportSymbol, androidSupportSymbol, iosSupportSymbol string

			// web
			if fd.comm.PlatformSupport.IsWebSupported {
				webSupportSymbol = ":white_check_mark:"
			} else {
				webSupportSymbol = ":x:"
			}

			// ios
			if fd.comm.PlatformSupport.IOSMaxVersion == 0 {
				iosMaxVersion = "Now"
			} else {
				iosMaxVersion = strconv.Itoa(fd.comm.PlatformSupport.IOSMaxVersion)
			}
			if fd.comm.PlatformSupport.IOSMaxVersion == 0 && fd.comm.PlatformSupport.IOSMinVersion == 0 {
				iosSupportSymbol = ":x:"
			} else {
				iosSupportSymbol = strconv.Itoa(fd.comm.PlatformSupport.IOSMinVersion) + " &rarr; " + iosMaxVersion
			}

			// android
			if fd.comm.PlatformSupport.AndroidMaxVersion == 0 {
				androidMaxVersion = "Now"
			} else {
				androidMaxVersion = strconv.Itoa(fd.comm.PlatformSupport.AndroidMaxVersion)
			}
			if fd.comm.PlatformSupport.AndroidMaxVersion == 0 && fd.comm.PlatformSupport.AndroidMinVersion == 0 {
				androidSupportSymbol = ":x:"
			} else {
				androidSupportSymbol = strconv.Itoa(fd.comm.PlatformSupport.AndroidMinVersion) + " &rarr; " + androidMaxVersion
			}

			g.P("|", fd.name, "|", fd.typ, "|",
				iosSupportSymbol, "|",
				androidSupportSymbol, "|",
				webSupportSymbol, "|")
		}
		g.P("<br>")
		g.P()

	}

	return nil
}

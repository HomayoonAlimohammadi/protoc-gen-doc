package comment

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

type WidgetStructure struct {
	PlatformSupport *PlatformSupportStructure
	DesignURL       *url.URL
	WidgetifyUrl    *url.URL
}

type WidgetFieldStructure struct {
	PlatformSupport *PlatformSupportStructure
}

type PlatformSupportStructure struct {
	AndroidMinVersion, AndroidMaxVersion int
	IOSMinVersion, IOSMaxVersion         int

	IsWebSupported bool
}

// ParseWidget parses specific widget's comment from protofile.
// The comment should have @support and @design cmds.
// For more information on cmds, see parseCMDs.
func ParseWidget(comm string) (*WidgetStructure, error) {
	cmds, err := parseCMDs(comm)
	if err != nil {
		return nil, err
	}

	if cmds["@support"] == "" {
		return nil, errors.New("parsing error: @support should be provided in comment")
	}
	if cmds["@design"] == "" {
		return nil, errors.New("parsing error: @design should be provided in comment")
	}
	if cmds["@widgetify"] == "" {
		return nil, errors.New("parsing error: @widgetify should be provided in comment")
	}

	// support
	platformSupport, err := parsePlatformSupportCMD(cmds["@support"])
	if err != nil {
		return nil, err
	}

	// design url
	u, err := url.Parse(cmds["@design"])
	if err != nil {
		return nil, errors.New("parsing error: @design url format is not correct")
	}

	// widgetify url
	widgetifyUrl, err := url.Parse(cmds["@widgetify"])
	if err != nil {
		return nil, errors.New("parsing error: @widgetify url format is not correct")
	}

	return &WidgetStructure{
		PlatformSupport: platformSupport,
		DesignURL:       u,
		WidgetifyUrl:    widgetifyUrl,
	}, nil
}

// ParseWidgetField parses comment of a specific field of a specific widget from the protofile.
// It should have @support cmd.
// For more information on cmds, see parseCMDs.
func ParseWidgetField(comm string) (*WidgetFieldStructure, error) {
	cmds, err := parseCMDs(comm)
	if err != nil {
		return nil, err
	}

	if cmds["@support"] == "" {
		return nil, errors.New("parsing error: @support should be provided in comment")
	}

	// support
	platformSupport, err := parsePlatformSupportCMD(cmds["@support"])
	if err != nil {
		return nil, err
	}

	return &WidgetFieldStructure{
		PlatformSupport: platformSupport,
	}, nil
}

// parseCMDs parse commands in the comment and return map of these commands with their values.
// commands start with "@" followed by ":" . They should appear at the beginning of the line.
// The following string after ":" till the "\n" is called value of the command.
// comments should be in "slash, asterisk" format.
// Example:
// /*
//
//	@cmd0: value0
//	@cmd1: value1
//
// */
func parseCMDs(comm string) (map[string]string, error) {
	// comment structure should be in "slash, asterisk" format
	comm = strings.TrimSpace(comm)
	lines := strings.Split(comm, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimSpace(l)
	}

	cmds := make(map[string]string)

	for _, l := range lines {
		ll := strings.SplitN(l, ":", 2)
		if len(ll) < 2 {
			continue
		}

		cmd := ll[0]
		cmd = strings.TrimSpace(cmd)

		if cmd[0] == '@' {
			// check multiple same cmd
			if cmds[cmd] != "" {
				return nil, errors.New("parsing error: multiple and same @[cmd] provided which should be at most one")
			}

			val := ll[1]
			val = strings.TrimSpace(val)
			cmds[cmd] = val
		}
	}

	return cmds, nil
}

// parsePlatformSupportCMD parse @support cmd value
// Example:
// /*
//
//	@support: android(8080-), ios(8081-8082), web
//
// */
func parsePlatformSupportCMD(commLine string) (*PlatformSupportStructure, error) {
	supportP := &PlatformSupportStructure{}
	supportS := strings.Split(commLine, ",")

	for _, s := range supportS {
		s = strings.TrimSpace(s)
		ss := strings.FieldsFunc(s, func(r rune) bool {
			return r == '(' || r == ')'
		})

		platform := ss[0]

		switch platform {
		case "android":
			if len(ss) != 2 {
				return nil, errors.New("parsing error: @support android didn't follow the specified format")
			}

			// check multiple @support android
			if supportP.AndroidMinVersion != 0 && supportP.AndroidMaxVersion != 0 {
				return nil, errors.New("parsing error: multiple @support android provdied, should be at most one")
			}

			minVersion, maxVersion, err := parseSupportClientVersions(ss[1])
			if err != nil {
				return nil, err
			}

			supportP.AndroidMinVersion = minVersion
			supportP.AndroidMaxVersion = maxVersion

		case "ios":
			if len(ss) != 2 {
				return nil, errors.New("parsing error: @support ios didn't follow the specified format")
			}

			// check multiple @support ios
			if supportP.IOSMinVersion != 0 && supportP.IOSMaxVersion != 0 {
				return nil, errors.New("parsing error: multiple @support ios provdied, should be at most one")
			}

			minVersion, maxVersion, err := parseSupportClientVersions(ss[1])
			if err != nil {
				return nil, err
			}

			supportP.IOSMinVersion = minVersion
			supportP.IOSMaxVersion = maxVersion

		case "web":
			if len(ss) != 1 {
				return nil, errors.New("parsing error: @support web didn't follow the specified format")
			}

			// check multiple @support web
			if supportP.IsWebSupported {
				return nil, errors.New("parsing error: multiple @support web provdied, should be at most one")
			}

			supportP.IsWebSupported = true

		default:
			return nil, errors.New("parsing error: @support should have a 'android','ios' or 'web' platform")
		}
	}

	return supportP, nil
}

func parseSupportClientVersions(versionS string) (minVersion int, maxVersion int, err error) {
	clientErr := errors.New("parsing error: @support android or ios didn't follow the specified format")

	versions := strings.Split(versionS, "-")
	if len(versions) != 2 {
		return 0, 0, clientErr
	}

	minVersionS, maxVersionS := versions[0], versions[1]

	minVersion, err = strconv.Atoi(minVersionS)
	if err != nil {
		return 0, 0, clientErr
	}
	if maxVersionS != "" {
		maxVersion, err = strconv.Atoi(maxVersionS)
		if err != nil {
			return 0, 0, clientErr
		}
	}

	return minVersion, maxVersion, nil
}

package pac

import (
	"chreddal/pkgs/logger"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var modlog = logger.StdLog.Specific("PAC agent")

// FromFile parse PAC script from a file
func FromFile(
	name string,
	adapt *HostAdapter,
	exts ...HostExtension,
) (Agent, error) {
	s, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	if !s.Mode().IsRegular() {
		return nil, fmt.Errorf("%s is not a regular file", name)
	}
	fp, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return FromStream(name, fp, adapt, exts...)
}

// FromString parse PAC string
func FromString(
	name string,
	content string,
	adapt *HostAdapter,
	exts ...HostExtension,
) (Agent, error) {
	return initAgent(name, content, adapt, exts...)
}

// FromStream parse PAC script from a stream IO
func FromStream(
	name string,
	fp io.Reader,
	adapt *HostAdapter,
	exts ...HostExtension,
) (Agent, error) {
	bin, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, fmt.Errorf("Failed to read stream - %s", err.Error())
	}
	return initAgent(name, string(bin), adapt, exts...)
}

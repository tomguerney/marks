package opener

import (
	"errors"
	"reflect"
	"testing"

	"github.com/abitofoldtom/marks/mocks"
)

func newTestOpener() *opener {
	return &opener{
		config:    mocks.NewConfig(),
		commander: newMockCommander(),
	}
}

type mockCommmander struct {
	commandFn         func(name string, arg ...string) combinedOutputter
	commandFnCalled   bool
	combinedOutputter combinedOutputter
}

func newMockCommander() *mockCommmander {
	return &mockCommmander{
		commandFn:         defaultCommandFn,
		combinedOutputter: newMockCombinedOutputter(),
	}
}

func (c *mockCommmander) Command(name string, arg ...string) combinedOutputter {
	c.commandFnCalled = true
	if combinedOutputter := c.commandFn(name, arg...); combinedOutputter != nil {
		return combinedOutputter
	}
	return c.combinedOutputter
}

var defaultCommandFn = func(name string, arg ...string) combinedOutputter {
	return nil
}

type mockCombinedOutputter struct {
	combinedOutputFn       func() ([]byte, error)
	combinedOutputFnCalled bool
}

func newMockCombinedOutputter() *mockCombinedOutputter {
	return &mockCombinedOutputter{
		combinedOutputFn: defaultCombinedOutputFn,
	}
}

func (c *mockCombinedOutputter) CombinedOutput() ([]byte, error) {
	c.combinedOutputFnCalled = true
	return c.combinedOutputFn()
}

var defaultCombinedOutputFn = func() ([]byte, error) {
	return []byte("combined output"), nil
}

func TestOpenSuccess(t *testing.T) {
	o := newTestOpener()
	url := "https://www.url.com"
	browser := "chrome"
	o.config.ChromeOpenArgs = "-a \"Google Chrome\" {{.Url}}"
	commandFn := func(actualName string, actualArgs ...string) combinedOutputter {
		expectedName := "open"
		expectedArgs := []string{"-a", "Google Chrome", url}
		if actualName != expectedName {
			t.Fatalf("expected %v, received %v", expectedName, actualName)
		}
		if !reflect.DeepEqual(actualArgs, expectedArgs) {
			t.Fatalf("expected %v, received %v", expectedArgs, actualArgs)
		}
		return nil
	}
	o.commander.(*mockCommmander).commandFn = commandFn
	err := o.Open(url, browser)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !o.commander.(*mockCommmander).commandFnCalled {
		t.Fatal("Command should have been called")
	}
}

func TestOpenTemplateFail(t *testing.T) {
	o := newTestOpener()
	url := "https://www.url.com"
	browser := "not a browser"
	err := o.Open(url, browser)
	t.Log("Expected error: ", err.Error())
	if err == nil {
		t.Fatal("should return error")
	}
	if o.commander.(*mockCommmander).commandFnCalled {
		t.Fatal("Command should not have been called")
	}
}

func TestOpenInterplateFail(t *testing.T) {
	o := newTestOpener()
	o.config.ChromeOpenArgs = "{{.Notafield}}"
	url := "https://www.url.com"
	browser := "chrome"
	err := o.Open(url, browser)
	t.Log("Expected error: ", err.Error())
	if err == nil {
		t.Fatal("should return error")
	}
	if o.commander.(*mockCommmander).commandFnCalled {
		t.Fatal("Command should not have been called")
	}
}

func TestOpenCombinedOutputFail(t *testing.T) {
	o := newTestOpener()
	url := "https://www.url.com"
	browser := "chrome"
	combinedOutputFn := func() ([]byte, error) {
		return nil, errors.New("combined output error")
	}
	o.commander.(*mockCommmander).
		combinedOutputter.(*mockCombinedOutputter).combinedOutputFn = combinedOutputFn
	err := o.Open(url, browser)
	t.Log("Expected error: ", err.Error())
	if err == nil {
		t.Fatal("should return error")
	}
	if !o.commander.(*mockCommmander).
		combinedOutputter.(*mockCombinedOutputter).combinedOutputFnCalled {
		t.Fatal("combinedOutput should have been called")
	}
}

func TestTemplateSuccess(t *testing.T) {
	o := newTestOpener()
	expected := "firefox args"
	o.config.FirefoxOpenArgs = expected
	actual, err := o.template("firefox")
	if err != nil {
		t.Fatal(err.Error())
	}
	if actual != expected {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}

func TestTemplateFail(t *testing.T) {
	o := newTestOpener()
	expected := ""
	actual, err := o.template("not a browser")
	if err == nil {
		t.Fatal("expected error")
	}
	t.Log("Expected error: ", err.Error())
	if actual != "" {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}

func TestInterpolateTemplateSuccess(t *testing.T) {
	o := newTestOpener()
	url := "testUrl"
	expected := "***testUrl***"
	template := "***{{.Url}}***"
	actual, err := o.interpolateTemplate(template, url)
	if err != nil {
		t.Fatal(err.Error())
	}
	if actual != expected {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}

func TestInterpolateTemplateFail(t *testing.T) {
	o := newTestOpener()
	url := "testUrl"
	template := "***{{.Other}}***"
	actual, err := o.interpolateTemplate(template, url)
	if err == nil {
		t.Fatal("expected error")
	}
	t.Log("Expected error: ", err.Error())
	if actual != "" {
		t.Fatal("string should be empty")
	}
	t.Log("Expected error: ", err.Error())
}

func TestSliceArgsSuccess(t *testing.T) {
	o := newTestOpener()
	args := "./foo --bar=baz \"blah blah\" -p p -b"
	expected := []string{"./foo", "--bar=baz", "blah blah", "-p", "p", "-b"}
	actual, err := o.sliceArgs(args)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, received %v", expected, actual)
	}
}

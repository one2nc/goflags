package goflags

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestGenerateDefaultConfig(t *testing.T) {
	flagSet := NewFlagSet()

	example := `# generated by https://github.com/projectdiscovery/goflags

# default value for a test flag example
#test: test-default-value

# string slice flag example value
#slice: ["item1", "item2"]`

	var data string
	var data2 StringSlice
	flagSet.StringVar(&data, "test", "test-default-value", "Default value for a test flag example")
	flagSet.StringSliceVar(&data2, "slice", []string{"item1", "item2"}, "String slice flag example value")
	defaultConfig := string(flagSet.generateDefaultConfig())
	parts := strings.SplitN(defaultConfig, "\n", 2)

	require.Equal(t, example, parts[1], "could not get correct default config")
	tearDown(t.Name())
}

func TestConfigFileDataTypes(t *testing.T) {
	flagSet := NewFlagSet()
	var data string
	var data2 StringSlice
	var data3 int
	var data4 bool

	flagSet.StringVar(&data, "string-value", "", "Default value for a test flag example")
	flagSet.StringSliceVar(&data2, "slice-value", []string{}, "String slice flag example value")
	flagSet.IntVar(&data3, "int-value", 0, "Int value example")
	flagSet.BoolVar(&data4, "bool-value", false, "Bool value example")

	configFileData := `
string-value: test
slice-value:
 - test
 - test2
severity:
 - info
 - high
int-value: 543
bool-value: true`
	err := ioutil.WriteFile("test.yaml", []byte(configFileData), os.ModePerm)
	require.Nil(t, err, "could not write temporary config")
	//defer os.Remove("test.yaml")

	err = flagSet.MergeConfigFile("test.yaml")
	require.Nil(t, err, "could not merge temporary config")

	require.Equal(t, "test", data, "could not get correct string")
	require.Equal(t, StringSlice{"test", "test2"}, data2, "could not get correct string slice")
	require.Equal(t, 543, data3, "could not get correct int")
	require.Equal(t, true, data4, "could not get correct bool")

	tearDown(t.Name())
}

func TestUsageOrder(t *testing.T) {
	flagSet := NewFlagSet()

	var stringData string
	var stringSliceData StringSlice
	var stringSliceData2 StringSlice
	var intData int
	var boolData bool

	flagSet.SetGroup("String", "String")
	flagSet.StringVar(&stringData, "string-value", "", "String example value example").Group("String")
	flagSet.StringVarP(&stringData, "", "ts2", "test-string", "String with default value example #2").Group("String")
	flagSet.StringVar(&stringData, "string-with-default-value", "test-string", "String with default value example").Group("String")
	flagSet.StringVarP(&stringData, "string-with-default-value2", "ts", "test-string", "String with default value example #2").Group("String")

	flagSet.SetGroup("StringSlice", "StringSlice")
	flagSet.StringSliceVar(&stringSliceData, "slice-value", []string{}, "String slice flag example value").Group("StringSlice")
	flagSet.StringSliceVarP(&stringSliceData, "slice-value2", "sv", []string{}, "String slice flag example value #2").Group("StringSlice")
	flagSet.StringSliceVar(&stringSliceData, "slice-with-default-value", []string{"a", "b", "c"}, "String slice flag with default example values").Group("StringSlice")
	flagSet.StringSliceVarP(&stringSliceData2, "slice-with-default-value2", "swdf", []string{"a", "b", "c"}, "String slice flag with default example values #2").Group("StringSlice")

	flagSet.SetGroup("Integer", "Integer")
	flagSet.IntVar(&intData, "int-value", 0, "Int value example").Group("Integer")
	flagSet.IntVarP(&intData, "int-value2", "iv", 0, "Int value example #2").Group("Integer")
	flagSet.IntVar(&intData, "int-with-default-value", 12, "Int with default value example").Group("Integer")
	flagSet.IntVarP(&intData, "int-with-default-value2", "iwdv", 12, "Int with default value example #2").Group("Integer")

	flagSet.SetGroup("Bool", "Boolean")
	flagSet.BoolVar(&boolData, "bool-value", false, "Bool value example").Group("Bool")
	flagSet.BoolVarP(&boolData, "bool-value2", "bv", false, "Bool value example #2").Group("Bool")
	flagSet.BoolVar(&boolData, "bool-with-default-value", true, "Bool with default value example").Group("Bool")
	flagSet.BoolVarP(&boolData, "bool-with-default-value2", "bwdv", true, "Bool with default value example #2").Group("Bool")

	output := &bytes.Buffer{}
	flag.CommandLine.SetOutput(output)

	flagSet.usageFunc()

	resultOutput := output.String()
	actual := resultOutput[strings.Index(resultOutput, "Flags:\n"):]

	expected :=
		`Flags:
STRING:
   -string-value string                     String example value example
   -ts2 string                              String with default value example #2 (default "test-string")
   -string-with-default-value string        String with default value example (default "test-string")
   -ts, -string-with-default-value2 string  String with default value example #2 (default "test-string")
STRINGSLICE:
   -slice-value string[]                       String slice flag example value
   -sv, -slice-value2 string[]                 String slice flag example value #2
   -slice-with-default-value string[]          String slice flag with default example values (default ["a", "b", "c"])
   -swdf, -slice-with-default-value2 string[]  String slice flag with default example values #2 (default ["a", "b", "c"])
INTEGER:
   -int-value int                       Int value example
   -iv, -int-value2 int                 Int value example #2
   -int-with-default-value int          Int with default value example (default 12)
   -iwdv, -int-with-default-value2 int  Int with default value example #2 (default 12)
BOOLEAN:
   -bool-value                       Bool value example
   -bv, -bool-value2                 Bool value example #2
   -bool-with-default-value          Bool with default value example (default true)
   -bwdv, -bool-with-default-value2  Bool with default value example #2 (default true)
`
	assert.Equal(t, actual, expected)

	tearDown(t.Name())
}

func TestIncorrectStringFlagsCausePanic(t *testing.T) {
	flagSet := NewFlagSet()
	var stringData string

	flagSet.StringVar(&stringData, "", "test-string", "String with default value example")
	assert.Panics(t, flagSet.usageFunc)

	// env GOOS=linux GOARCH=amd64 go build main.go -o nuclei

	tearDown(t.Name())
}

func TestIncorrectFlagsCausePanic(t *testing.T) {
	type flagPair struct {
		Short, Long string
	}

	createTestParameters := func() []flagPair {
		var result []flagPair
		result = append(result, flagPair{"", ""})

		badFlagNames := [...]string{" ", "\t", "\n"}
		for _, badFlagName := range badFlagNames {
			result = append(result, flagPair{"", badFlagName})
			result = append(result, flagPair{badFlagName, ""})
			result = append(result, flagPair{badFlagName, badFlagName})
		}
		return result
	}

	for index, tuple := range createTestParameters() {
		uniqueName := strconv.Itoa(index)
		t.Run(uniqueName, func(t *testing.T) {
			assert.Panics(t, func() {
				tearDown(uniqueName)

				flagSet := NewFlagSet()
				var stringData string

				flagSet.StringVarP(&stringData, tuple.Short, tuple.Long, "test-string", "String with default value example")
				flagSet.usageFunc()
			})
		})
	}
}

type testSliceValue []interface{}

func (value testSliceValue) String() string   { return "" }
func (value testSliceValue) Set(string) error { return nil }

func TestCustomSliceUsageType(t *testing.T) {
	testCases := map[string]flag.Flag{
		"string[]": {Value: &StringSlice{}},
		"value[]":  {Value: &testSliceValue{}},
	}

	for expected, currentFlag := range testCases {
		result := createUsageTypeAndDescription(&currentFlag, reflect.TypeOf(currentFlag.Value))
		assert.Equal(t, expected, strings.TrimSpace(result))
	}
}

func TestParseStringSlice(t *testing.T) {
	flagSet := NewFlagSet()

	var stringSlice StringSlice
	flagSet.StringSliceVarP(&stringSlice, "header", "H", []string{}, "Header values. Expected usage: -H \"header1\":\"value1\" -H \"header2\":\"value2\"")

	header1 := "\"header1:value1\""
	header2 := "\" HEADER 2: VALUE2  \""
	header3 := "\"header3\":\"value3, value4\""

	os.Args = []string{
		"./appName",
		"-H", header1,
		"-header", header2,
		"-H", header3,
	}

	err := flagSet.Parse()
	assert.Nil(t, err)

	assert.Equal(t, stringSlice, StringSlice{header1, header2, header3})
}

func tearDown(uniqueValue string) { // sadly there is no official support for setup/teardown/test
	flag.CommandLine = flag.NewFlagSet(uniqueValue, flag.ContinueOnError)
	flag.CommandLine.Usage = flag.Usage
}

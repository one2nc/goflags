package goflags

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizedStringSlicePositive(t *testing.T) {
	expectedABC := []string{"aa", "bb", "cc"}
	expectedFilePath := []string{"/root/home/file0"}

	slices := map[string][]string{
		"aa,bb,cc":                 expectedABC,
		"  aa, bb,  cc   ":         expectedABC,
		"  `aa`, 'bb',  \"cc\"   ": expectedABC,
		"  `aa`, bb,  \"cc\"   ":   expectedABC,
		"  `aa, bb,  cc\"   ":      expectedABC,
		"  \"aa\", bb,  cc\"   ":   expectedABC,
		"\n  aa, \tbb,  cc\r   ":   expectedABC,

		"\"value1\",value,'value3'": {"value1", "value", "value3"},
		"\"value1\",VALUE,'value3'": {"value1", "value", "value3"},

		"\"/root/home/file0\"":       expectedFilePath,
		"'/root/home/file0'":         expectedFilePath,
		"`/root/home/file0`":         expectedFilePath,
		"\"/root/home/file0\",":      expectedFilePath,
		",\"/root/home/file0\",":     expectedFilePath,
		",\"/root/home/file0\"":      expectedFilePath,
		",,\"/root/home/file0\"":     expectedFilePath,
		"\"\",,\"/root/home/file0\"": expectedFilePath,
		"\" \",\"/root/home/file0\"": expectedFilePath,
		"\"/root/home/file0\",\"\"":  expectedFilePath,
		"/root/home/file0":           expectedFilePath,

		"\"/root/home/file2\",\"/root/home/file3\"":             {"/root/home/file2", "/root/home/file3"},
		"/root/home/file4,/root/home/file5":                     {"/root/home/file4", "/root/home/file5"},
		"\"/root/home/file4,/root/home/file5\"":                 {"/root/home/file4,/root/home/file5"},
		"\"/root/home/file6\",/root/home/file7":                 {"/root/home/file6", "/root/home/file7"},
		"\"c:\\my files\\bug,bounty\"":                          {"c:\\my files\\bug,bounty"},
		"\"c:\\my files\\bug,bounty\",c:\\my_files\\bug bounty": {"c:\\my files\\bug,bounty", "c:\\my_files\\bug bounty"},
	}

	for value, expected := range slices {
		result, err := ToStringSlice(value, NormalizedStringSliceOptions)
		fmt.Println(result)
		assert.Nil(t, err)
		assert.Equal(t, result, expected)
	}
}

func TestNormalizedStringSliceNegative(t *testing.T) {
	slices := []string{
		"\"/root/home/file0",
		"'/root/home/file0",
		"`/root/home/file0",
		"\"/root/home/file0'",
		"\"/root/home/file0`",
	}

	for _, value := range slices {
		result, err := ToStringSlice(value, NormalizedStringSliceOptions)
		assert.Nil(t, result)
		assert.NotNil(t, err)
	}
}

func TestNormalizedOriginalStringSlice(t *testing.T) {
	result, err := ToStringSlice("/Users/Home/Test/test.yaml", NormalizedOriginalStringSliceOptions)
	assert.Nil(t, err)
	assert.Equal(t, []string{"/Users/Home/Test/test.yaml"}, result, "could not get correct path")

	result, err = ToStringSlice("'test user'", NormalizedOriginalStringSliceOptions)
	assert.Nil(t, err)
	assert.Equal(t, []string{"test user"}, result, "could not get correct path")
}

func TestFileNormalizedStringSliceOptions(t *testing.T) {
	result, err := ToStringSlice("/Users/Home/Test/test.yaml", FileNormalizedStringSliceOptions)
	assert.Nil(t, err)
	assert.Equal(t, []string{"/users/home/test/test.yaml"}, result, "could not get correct path")

	result, err = ToStringSlice("'Test User'", FileNormalizedStringSliceOptions)
	assert.Nil(t, err)
	assert.Equal(t, []string{"test user"}, result, "could not get correct path")
}

func TestFileStringSliceOptions(t *testing.T) {
	filename := "test.txt"
	_ = os.WriteFile(filename, []byte("value1,value2\nvalue3"), 0644)
	defer os.RemoveAll(filename)

	result, err := ToStringSlice(filename, FileStringSliceOptions)
	assert.Nil(t, err)
	assert.Equal(t, []string{"value1,value2", "value3"}, result, "could not get correct path")
}

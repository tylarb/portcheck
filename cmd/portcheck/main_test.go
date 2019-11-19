/*
Released under MIT license, copyright 2019 Tyler Ramer
*/

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePorts(t *testing.T) {
	validString := "1,2-4,10-10"
	expectedResult := []int{1, 2, 3, 4, 10}

	actual, err := parsePorts(validString)
	if err != nil {
		t.Logf("expected successful parse but instead got err %q", err)
		t.Fail()
	}
	assert.Equal(t, expectedResult, actual)

	invalidString1 := "1,2--10,"
	_, err = parsePorts(invalidString1)
	assert.Error(t, err)

	invalidString2 := "10-2"
	_, err = parsePorts(invalidString2)
	assert.Error(t, err)

}

func TestComposeByteMessage(t *testing.T) {
	headerString := "header"
	footerString := "footer"

	header := []byte(headerString)
	footer := []byte(footerString)

	composedSize := 500

	composedBytes := composeByteMessage(header, footer, composedSize)

	// length of composed bytes is size passed into function
	assert.Equal(t, 500, len(composedBytes))

	// first len(header) bytes of composed == header
	assert.Equal(t, header, composedBytes[:len(header)])

	// last len(footer) bytes of composed == footer
	assert.Equal(t, footer, composedBytes[composedSize-len(footer):])

}

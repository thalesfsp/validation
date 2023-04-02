package validation_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thalesfsp/validation"
)

type TestStruct struct {
	Name  string    `validate:"stringContains=example"`
	Date1 time.Time `validate:"dateBefore=now"`
	Date2 time.Time `validate:"dateAfter=now"`
	Date3 time.Time `validate:"dateBefore=2030-01-01"`
	Date4 time.Time `validate:"dateAfter=2000-01-01"`
}

func TestValidation(t *testing.T) {
	testCases := []struct {
		name          string
		testStruct    TestStruct
		expectedError bool
	}{
		{
			name: "Valid struct",
			testStruct: TestStruct{
				Name:  "This is an example string",
				Date1: time.Now().Add(-time.Minute),
				Date2: time.Now().Add(time.Minute),
				Date3: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				Date4: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedError: false,
		},
		{
			name: "Invalid name",
			testStruct: TestStruct{
				Name:  "This is not a valid string",
				Date1: time.Now().Add(-time.Minute),
				Date2: time.Now().Add(time.Minute),
				Date3: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				Date4: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedError: true,
		},
		{
			name: "Invalid Date1",
			testStruct: TestStruct{
				Name:  "This is an example string",
				Date1: time.Now().Add(time.Minute),
				Date2: time.Now().Add(time.Minute),
				Date3: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				Date4: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedError: true,
		},
		{
			name: "Invalid Date2",
			testStruct: TestStruct{
				Name:  "This is an example string",
				Date1: time.Now().Add(-time.Minute),
				Date2: time.Now().Add(-time.Minute),
				Date3: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				Date4: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedError: true,
		},
		{
			name: "Invalid Date3",
			testStruct: TestStruct{
				Name:  "This is an example string",
				Date1: time.Now().Add(-time.Minute),
				Date2: time.Now().Add(time.Minute),
				Date3: time.Date(2040, 1, 1, 0, 0, 0, 0, time.UTC),
				Date4: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedError: true,
		},
		{
			name: "Invalid Date4",
			testStruct: TestStruct{
				Name:  "This is an example string",
				Date1: time.Now().Add(-time.Minute),
				Date2: time.Now().Add(time.Minute),
				Date3: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				Date4: time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedError: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validation.Validate(&tc.testStruct)
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

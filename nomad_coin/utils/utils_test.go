package utils

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestHashObject(t *testing.T) {
	t.Run("Hash is always same", func(t *testing.T) {
		// Given
		s := struct{ Test string }{Test: "test"}
		expected_hash_value := "654209534fb59c0f07ae82d4a57289e6b247e0e103778a251e885030a1cfd038"

		// When
		result := HashObject(s)

		// Then
		if result != expected_hash_value {
			t.Errorf("Expected %s, got %s", expected_hash_value, result)
		}
	})
	t.Run("Hash is hex encoded", func(t *testing.T) {
		// Given
		s := struct{ Test string }{Test: "random_string"}

		// When
		result := HashObject(s)

		// Then
		_, err := hex.DecodeString(result)
		if err != nil {
			t.Error("Hash should be hex encoded")
		}
	})
}

func TestErrHandler(t *testing.T) {
	oldLogFn := LogFn
	defer func() {
		LogFn = oldLogFn
	}()

	// Given
	logFnCalled := false
	LogFn = func(v ...interface{}) {
		logFnCalled = true
	}
	err := errors.New("test")

	// When
	ErrHandler(err)

	// Then
	if logFnCalled != true {
		t.Error("ErrHandler should call log.Panic")
	}
}

func TestGetNowUnixTimestamp(t *testing.T) {
	oldFn := timeNowFn
	defer func() {
		timeNowFn = oldFn
	}()

	// Given
	timeNowFnCalled := false
	timeNowFn = func() time.Time {
		timeNowFnCalled = true
		return time.Now()
	}

	// When
	GetNowUnixTimestamp()

	// Then
	if timeNowFnCalled != true {
		t.Error("GetNowUnixTimestamp should call time.now")
	}
}

func TestSplitter(t *testing.T) {
	// Given
	type test struct {
		input          string
		sep            string
		index          int
		expectedOutput string
		actualOutput   string
	}
	testCaseList := []*test{
		{input: "0:6:0", sep: ":", index: 1, expectedOutput: "6"},
		{input: "0:6:0", sep: ":", index: 10, expectedOutput: ""},
		{input: "0:6:0", sep: "/", index: 0, expectedOutput: "0:6:0"},
	}

	// When
	for _, testCase := range testCaseList {
		testCase.actualOutput = Splitter(testCase.input, testCase.sep, testCase.index)
	}

	// Then
	for _, testCase := range testCaseList {
		if testCase.expectedOutput != testCase.actualOutput {
			t.Errorf("Expected %s but got %s", testCase.expectedOutput, testCase.actualOutput)
		}
	}
}

func TestToJson(t *testing.T) {
	t.Run("ToJson should return slice", func(t *testing.T) {
		// Given
		type testStruct struct {
			Value string
		}
		testCase := testStruct{"test"}

		// When
		result := ToJson(testCase)

		// Then
		typeOfResult := reflect.TypeOf(result).Kind()
		if typeOfResult != reflect.Slice {
			t.Errorf("Expected %v but got %v", reflect.Slice, typeOfResult)
		}
	})

	t.Run("ToJson should be restored by json.Unmarshal", func(t *testing.T) {
		// Given
		type testStruct struct {
			Value string
		}
		testCase := testStruct{"test"}

		// When
		result := ToJson(testCase)

		// Then
		var restored testStruct
		json.Unmarshal(result, &restored)
		if !reflect.DeepEqual(testCase, restored) {
			t.Errorf("Expected %v but got %v", testCase, restored)
		}
	})
}

func TestObjectToBytes(t *testing.T) {
	t.Run("ObjectToBytes should return error", func(t *testing.T) {
		// When
		_, err := ObjectToBytes(nil)

		// Then
		if err == nil {
			t.Error("ObjectToBytes didn't return error")
		}
	})

	t.Run("ObjectToBytes should return bytes", func(t *testing.T) {
		// Given
		type testStruct struct {
			Value string
		}
		testCase := testStruct{"test"}

		// When
		result, err := ObjectToBytes(testCase)

		// Then
		if err != nil {
			t.Error("ObjectToBytes returned error")
		}

		typeOfResult := reflect.TypeOf(result).Kind()
		if typeOfResult != reflect.Slice {
			t.Errorf("Expected %v but got %v", reflect.Slice, typeOfResult)
		}
	})

}

func TestObjectFromBytes(t *testing.T) {
	t.Run("ObjectFromBytes should restore ObjectToBytes", func(t *testing.T) {
		// Given
		type testStruct struct {
			Value string
		}
		testCase := testStruct{"test"}

		// When
		result, err := ObjectToBytes(testCase)
		if err != nil {
			t.Error("ObjectToBytes returned error")
		}

		var restored testStruct
		err = ObjectFromBytes(&restored, result)
		if err != nil {
			t.Error("ObjectFromBytes returned error")
		}

		// Then
		if !reflect.DeepEqual(testCase, restored) {
			t.Errorf("Expected %v but got %v", testCase, restored)
		}
	})
}

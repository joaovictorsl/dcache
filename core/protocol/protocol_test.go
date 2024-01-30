package protocol

import (
	"bytes"
	"testing"
	"time"

	"github.com/joaovictorsl/dcache/core"
	"github.com/joaovictorsl/dcache/core/command"
)

func TestParseCommandSet(t *testing.T) {
	t.Run("should return a set command with ttl 5000ms", func(t *testing.T) {
		// "SET Foo Bar 5000"
		foo := "Foo"
		bar := []byte("Bar")
		cmd := CreateSetCmd(foo, bar, 5000)

		expected := &command.SetCommand{
			Key:   foo,
			Value: bar,
			TTL:   5000 * time.Millisecond,
		}

		actual, err := ParseCommand(cmd)
		if err != nil {
			t.Errorf("parseCommand(%q) returned error %q", cmd, err)
		}

		actualSet := actual.(*command.SetCommand)
		if actualSet.Key != expected.Key || !bytes.Equal(actualSet.Value, expected.Value) || actualSet.TTL != expected.TTL {
			t.Errorf("parseCommand(%q) = %v, want %v", cmd, actual, expected)
		}
	})

	t.Run("should return an error if command is invalid", func(t *testing.T) {
		cmdInvalidKeySizeOver := CreateSetCmd("Foo", []byte("Bar"), 5000)
		cmdInvalidKeySizeOver[1] = 4
		_, err := ParseCommand(cmdInvalidKeySizeOver)
		expected := core.INVALID_SET_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdInvalidKeySizeOver)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdInvalidKeySizeOver, err, expected)
		}

		cmdInvalidKeySizeUnder := CreateSetCmd("Foo", []byte("Bar"), 5000)
		cmdInvalidKeySizeUnder[1] = 2
		_, err = ParseCommand(cmdInvalidKeySizeUnder)
		expected = core.INVALID_SET_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdInvalidKeySizeUnder)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdInvalidKeySizeUnder, err, expected)
		}

		cmdInvalidValueSizeOver := CreateSetCmd("Foo", []byte("Bar"), 5000)
		cmdInvalidValueSizeOver[5] = 4
		_, err = ParseCommand(cmdInvalidValueSizeOver)
		expected = core.INVALID_SET_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdInvalidValueSizeOver)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdInvalidValueSizeOver, err, expected)
		}

		cmdInvalidValueSizeUnder := CreateSetCmd("Foo", []byte("Bar"), 5000)
		cmdInvalidValueSizeUnder[5] = 2
		_, err = ParseCommand(cmdInvalidValueSizeUnder)
		expected = core.INVALID_SET_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdInvalidValueSizeUnder)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdInvalidValueSizeUnder, err, expected)
		}

		cmdEmpty := make([]byte, 0)
		_, err = ParseCommand(cmdEmpty)
		expected = core.INVALID_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdEmpty)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdEmpty, err, expected)
		}

		cmdOnlyKey := CreateSetCmd("Foo", []byte("Bar"), 5000)[0:5]
		_, err = ParseCommand(cmdOnlyKey)
		expected = core.INVALID_SET_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdOnlyKey)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdOnlyKey, err, expected)
		}

		cmdOnlyMissingTtl := CreateSetCmd("Foo", []byte("Bar"), 5000)[0:9]
		_, err = ParseCommand(cmdOnlyMissingTtl)
		expected = core.INVALID_SET_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdOnlyMissingTtl)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdOnlyMissingTtl, err, expected)
		}

		cmdInvalidOneByte := []byte{core.CMD_SET}
		_, err = ParseCommand(cmdInvalidOneByte)
		expected = core.INVALID_SET_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdInvalidOneByte)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdInvalidOneByte, err, expected)
		}

		cmdKeyLen0 := CreateSetCmd("Foo", []byte("Bar"), 5000)
		cmdKeyLen0[1] = 0
		_, err = ParseCommand(cmdKeyLen0)
		expected = core.INVALID_SET_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdKeyLen0)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdKeyLen0, err, expected)
		}
	})
}

func TestParseCommandDelete(t *testing.T) {
	t.Run("should return a delete command for key Foo", func(t *testing.T) {
		// "DELETE Foo"
		foo := "Foo"
		cmd := CreateDeleteCmd(foo)

		expected := &command.DeleteCommand{
			Key: foo,
		}

		actual, err := ParseCommand(cmd)
		if err != nil {
			t.Errorf("parseCommand(%q) returned error %q", cmd, err)
		}

		actualSet := actual.(*command.DeleteCommand)
		if actualSet.Key != expected.Key {
			t.Errorf("parseCommand(%q) = %v, want %v", cmd, actual, expected)
		}
	})

	t.Run("should return an error if command is invalid", func(t *testing.T) {
		cmdInvalidOneByte := []byte{core.CMD_DELETE}
		_, err := ParseCommand(cmdInvalidOneByte)
		expected := core.INVALID_DELETE_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdInvalidOneByte)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdInvalidOneByte, err, expected)
		}

		cmdInvalidNoKey := []byte{core.CMD_DELETE, 2}
		_, err = ParseCommand(cmdInvalidNoKey)
		expected = core.INVALID_DELETE_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdInvalidNoKey)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdInvalidNoKey, err, expected)
		}

		cmdInvalidKeySizeOver := CreateDeleteCmd("Foo")
		cmdInvalidKeySizeOver[1] = 4
		_, err = ParseCommand(cmdInvalidKeySizeOver)
		expected = core.INVALID_DELETE_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdInvalidKeySizeOver)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdInvalidKeySizeOver, err, expected)
		}

		cmdInvalidKeySizeUnder := CreateDeleteCmd("Foo")
		cmdInvalidKeySizeUnder[1] = 2
		_, err = ParseCommand(cmdInvalidKeySizeUnder)
		expected = core.INVALID_DELETE_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdInvalidKeySizeUnder)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdInvalidKeySizeUnder, err, expected)
		}

		cmdEmpty := make([]byte, 0)
		_, err = ParseCommand(cmdEmpty)
		expected = core.INVALID_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdEmpty)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdEmpty, err, expected)
		}

		cmdKeyLen0 := CreateDeleteCmd("Foo")
		cmdKeyLen0[1] = 0
		_, err = ParseCommand(cmdKeyLen0)
		expected = core.INVALID_DELETE_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdKeyLen0)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdKeyLen0, err, expected)
		}
	})
}

func TestParseCommandGet(t *testing.T) {
	t.Run("should return a get command for key Foo", func(t *testing.T) {
		// "GET Foo"
		foo := "Foo"
		cmd := CreateGetCmd(foo)

		expected := &command.GetCommand{
			Key: foo,
		}

		actual, err := ParseCommand(cmd)
		if err != nil {
			t.Errorf("parseCommand(%q) returned error %q", cmd, err)
		}

		actualSet := actual.(*command.GetCommand)
		if actualSet.Key != expected.Key {
			t.Errorf("parseCommand(%q) = %v, want %v", cmd, actual, expected)
		}
	})

	t.Run("should return an error if command is invalid", func(t *testing.T) {
		cmdInvalidOneByte := []byte{core.CMD_GET}
		_, err := ParseCommand(cmdInvalidOneByte)
		expected := core.INVALID_GET_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdInvalidOneByte)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdInvalidOneByte, err, expected)
		}

		cmdInvalidNoKey := []byte{core.CMD_GET, 2}
		_, err = ParseCommand(cmdInvalidNoKey)
		expected = core.INVALID_GET_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdInvalidNoKey)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdInvalidNoKey, err, expected)
		}

		cmdInvalidKeySizeOver := CreateGetCmd("Foo")
		cmdInvalidKeySizeOver[1] = 4
		_, err = ParseCommand(cmdInvalidKeySizeOver)
		expected = core.INVALID_GET_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdInvalidKeySizeOver)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdInvalidKeySizeOver, err, expected)
		}

		cmdInvalidKeySizeUnder := CreateGetCmd("Foo")
		cmdInvalidKeySizeUnder[1] = 2
		_, err = ParseCommand(cmdInvalidKeySizeUnder)
		expected = core.INVALID_GET_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdInvalidKeySizeUnder)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdInvalidKeySizeUnder, err, expected)
		}

		cmdEmpty := make([]byte, 0)
		_, err = ParseCommand(cmdEmpty)
		expected = core.INVALID_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdEmpty)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdEmpty, err, expected)
		}

		cmdKeyLen0 := CreateGetCmd("Foo")
		cmdKeyLen0[1] = 0
		_, err = ParseCommand(cmdKeyLen0)
		expected = core.INVALID_GET_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdKeyLen0)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdKeyLen0, err, expected)
		}
	})
}

func TestParseCommandHas(t *testing.T) {
	t.Run("should return a has command for key Foo", func(t *testing.T) {
		// "HAS Foo"
		foo := "Foo"
		cmd := CreateHasCmd(foo)
		expected := &command.HasCommand{
			Key: foo,
		}

		actual, err := ParseCommand(cmd)
		if err != nil {
			t.Errorf("parseCommand(%q) returned error %q", cmd, err)
		}

		actualSet := actual.(*command.HasCommand)
		if actualSet.Key != expected.Key {
			t.Errorf("parseCommand(%q) = %v, want %v", cmd, actual, expected)
		}
	})

	t.Run("should return an error if command is invalid", func(t *testing.T) {
		cmdInvalidOneByte := []byte{core.CMD_HAS}
		_, err := ParseCommand(cmdInvalidOneByte)
		expected := core.INVALID_HAS_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdInvalidOneByte)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdInvalidOneByte, err, expected)
		}

		cmdInvalidNoKey := []byte{core.CMD_HAS, 2}
		_, err = ParseCommand(cmdInvalidNoKey)
		expected = core.INVALID_HAS_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdInvalidNoKey)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdInvalidNoKey, err, expected)
		}

		cmdInvalidKeySizeOver := CreateHasCmd("Foo")
		cmdInvalidKeySizeOver[1] = 4
		_, err = ParseCommand(cmdInvalidKeySizeOver)
		expected = core.INVALID_HAS_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdInvalidKeySizeOver)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdInvalidKeySizeOver, err, expected)
		}

		cmdInvalidKeySizeUnder := CreateHasCmd("Foo")
		cmdInvalidKeySizeUnder[1] = 2
		_, err = ParseCommand(cmdInvalidKeySizeUnder)
		expected = core.INVALID_HAS_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdInvalidKeySizeUnder)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdInvalidKeySizeUnder, err, expected)
		}

		cmdEmpty := make([]byte, 0)
		_, err = ParseCommand(cmdEmpty)
		expected = core.INVALID_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdEmpty)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdEmpty, err, expected)
		}

		cmdKeyLen0 := CreateHasCmd("Foo")
		cmdKeyLen0[1] = 0
		_, err = ParseCommand(cmdKeyLen0)
		expected = core.INVALID_HAS_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdKeyLen0)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdKeyLen0, err, expected)
		}
	})
}

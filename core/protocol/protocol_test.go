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
		cmd := command.SetCmdAsBytes(foo, bar, 5000)

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
		cmdInvalidKeySizeOver := command.SetCmdAsBytes("Foo", []byte("Bar"), 5000)
		cmdInvalidKeySizeOver[1] = 4
		_, err := ParseCommand(cmdInvalidKeySizeOver)
		expected := core.INVALID_SET_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdInvalidKeySizeOver)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdInvalidKeySizeOver, err, expected)
		}

		cmdInvalidKeySizeUnder := command.SetCmdAsBytes("Foo", []byte("Bar"), 5000)
		cmdInvalidKeySizeUnder[1] = 2
		_, err = ParseCommand(cmdInvalidKeySizeUnder)
		expected = core.INVALID_SET_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdInvalidKeySizeUnder)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdInvalidKeySizeUnder, err, expected)
		}

		cmdInvalidValueSizeOver := command.SetCmdAsBytes("Foo", []byte("Bar"), 5000)
		cmdInvalidValueSizeOver[5] = 4
		_, err = ParseCommand(cmdInvalidValueSizeOver)
		expected = core.INVALID_SET_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdInvalidValueSizeOver)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdInvalidValueSizeOver, err, expected)
		}

		cmdInvalidValueSizeUnder := command.SetCmdAsBytes("Foo", []byte("Bar"), 5000)
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

		cmdOnlyKey := command.SetCmdAsBytes("Foo", []byte("Bar"), 5000)[0:5]
		_, err = ParseCommand(cmdOnlyKey)
		expected = core.INVALID_SET_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdOnlyKey)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdOnlyKey, err, expected)
		}

		cmdOnlyMissingTtl := command.SetCmdAsBytes("Foo", []byte("Bar"), 5000)[0:9]
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

		cmdKeyLen0 := command.SetCmdAsBytes("Foo", []byte("Bar"), 5000)
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
		cmd := command.DeleteCmdAsBytes(foo)

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

		cmdInvalidKeySizeOver := command.DeleteCmdAsBytes("Foo")
		cmdInvalidKeySizeOver[1] = 4
		_, err = ParseCommand(cmdInvalidKeySizeOver)
		expected = core.INVALID_DELETE_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdInvalidKeySizeOver)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdInvalidKeySizeOver, err, expected)
		}

		cmdInvalidKeySizeUnder := command.DeleteCmdAsBytes("Foo")
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

		cmdKeyLen0 := command.DeleteCmdAsBytes("Foo")
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
		cmd := command.GetCmdAsBytes(foo)

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

		cmdInvalidKeySizeOver := command.GetCmdAsBytes("Foo")
		cmdInvalidKeySizeOver[1] = 4
		_, err = ParseCommand(cmdInvalidKeySizeOver)
		expected = core.INVALID_GET_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdInvalidKeySizeOver)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdInvalidKeySizeOver, err, expected)
		}

		cmdInvalidKeySizeUnder := command.GetCmdAsBytes("Foo")
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

		cmdKeyLen0 := command.GetCmdAsBytes("Foo")
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
		cmd := command.HasCmdAsBytes(foo)
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

		cmdInvalidKeySizeOver := command.HasCmdAsBytes("Foo")
		cmdInvalidKeySizeOver[1] = 4
		_, err = ParseCommand(cmdInvalidKeySizeOver)
		expected = core.INVALID_HAS_COMMAND
		if err == nil {
			t.Errorf("parseCommand(%q) should return error", cmdInvalidKeySizeOver)
		}

		if err.Error() != expected {
			t.Errorf("parseCommand(%q) = %q, want %q", cmdInvalidKeySizeOver, err, expected)
		}

		cmdInvalidKeySizeUnder := command.HasCmdAsBytes("Foo")
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

		cmdKeyLen0 := command.HasCmdAsBytes("Foo")
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

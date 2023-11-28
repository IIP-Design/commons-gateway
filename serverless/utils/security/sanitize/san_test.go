package sanitize

// Adapted from JS package sanitize-s3-objectkey (https://github.com/Advanon/sanitize-s3-objectkey), which is Apache 2.0 licensed

import (
	"testing"
)

func TestValidObjectKey(t *testing.T) {
	objectKey := "my.great_photos-2014-jan-myvacation.jpg"
	sanitized := DefaultKeySanitizer(objectKey)
	if objectKey != sanitized {
		t.Fatalf(`DefaultKeySanitizer modified valid key from %s to %s`, objectKey, sanitized)
	}
}

func TestRemoveSpaces(t *testing.T) {
	expected := "my.great_photos-2014janmyvacation.jpg"
	objectKey := "    my.great_photos 2014/jan/myvacation.jpg"
	sanitized := DefaultKeySanitizer(objectKey)
	if sanitized != expected {
		t.Fatalf(`DefaultKeySanitizer modified key to %s, expected %s`, sanitized, expected)
	}
}

func TestForbiddenChars(t *testing.T) {
	expected := "123456-_"
	objectKey := "123#@%$^&@456!-+=*_"
	sanitized := DefaultKeySanitizer(objectKey)
	if sanitized != expected {
		t.Fatalf(`DefaultKeySanitizer modified key to %s, expected %s`, sanitized, expected)
	}
}

func TestPercentEncodeChars(t *testing.T) {
	expected := "-_."
	objectKey := "-_.!#$&'()*+,/:;=?@[]"
	sanitized := DefaultKeySanitizer(objectKey)
	if sanitized != expected {
		t.Fatalf(`DefaultKeySanitizer modified key to %s, expected %s`, sanitized, expected)
	}
}

func TestAccentChars(t *testing.T) {
	expected := "aeiou"
	objectKey := "áêīòü"
	sanitized := DefaultKeySanitizer(objectKey)
	if sanitized != expected {
		t.Fatalf(`DefaultKeySanitizer modified key to %s, expected %s`, sanitized, expected)
	}
}

func TestOtherSeparators(t *testing.T) {
	expected := "test/test/test"
	objectKey := "test test test"
	sanitized := SanitizeObjectKey(objectKey, "/")
	if sanitized != expected {
		t.Fatalf(`SanitizeObjectKey modified key to %s, expected %s`, sanitized, expected)
	}
}

package chain

import (
	"reflect"
	"testing"
)

func TestParseBlacklistRule(t *testing.T) {
	awardID, users, err := parseBlacklistRule("101:user001, user002")
	if err != nil {
		t.Fatalf("parse blacklist rule: %v", err)
	}

	if awardID != 101 {
		t.Fatalf("expected award 101, got %d", awardID)
	}
	wantUsers := []string{"user001", "user002"}
	if !reflect.DeepEqual(users, wantUsers) {
		t.Fatalf("expected users %v, got %v", wantUsers, users)
	}
}

func TestParseBlacklistRuleInvalid(t *testing.T) {
	cases := []string{"101", "abc:user001"}

	for _, value := range cases {
		if _, _, err := parseBlacklistRule(value); err == nil {
			t.Fatalf("expected parse error for %q", value)
		}
	}
}

func TestParseWeightRule(t *testing.T) {
	got, err := parseWeightRule("4000:102,103 5000:104,105")
	if err != nil {
		t.Fatalf("parse weight rule: %v", err)
	}

	want := map[int]string{
		4000: "4000:102,103",
		5000: "5000:104,105",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected weight rule %v, got %v", want, got)
	}
}

func TestParseWeightRuleInvalid(t *testing.T) {
	cases := []string{"4000", "abc:102,103"}

	for _, value := range cases {
		if _, err := parseWeightRule(value); err == nil {
			t.Fatalf("expected parse error for %q", value)
		}
	}
}

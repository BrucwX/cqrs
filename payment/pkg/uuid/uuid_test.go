package uuid

import (
	"testing"

	guuid "github.com/google/uuid"
)

func TestNewInt64IsPositiveAndUnique(t *testing.T) {
	seen := make(map[int64]struct{}, 1000)
	for i := 0; i < 1000; i++ {
		id := NewInt64()
		if id <= 0 {
			t.Fatalf("NewInt64() = %d，必须是正数", id)
		}
		if _, dup := seen[id]; dup {
			t.Fatalf("NewInt64() 撞了：%d", id)
		}
		seen[id] = struct{}{}
	}
}

func TestNewStringIsValidUUID(t *testing.T) {
	s := NewString()
	if len(s) != 36 {
		t.Fatalf("NewString() 长度 = %d，想要 36：%q", len(s), s)
	}

	u, err := guuid.Parse(s)
	if err != nil {
		t.Fatalf("NewString() = %q 不是合法 UUID：%v", s, err)
	}
	if got := StringOf(u); got != s {
		t.Fatalf("StringOf(Parse(s)) = %q，想要 %q", got, s)
	}
	if Int64(u) <= 0 {
		t.Fatalf("Int64(解析回来的 UUID) 必须是正数")
	}
}

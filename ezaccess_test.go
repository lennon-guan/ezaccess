package ezaccess

import (
	"testing"
)

type DBConfig struct {
	DSN string
	KV  map[string]string
}

var c = struct {
	DB *DBConfig
}{
	DB: &DBConfig{
		DSN: "xxx",
		KV: map[string]string{
			"name": "zgg",
		},
	},
}

func TestGet(t *testing.T) {
	t.Log(MustGet[string](nil, &c, "DB.DSN"))
	t.Log(MustGet[string](nil, &c, "DB.KV.name"))
	t.Log(TryGet[int](nil, &c, "DB.KV.id"))
	t.Log(DefaultGet(nil, &c, "DB.KV.id", "N/A"))
}

func BenchmarkDirectlyGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if c.DB != nil {
			_ = c.DB.DSN
		}
		if c.DB != nil {
			_ = c.DB.KV["name"]
		}
		if c.DB != nil {
			_ = c.DB.KV["id"]
		}
		if c.DB != nil {
			_ = c.DB.KV["id"]
		}
	}
}

func BenchmarkUsingCache(b *testing.B) {
	var s PathStore
	_ = MustGet[string](&s, &c, "DB.DSN")
	_ = MustGet[string](&s, &c, "DB.KV.name")
	_, _ = TryGet[int](&s, &c, "DB.KV.id")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MustGet[string](&s, &c, "DB.DSN")
		_ = MustGet[string](&s, &c, "DB.KV.name")
		_, _ = TryGet[int](&s, &c, "DB.KV.id")
		_ = DefaultGet(&s, &c, "DB.KV.id", "N/A")
	}
}

func BenchmarkNoCache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = MustGet[string](nil, &c, "DB.c)N")
		_ = MustGet[string](nil, &c, "DB.KV.name")
		_, _ = TryGet[int](nil, &c, "DB.KV.id")
		_ = DefaultGet(nil, &c, "DB.KV.id", "N/A")
	}
}

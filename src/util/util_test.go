package util

import "testing"

func TestCompareVersion(t *testing.T) {
	// 正常长度比较
	if CompareVersion("1.0.0", "1.0.0") != 0 {
		t.Errorf("util.CompareVersion(\"1.0.0\", \"1.0.0\") != 0")
	}
	if CompareVersion("1.0.0", "1.0.1") != -1 {
		t.Errorf("util.CompareVersion(\"1.0.0\", \"1.0.1\") != -1")
	}
	if CompareVersion("1.0.1", "1.0.0") != 1 {
		t.Errorf("util.CompareVersion(\"1.0.1\", \"1.0.0\") != 1")
	}
	// 缺位比较
	if CompareVersion("1.0", "1.0.0") != 0 {
		t.Error("util.CompareVersion(\"1.0\", \"1.0.0\") != 0")
	}
	if CompareVersion("1.0.0", "1.0") != 0 {
		t.Error("util.CompareVersion(\"1.0.0\", \"1.0\") != 0")
	}
	if CompareVersion("1", "1.0.0") != 0 {
		t.Error("util.CompareVersion(\"1\", \"1.0.0\") != 0")
	}
	// 不同量级比较
	if CompareVersion("10.0.0", "1.0.0") != 1 {
		t.Error("util.CompareVersion(\"10.0.0\", \"1.0.0\") != 1")
	}
	// 其他字符比较
	if CompareVersion("1a.0.0", "1.0.0") != 0 {
		t.Error("util.CompareVersion(\"1a.0.0\", \"1.0.0\") != 0")
	}
	if CompareVersion("1.0.0a", "1.0.0") != 0 {
		t.Error("util.CompareVersion(\"1.0.0a\", \"1.0.0\") != 0")
	}
	if CompareVersion("a1.0.0", "1.0.0") != 0 {
		t.Error("util.CompareVersion(\"a1.0.0\", \"1.0.0\") != 0")
	}
	if CompareVersion("1.0.0", "a1.0.0") != 0 {
		t.Error("util.CompareVersion(\"1.0.0\", \"a1.0.0\") != 0")
	}

}
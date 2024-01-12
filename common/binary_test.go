package common_test

import (
	"encoding/binary"
	"testing"
	"tiangong/common"
)

var (
	void   = byte(0)
	_1void = []byte{void}
	_2void = []byte{void, void}
	_3void = []byte{void, void, void}
	_4void = []byte{void, void, void, void}
	_5void = []byte{void, void, void, void, void}
	_6void = []byte{void, void, void, void, void, void}
	_7void = []byte{void, void, void, void, void, void, void}
	_8void = []byte{void, void, void, void, void, void, void, void}

	value  byte
	_1byte []byte
	_2byte []byte
	_3byte []byte
	_4byte []byte
	_5byte []byte
	_6byte []byte
	_7byte []byte
	_8byte []byte
)

func Reset(v byte) {
	value = v
	_1byte = []byte{value}
	_2byte = []byte{value, value}
	_3byte = []byte{value, value, value}
	_4byte = []byte{value, value, value, value}
	_5byte = []byte{value, value, value, value, value}
	_6byte = []byte{value, value, value, value, value, value}
	_7byte = []byte{value, value, value, value, value, value, value}
	_8byte = []byte{value, value, value, value, value, value, value, value}
}

func testUint16(t *testing.T) bool {
	// 2bytes test
	expect := common.Uint16(_2byte)
	actual := binary.BigEndian.Uint16(_2byte)
	if expect != actual {
		t.Errorf("TestUint16 error, expect:%d,actual:%d", expect, actual)
		return false
	}

	// 1bytes test
	expect = common.Uint16(_1byte)
	actual = uint16(value)
	if expect != actual {
		t.Errorf("TestUint16 error, expect:%d,actual:%d", expect, actual)
		return false
	}
	return true
}

func testUint32(t *testing.T) bool {
	// 4bytes test
	expect := common.Uint32(_4byte)
	actual := binary.BigEndian.Uint32(_4byte)
	if expect != actual {
		t.Errorf("TestUint32 error, expect:%d,actual:%d", expect, actual)
		return false
	}

	// 3bytes test
	expect = common.Uint32(_3byte)
	actual = binary.BigEndian.Uint32(append(_1void, _3byte...))
	if expect != actual {
		t.Errorf("TestUint32 error, expect:%d,actual:%d", expect, actual)
		return false
	}

	// 2bytes test
	expect = common.Uint32(_2byte)
	actual = binary.BigEndian.Uint32(append(_2void, _2byte...))
	if expect != actual {
		t.Errorf("TestUint32 error, expect:%d,actual:%d", expect, actual)
		return false
	}

	// 1bytes test
	expect = common.Uint32(_1byte)
	actual = binary.BigEndian.Uint32(append(_3void, _1byte...))
	if expect != actual {
		t.Errorf("TestUint32 error, expect:%d,actual:%d", expect, actual)
		return false
	}
	return true
}

func testUint64(t *testing.T) bool {
	// 8bytes test
	expect := common.Uint64(_8byte)
	actual := binary.BigEndian.Uint64(_8byte)
	if expect != actual {
		t.Errorf("TestUint64 error, expect:%d,actual:%d", expect, actual)
		return false
	}

	// 7bytes test
	expect = common.Uint64(_7byte)
	actual = binary.BigEndian.Uint64(append(_1void, _7byte...))
	if expect != actual {
		t.Errorf("TestUint64 error, expect:%d,actual:%d", expect, actual)
		return false
	}

	// 6bytes test
	expect = common.Uint64(_6byte)
	actual = binary.BigEndian.Uint64(append(_2void, _6byte...))
	if expect != actual {
		t.Errorf("TestUint64 error, expect:%d,actual:%d", expect, actual)
		return false
	}

	// 5bytes test
	expect = common.Uint64(_5byte)
	actual = binary.BigEndian.Uint64(append(_3void, _5byte...))
	if expect != actual {
		t.Errorf("TestUint64 error, expect:%d,actual:%d", expect, actual)
		return false
	}

	// 4bytes test
	expect = common.Uint64(_4byte)
	actual = binary.BigEndian.Uint64(append(_4void, _4byte...))
	if expect != actual {
		t.Errorf("TestUint64 error, expect:%d,actual:%d", expect, actual)
		return false
	}

	// 3bytes test
	expect = common.Uint64(_3byte)
	actual = binary.BigEndian.Uint64(append(_5void, _3byte...))
	if expect != actual {
		t.Errorf("TestUint64 error, expect:%d,actual:%d", expect, actual)
		return false
	}

	// 2bytes test
	expect = common.Uint64(_2byte)
	actual = binary.BigEndian.Uint64(append(_6void, _2byte...))
	if expect != actual {
		t.Errorf("TestUint64 error, expect:%d,actual:%d", expect, actual)
		return false
	}

	// 1bytes test
	expect = common.Uint64(_1byte)
	actual = binary.BigEndian.Uint64(append(_7void, _1byte...))
	if expect != actual {
		t.Errorf("TestUint64 error, expect:%d,actual:%d", expect, actual)
		return false
	}
	return true
}

func TestUint(t *testing.T) {
	for i := 0; i < 255; i++ {
		Reset(byte(i))
		if !testUint16(t) {
			return
		}
		if !testUint32(t) {
			return
		}
		if !testUint64(t) {
			return
		}
	}
}

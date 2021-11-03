package geecache

type ByteView struct {
	b []byte
}

func (v ByteView) Len() int {
	return len(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

func (v ByteView) ByteSlice() []byte {
	return v.cloneByte()
}

func (v ByteView) cloneByte() []byte {
	temp := make([]byte, 0)
	copy(temp, v.b)
	return temp
}

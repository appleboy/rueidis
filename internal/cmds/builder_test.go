package cmds

import "testing"

func TestBuilder_Put(t *testing.T) {
	builder := NewBuilder(InitSlot)
	cs1 := builder.get()
	cs1.s = append(cs1.s, "1", "1", "1", "1", "1")
	builder.Put(cs1)
	cs2 := builder.get()
	if cs1 != cs2 {
		t.Fatalf("Put doesn't recycle the CommandSlice")
	}
	if len(cs2.s) != 0 {
		t.Fatalf("Put doesn't clean the CommandSlice")
	}
}

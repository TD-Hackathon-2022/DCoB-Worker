package executor

import (
	"bytes"
	"context"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func ShouldByteHas(actual interface{}, expected ...interface{}) string {
	data, _ := actual.([]byte)
	subData2, _ := expected[0].([]byte)
	if !bytes.HasPrefix(data, subData2) {
		return fmt.Sprintf("Data does not start with %v\n", expected[0])
	}
	return ""
}

func Test_miningCoins(t *testing.T) {
	Convey("The data starts with [0,0,0]", t, func() {
		data := miningCoins(3)
		So(data, ShouldByteHas, []byte{0, 0, 0})
	})
}

func TestNewHashCoins(t *testing.T) {
	Convey("test", t, func() {
		ctx, cancel := context.WithCancel(context.Background())
		ctx2 := context.WithValue(ctx, "var2", "2")
		ctx3 := context.WithValue(ctx2, "var3", "3")
		data := ctx3.Value("var2")
		go func() {
			time.Sleep(time.Second * 3)
			cancel()
		}()
		select {
		case <-ctx3.Done():
			fmt.Println(data)
			fmt.Println(ctx3.Err())
		}

		So("2", ShouldEqual, "2")
	})
}

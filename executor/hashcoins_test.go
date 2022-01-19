package executor

import (
	"context"
	"encoding/hex"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
	"time"
)

func ShouldByteHas(actual interface{}, expected ...interface{}) string {
	data, _ := actual.([]byte)
	subData2, _ := expected[0].([]byte)
	if !strings.HasPrefix(hex.EncodeToString(data), string(subData2)) {
		return fmt.Sprintf("Data does not start with %v\n", expected[0])
	}
	return ""
}

func Test_miningCoins(t *testing.T) {
	Convey("The data starts with [0,0,0]", t, func() {
		bits := 4
		data := miningCoins(bits)
		fmt.Println(hex.EncodeToString(data.([]byte)[:8]))
		So(data, ShouldByteHas, strings.Repeat("0", bits))
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

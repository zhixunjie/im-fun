package gen_id

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/zhixunjie/im-fun/pkg/gmodel"
	"math"
	"testing"
	"time"
)

var client = redis.NewClient(&redis.Options{
	Addr:     "127.0.0.1:6379",
	Password: "",
	DB:       0,
})

func TestIdMsg(t *testing.T) {
	ctx := context.Background()

	id1 := gmodel.NewUserComponentId(1001)
	id2 := gmodel.NewUserComponentId(1002)
	id3 := gmodel.NewGroupComponentId(10)

	// test1
	msgId, err := MsgId(ctx, client, id1, id2)
	fmt.Printf("单聊,msgId=%v,err=%v\n", msgId, err)

	// test2
	msgId, err = MsgId(ctx, client, id2, id1)
	fmt.Printf("单聊,msgId=%v,err=%v\n", msgId, err)

	// test3
	msgId, err = MsgId(ctx, client, id1, id3)
	fmt.Printf("群聊,msgId=%v,err=%v\n", msgId, err)

	// test4
	msgId, err = MsgId(ctx, client, id3, id1)
	fmt.Printf("群聊,msgId=%v,err=%v\n", msgId, err)
}

// 测试位数
// baseTimeStampOffset = 2023-02-22 02:31:47
func TestMsgIdBit(t *testing.T) {
	ctx := context.Background()

	// 结论：uint64 可以支持58年
	typ := "uint64"
	fmt.Printf("测试%v，最大值为: %v\n", typ, "18446744073709551615")
	testGen(ctx, "2080-01-01", typ) // 输出：17942884930000011001 <nil>
	testGen(ctx, "2081-01-01", typ) // 输出：18259108930000011001 <nil>	这里已经是极限了！！！
	testGen(ctx, "2082-01-01", typ) // 输出：18446744073709551615 <nil>
	testGen(ctx, "2083-01-01", typ) // 输出：18446744073709551615 <nil>
	testGen(ctx, "2084-01-01", typ) // 输出：18446744073709551615 <nil>

	// /////////////////////////////////////
	// 结论：int64 可以支持28年
	fmt.Println("======================================================")
	typ = "int64"
	fmt.Printf("测试%v，最大值为: %v\n", typ, math.MaxInt64)
	testGen(ctx, "2050-01-01", typ) // 输出：8476036930000011001 <nil>
	testGen(ctx, "2051-01-01", typ) // 输出：8791396930000011001 <nil>
	testGen(ctx, "2052-01-01", typ) // 输出：9106756930000011001 <nil> 这里已经是极限了！！！
	testGen(ctx, "2053-01-01", typ) // 输出：9223372036854775807 <nil>
	testGen(ctx, "2054-01-01", typ) // 输出：9223372036854775807 <nil>

	// /////////////////////////////////////
	// 结论：uint32 可以支持20年
	fmt.Println("======================================================")
	typ = "uint32"
	fmt.Printf("测试%v，最大值为: %v\n", typ, math.MaxUint32)
	testGen(ctx, "2043-01-01", typ) // 输出：922763001 <nil> 这里已经是极限了！！！
	testGen(ctx, "2044-01-01", typ) // 输出：4284366585 <nil>
	testGen(ctx, "2045-01-01", typ) // 输出：1959935737 <nil>
	testGen(ctx, "2046-01-01", typ) // 输出：1026572025 <nil>
	testGen(ctx, "2047-01-01", typ) // 输出：93208313 <nil>

	// /////////////////////////////////////
	// 结论：int32 完全玩不了
	fmt.Println("======================================================")
	typ = "int32"
	fmt.Printf("测试%v，最大值为: %v\n", typ, math.MaxInt32)
	testGen(ctx, "2024-01-01", typ) // 输出：8476036930000011001 <nil> 这里已经是极限了！！！
	testGen(ctx, "2025-01-01", typ) // 输出：8791396930000011001 <nil>
	testGen(ctx, "2026-01-01", typ) // 输出：9106756930000011001 <nil>
	testGen(ctx, "2027-01-01", typ) // 输出：9223372036854775807 <nil>
	testGen(ctx, "2028-01-01", typ) // 输出：9223372036854775807 <nil>

}

func testGen(ctx context.Context, dateStr, typ string) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	date := getDate(dateStr, loc)
	res, err := genMsgId1(ctx, client, 1001, date, typ)
	fmt.Println("date: ", date, res, err)
}

// 解析日期字符串为时间对象
func getDate(dateStr string, loc *time.Location) (date time.Time) {
	var err error
	date, err = time.Parse("2006-01-02", dateStr)
	if err != nil {
		fmt.Println("解析日期失败:", err)
		return
	}
	date = date.In(loc)

	return
}

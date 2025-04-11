package snowflake

import (
	sf "github.com/bwmarrin/snowflake"
	"time"
)

var (
	node *sf.Node
)

func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime) // 将传入的开始时间解析为time固定格式
	if err != nil {
		return err
	}
	//	// 步骤3：计算起始时间的毫秒时间戳
	//	// 1. st.UnixNano() 获取纳秒级时间戳（从1970-01-01起）
	//	// 2. 除以 1,000,000 转换为毫秒（Snowflake算法通常使用毫秒）
	sf.Epoch = st.UnixNano() / 1000000
	// 传入机器ID（用于分布式系统中不同节点的唯一标识）
	node, err = sf.NewNode(machineID)
	return err
}

// 使用node创建id
func GetID() int64 {
	return node.Generate().Int64()
}

//func main() {
//	if err := Init("2023-09-09", 1); err != nil {
//		fmt.Printf("Init err: %v\n", err)
//		return
//	}
//	id := GetID()
//	fmt.Printf("id: %d\n", id)
//}

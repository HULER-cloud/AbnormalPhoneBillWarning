package email

import (
	"AbnormalPhoneBillWarning/DataAnalysis/RabbitMQ"
	"fmt"
	"github.com/streadway/amqp"
	"testing"
)

func Test_handler_MultipleSend(t *testing.T) {
	type args struct {
		delivery amqp.Delivery
		map_pub  map[string]*RabbitMQ.RabbitMQ
		pub_name []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "1", args: args{
			delivery: amqp.Delivery{
				Body: []byte(`{"user_id":1,"email":"2799591178@qq.com","missionID":0,"mission":{"timeStamp":"2024/06/25 19:12:31","balance":56.82,"cost":0,"abnormal_consumption":null}}`),
			},
			map_pub:  nil,
			pub_name: nil,
		}}, // 余额异常测试用例
		{name: "2", args: args{
			delivery: amqp.Delivery{
				Body: []byte(`{"user_id":1,"email":"2799591178@qq.com","missionID":1,"mission":{"timeStamp":"2024/06/25 19:12:31","balance":56.82,"cost":32,"abnormal_consumption":[{"consumption_name": "业务1","consumption_amount": 15},{"consumption_name": "业务2","consumption_amount": 17}]}}`),
			},
			map_pub:  nil,
			pub_name: nil,
		}}, // 消费异常测试用例
		{name: "3", args: args{
			delivery: amqp.Delivery{
				Body: []byte(`{"user_id":1,"email":"2799591178@qq.com","missionID":1,"mission":{"timeStamp":"2024/06/25 19:12:31","balance":156.82,"cost":0.5,"abnormal_consumption":[{"consumption_name": "业务1","consumption_amount": 15},{"consumption_name": "业务2","consumption_amount": 17}]}}`),
			},
			map_pub:  nil,
			pub_name: nil,
		}}, // 没有异常（其实没有异常就不该传到这里来的，只能强制选一个异常属性）
		{name: "4", args: args{
			delivery: amqp.Delivery{
				Body: []byte(`{"userid":1,"email":"2799591178@qq.com","missionID":1,"mission":{"timeStamp":"2024/06/25 19:12:31","balance":56.82,"cost":32,"abnormal_consumption":[{"consumption_name": "业务1","consumption_amount": 15},{"consumption_name": "业务2","consumption_amount": 17}]}}`),
			},
			map_pub:  nil,
			pub_name: nil,
		}}, // 参数名错误（缺失某个属性）
		{name: "5", args: args{
			delivery: amqp.Delivery{
				Body: []byte(`{"user_id":1,"extra":231,"email":"2799591178@qq.com","missionID":1,"mission":{"timeStamp":"2024/06/25 19:12:31","balance":56.82,"cost":32,"abnormal_consumption":[{"consumption_name": "业务1","consumption_amount": 15},{"consumption_name": "业务2","consumption_amount": 17}]}}`),
			},
			map_pub:  nil,
			pub_name: nil,
		}}, // 多出某个属性
		{name: "6", args: args{
			delivery: amqp.Delivery{
				Body: []byte(`{"user_id":1,"","email":"2799591178@qq.com","missionID":1,"mission":{"timeStamp":"2024/06/25 19:12:31","balance":56.82,"cost":32,"abnormal_consumption":[{"consumption_name": "业务1","consumption_amount": 15},{"consumption_name": "业务2","consumption_amount": 17}]}}`),
			},
			map_pub:  nil,
			pub_name: nil,
		}}, // json格式错误
	}
	initdb()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := handler_MultipleSend(tt.args.delivery, tt.args.map_pub, tt.args.pub_name); (err != nil) != tt.wantErr {
				t.Errorf("handler_MultipleSend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Benchmark_handler_MultipleSend(b *testing.B) {
	b.StopTimer() //调用该函数停止压力测试的时间计数

	initdb()
	//做一些初始化的工作,例如读取文件数据,数据库连接之类的,
	//这样这些时间不影响我们测试函数本身的性能

	b.StartTimer() //重新开始时间
	for i := 0; i < b.N; i++ {
		fmt.Println(i)
		delivery := amqp.Delivery{
			Body: []byte(`{"user_id":1,"email":"2799591178@qq.com","missionID":0,"mission":{"timeStamp":"2024/06/25 19:12:31","balance":56.82,"cost":0,"abnormal_consumption":null}}`),
		}
		var map_pub map[string]*RabbitMQ.RabbitMQ
		var pub_name []string
		handler_MultipleSend(delivery, map_pub, pub_name)
	}
}

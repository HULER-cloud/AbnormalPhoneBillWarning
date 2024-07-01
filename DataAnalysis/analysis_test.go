package dataanalysis

import (
	"AbnormalPhoneBillWarning/DataAnalysis/RabbitMQ"
	"fmt"
	"github.com/streadway/amqp"
	"testing"
)

func Test_initdb(t *testing.T) {
	type test struct {
		name string
	}
	tests := []test{
		// TODO: Add test cases.
		{name: "1"},
		{name: "2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initdb()
		})
	}
}

func Benchmark_initdb(b *testing.B) {
	b.StopTimer() //调用该函数停止压力测试的时间计数

	//做一些初始化的工作,例如读取文件数据,数据库连接之类的,
	//这样这些时间不影响我们测试函数本身的性能

	b.StartTimer() //重新开始时间
	for i := 0; i < b.N; i++ {
		fmt.Println(i)
		initdb()
	}
}

func Test_handler_DataAnalysis(t *testing.T) {
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
				Body: []byte(`{
  "user_id": 1,
  "balance": 56.82,
  "timeStamp": "2024/06/25 19:12:31",
  "consumption_condition": {
    "consumption_amount": 3.15,
    "consumption_set": [
      {
        "subscription_name": "20元包20G国内流量包",
        "subscription_amount": 0.0
      },
      {
        "subscription_name": "iFree卡（2016版）",
        "subscription_amount": 3.0
      },
      {
        "subscription_name": "国内通话费",
        "subscription_amount": 0.15
      }
    ]
  }
}`),
			},
			map_pub:  nil,
			pub_name: nil,
		}}, // 异常1、2
		{name: "2", args: args{
			delivery: amqp.Delivery{
				Body: []byte(`{
  "user_id": 1,
  "balance": 156.82,
  "timeStamp": "2024/06/25 19:12:31",
  "consumption_condition": {
    "consumption_amount": 3.15,
    "consumption_set": [
      {
        "subscription_name": "20元包20G国内流量包",
        "subscription_amount": 0.0
      },
      {
        "subscription_name": "iFree卡（2016版）",
        "subscription_amount": 3.0
      },
      {
        "subscription_name": "国内通话费",
        "subscription_amount": 0.15
      }
    ]
  }
}`),
			},
			map_pub:  nil,
			pub_name: nil,
		}}, // 异常2
		{name: "3", args: args{
			delivery: amqp.Delivery{
				Body: []byte(`{
  "user_id": 1,
  "balance": 56.82,
  "timeStamp": "2024/06/25 19:12:31",
  "consumption_condition": {
    "consumption_amount": 0.15,
    "consumption_set": [
      {
        "subscription_name": "20元包20G国内流量包",
        "subscription_amount": 0.0
      },
      {
        "subscription_name": "iFree卡（2016版）",
        "subscription_amount": 3.0
      },
      {
        "subscription_name": "国内通话费",
        "subscription_amount": 0.15
      }
    ]
  }
}`),
			},
			map_pub:  nil,
			pub_name: nil,
		}}, // 异常1
		{name: "4", args: args{
			delivery: amqp.Delivery{
				Body: []byte(`{
  "user_id": 1,
  "balance": 156.82,
  "timeStamp": "2024/06/25 19:12:31",
  "consumption_condition": {
    "consumption_amount": 0.15,
    "consumption_set": [
      {
        "subscription_name": "20元包20G国内流量包",
        "subscription_amount": 0.0
      },
      {
        "subscription_name": "iFree卡（2016版）",
        "subscription_amount": 3.0
      },
      {
        "subscription_name": "国内通话费",
        "subscription_amount": 0.15
      }
    ]
  }
}`),
			},
			map_pub:  nil,
			pub_name: nil,
		}}, // 无异常
		{name: "5", args: args{
			delivery: amqp.Delivery{
				Body: []byte(`{
  "userid": 1,
  "balance": 156.82,
  "timeStamp": "2024/06/25 19:12:31",
  "consumption_condition": {
    "consumption_amount": 0.15,
    "consumption_set": [
      {
        "subscription_name": "20元包20G国内流量包",
        "subscription_amount": 0.0
      },
      {
        "subscription_name": "iFree卡（2016版）",
        "subscription_amount": 3.0
      },
      {
        "subscription_name": "国内通话费",
        "subscription_amount": 0.15
      }
    ]
  }
}`),
			},
			map_pub:  nil,
			pub_name: nil,
		}}, // 参数名错误（缺失某个属性）
		{name: "6", args: args{
			delivery: amqp.Delivery{
				Body: []byte(`{
  "user_id": 1,
  "extra":123,
  "balance": 156.82,
  "timeStamp": "2024/06/25 19:12:31",
  "consumption_condition": {
    "consumption_amount": 0.15,
    "consumption_set": [
      {
        "subscription_name": "20元包20G国内流量包",
        "subscription_amount": 0.0
      },
      {
        "subscription_name": "iFree卡（2016版）",
        "subscription_amount": 3.0
      },
      {
        "subscription_name": "国内通话费",
        "subscription_amount": 0.15
      }
    ]
  }
}`),
			},
			map_pub:  nil,
			pub_name: nil,
		}}, // 多出某个属性
		{name: "7", args: args{
			delivery: amqp.Delivery{
				Body: []byte(`{
  "user_id": 1
  "balance": 156.82,
  "timeStamp": "2024/06/25 19:12:31",
  "consumption_condition": {
    "consumption_amount": 0.15,
    "consumption_set": [
      {
        "subscription_name": "20元包20G国内流量包",
        "subscription_amount": 0.0
      },
      {
        "subscription_name": "iFree卡（2016版）",
        "subscription_amount": 3.0
      },
      {
        "subscription_name": "国内通话费",
        "subscription_amount": 0.15
      }
    ]
  }
}`),
			},
			map_pub:  nil,
			pub_name: nil,
		}}, // json格式错误
	}
	initdb()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := handler_DataAnalysis(tt.args.delivery, tt.args.map_pub, tt.args.pub_name); (err != nil) != tt.wantErr {
				t.Errorf("handler_DataAnalysis() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Benchmark_handler_DataAnalysis(b *testing.B) {
	b.StopTimer() //调用该函数停止压力测试的时间计数

	initdb()
	//做一些初始化的工作,例如读取文件数据,数据库连接之类的,
	//这样这些时间不影响我们测试函数本身的性能

	b.StartTimer() //重新开始时间
	for i := 0; i < b.N; i++ {
		fmt.Println(i)
		delivery := amqp.Delivery{
			Body: []byte(`{
  "user_id": 1,
  "balance": 56.82,
  "timeStamp": "2024/06/25 19:12:31",
  "consumption_condition": {
    "consumption_amount": 3.15,
    "consumption_set": [
      {
        "subscription_name": "20元包20G国内流量包",
        "subscription_amount": 0.0
      },
      {
        "subscription_name": "iFree卡（2016版）",
        "subscription_amount": 3.0
      },
      {
        "subscription_name": "国内通话费",
        "subscription_amount": 0.15
      }
    ]
  }
}`),
		}
		var map_pub map[string]*RabbitMQ.RabbitMQ
		var pub_name []string
		handler_DataAnalysis(delivery, map_pub, pub_name)
	}
}

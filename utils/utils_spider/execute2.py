import pika
import sys
import json

province=sys.argv[1]
userID=sys.argv[2]
phoneNumber=sys.argv[3]
password=sys.argv[4]

data={'province':province,"userid":int(userID),"phone_number":phoneNumber,"password":password}

print(data)

# 建立到RabbitMQ服务器的连接
user_info = pika.PlainCredentials('guest', 'guest')
connection = pika.BlockingConnection(pika.ConnectionParameters('localhost', 5672, '/', user_info))
channel = connection.channel()
# 声明队列（如果不存在的话）
channel.queue_declare(queue='TimerAwakeScript',durable=True)

def send_json(myjson):
    # 将Python字典转换为JSON字符串
    json_data = json.dumps(myjson, ensure_ascii=False)  # ensure_ascii=False用于支持中文显示
    # 发布JSON消息到RabbitMQ
    channel.basic_publish(exchange='',
                          routing_key='PythonCrawlerResult',
                          body=json_data.encode('utf-8'))  # 需要将字符串编码为字节流
    print(f"已发送JSON消息：{json_data}")

# 示例用法：
# my_data = {"姓名": "张三", "年龄": 30, "城市": "北京"}
# send_json(my_data)

send_json(data)

# 关闭与RabbitMQ的连接
connection.close()
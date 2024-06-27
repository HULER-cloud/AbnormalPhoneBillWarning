# 客户名称：任**
# 账单周期：2024年06月01日-2024年06月24日
# 打印日期：2024年06月24日
# 本期实时消费(元)
# 3.15
# 20元包20G国内流量包及叠加包
# 手机号:17353825020
# 原价 减免 实际消费
# 套餐及固定费
# 20元包20G国内流量包 20.00 -20.00 0.00
# iFree卡（2016版） 3.00 3.00
# 套餐外语音通信费
# 国内通话费 0.15 0.15
# 小计 23.15 -20.00 3.15
# 合计 23.15 -20.00 3.15
import base64
import sys
import time
from io import BytesIO
import cv2
import numpy as np
from PIL import Image
from selenium import webdriver
from selenium.webdriver import ActionChains
from selenium.webdriver.common.by import By
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.support.wait import WebDriverWait
from selenium.webdriver.firefox.options import Options
import json
from datetime import datetime
import pika

# 登录的网站
url = 'https://login.189.cn/web/login'

# 登录的账号和密码
# username = "17353825020"  # 山东号码
# password = ""
userID=sys.argv[1]
username=sys.argv[2]
password=sys.argv[3]
# low_threshold = 10  # 低阈值为10元
# high_threshold = 500  # 高阈值为500元


class CrackSlider():

    def __init__(self):  # 通过浏览器截图，识别验证码中缺口位置，获取需要滑动距离，并破解滑动验证码
        super(CrackSlider, self).__init__()
        options = Options()
        options.headless = True    # 设置Firefox浏览器为无头模式（即不直接打开浏览器，只要对应页面的html代码）
#         options.add_argument("--headless") #设置火狐为headless无界面模式
        self.opts = options
        self.driver = webdriver.Firefox(options=self.opts)
        self.url = url
        self.wait = WebDriverWait(self.driver, 10)  # 设置全局等待时间

    def get_pic(self):  # 获取滑块和背景图片
        self.driver.get(self.url)  # 登录网站

        # 通过class属性选择元素
        self.driver.find_element(By.ID, 'txtAccount').send_keys(username)  # gLFyf

        # 为了使密码输入框变成txtPassword
        trigger_element = self.driver.find_element(By.ID, 'txtShowPwd')
        trigger_element.click()
        self.driver.find_element(By.ID, 'txtPassword').send_keys(password)

        # 勾选同意协议
        trigger_element = self.driver.find_element(By.ID, 'u2_input')
        trigger_element.click()
        time.sleep(1)

        # 点击登录
        trigger_element = self.driver.find_element(By.ID, 'loginbtn')
        trigger_element.click()
        time.sleep(2)

        # 下载相应的滑块和背景图片
        target_link = self.driver.find_element(By.CLASS_NAME, "code-back-img").get_attribute('src')
        template_link = self.driver.find_element(By.CLASS_NAME, "code-front-img").get_attribute('src')
        target_img = Image.open(BytesIO(base64.b64decode(target_link.split(',')[1])))
        template_img = Image.open(BytesIO(base64.b64decode(template_link.split(',')[1])))
        target_img.save('target.jpg')
        template_img.save('template.png')

    # 使用ActionChains来模拟人滑动滑块的操作

    def crack_slider(self, distance):
        slider = self.wait.until(EC.element_to_be_clickable((By.CLASS_NAME, 'code-btn-img.code-btn-m')))
        ActionChains(self.driver).click_and_hold(slider).perform()
        ActionChains(self.driver).move_by_offset(xoffset=distance, yoffset=0).perform()
        time.sleep(1)
        ActionChains(self.driver).release().perform()
        time.sleep(1)

        # 等待页面跳转
        WebDriverWait(self.driver, 10).until(EC.url_changes(self.url))
        # self.driver.switch_to.window(self.driver.window_handles[1])  # 切换当前页面标签

        # 获取当前页面的URL地址
        current_url = self.driver.current_url
        # print("跳转后的页面URL地址111:", current_url)

        # 如果地址没有跳转说明滑块出错，重新执行程序
        if current_url == url:
           sys.exit(1)

        self.url = current_url
        self.get_sd()

        # 停留10秒后关闭浏览器
        time.sleep(10)
        self.driver.quit()  # 关闭浏览器
        return 0

    def get_sd(self):
        global telephone_text  # 实时话费模块
        global query_result_text  # 实时话费模块
        global balance_info_text  # 余额查询模块
        global balance  # 余额查询模块中的余额
        global combo_info_text  # 套餐使用情况查询模块
        # 选择并点击按钮
        time.sleep(2)
        button = self.driver.find_element(By.LINK_TEXT, '话费查询')
        button_url = button.get_attribute("href")
        # print('按钮的URL:', button_url)
        self.driver.get(button_url)  # 如果按钮有href属性，则使用driver.get访问该链接
        # 获取新页面的URL
        new_page_url = self.driver.current_url
        # print("新页面的URL:", new_page_url)
        '''实时话费板块'''
        self.driver.find_element(By.ID, 'aa').click()
        time.sleep(3)
        telephone = self.driver.find_element(By.ID, 'showUserID')
        telephone_text = telephone.text
        # print('telephone_text: ', telephone.text)  # 电话号码
        query_result = self.driver.find_element(By.CSS_SELECTOR,
                                                'html body div.myorder_nav div.myorder_nav div.myorder div.myorder_right div.selser_common_wrap div div div#resultHtml.orderlist_t div.com_box div#queryResult')
        query_result_text = query_result.text
        # print('query_result_text: ', query_result.text)  # 包含套餐及固定费
        '''余额查询模块'''
        self.driver.find_element(By.ID, 'queryBalance').click()
        time.sleep(3)
        balance_info = self.driver.find_element(By.CSS_SELECTOR,
                                                'html body div.myorder_nav div.myorder_nav div.myorder div.myorder_right div.selser_common_wrap div div div#resultHtml.orderlist_t div.com_box div#queryResult table.com_table')
        balance_info_text = balance_info.text
        balance = balance_info_text.split("账户余额（元）：")[1].split()[0]  # 分割文本并获取账户余额信息
        # print("账户余额为: ", balance) # 输出账户余额
        # print('余额查询文本: ', balance_info_text)
        '''套餐使用情况查询模块'''
        self.driver.find_element(By.ID, 'queryOfferInfo').click()
        time.sleep(3)
        combo_info = self.driver.find_element(By.CSS_SELECTOR,
                                              'html body div.myorder_nav div.myorder_nav div.myorder div.myorder_right div.selser_common_wrap div div div#resultHtml.orderlist_t div.com_box div#queryResult table.com_table')
        combo_info_text = combo_info.text
        # print('套餐使用情况查询: ')
        # print(combo_info_text)

        txt_content=("电话号码: " + "\n" + telephone_text + "\n\n"
                     +"实时话费: " + "\n" + query_result_text + "\n\n"
                     +"余额信息: " + "\n" + balance_info_text + "\n\n"
                     +"具体余额: " + "\n" + balance + "\n\n"
                     +"套餐使用情况: " + combo_info_text)

        # with open(f"./store/{username}.txt", "w", encoding="utf-8") as file:  # 将全局变量写入文件
        #     file.write("电话号码: " + "\n" + telephone_text + "\n\n")
        #     file.write("实时话费: " + "\n" + query_result_text + "\n\n")
        #     file.write("余额信息: " + "\n" + balance_info_text + "\n\n")
        #     file.write("具体余额: " + "\n" + balance + "\n\n")
        #     file.write("套餐使用情况: " + combo_info_text)

        # print("全局变量已写入文件 query_result.txt")
        txt2json(txt_content)

def is_float(string):
    try:
        float(string)
        return True
    except ValueError:
        return False

def txt2json(input_file):
    # 读取文本文件内容
    # with open(input_file, 'r', encoding='utf-8') as f:
    #     txt_content = f.readlines()
    txt_content=input_file.split('\n')

    # 解析文本内容，提取需要的信息
    data = {}
    consumption_amount = None
    subscription_set = []
    skip_line = False
    for idx, line in enumerate(txt_content):
        if skip_line:
            skip_line = False
            continue
        if "本期实时消费" in line:
            consumption_amount = float(txt_content[idx + 1].split()[-1])
        elif "套餐及固定费" in line:
            for i in range(idx + 1, len(txt_content)):
                line_parts = txt_content[i].split()
                if "综合信息服务费" in line_parts[0]:
                    continue
                elif "小计" in line_parts[0]:
                    break
                elif not is_float(line_parts[-1]):
                    continue
                subscription_name = line_parts[0]
                subscription_amount = float(line_parts[-1])

                subscription_set.append ( {
                    "subscription_name": subscription_name,
                    "subscription_amount": subscription_amount
                })
                if len(line_parts) < 4:
                    skip_line = True
        elif "小计" in line:
            break

    timestamp = datetime.now().strftime("%Y/%m/%d %H:%M:%S")  # 获取当前时间戳
    # 构造 JSON 数据
    data["user_id"] = int(userID)
    data["balance"] = float(balance)
    data["timeStamp"] = timestamp
    data["consumption_condition"] = {
        "consumption_amount": consumption_amount,
        "consumption_set":
            subscription_set
    }
    # print(json.dumps(data))
    # 写入 JSON 文件
    # with open(output_file, 'w', encoding='utf-8') as f:
    #     json.dump(data, f, ensure_ascii=False, indent=2)


    # 建立到RabbitMQ服务器的连接
    user_info = pika.PlainCredentials('guest', 'guest')
    connection = pika.BlockingConnection(pika.ConnectionParameters('localhost', 5672, '/', user_info))
    channel = connection.channel()
    # 声明队列（如果不存在的话）
    channel.queue_declare(queue='PythonCrawlerResult',durable=True)
    def send_json(myjson):
        # 将Python字典转换为JSON字符串
        print(1)
        json_data = json.dumps(myjson, ensure_ascii=False)  # ensure_ascii=False用于支持中文显示
        print(2)
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





def add_alpha_channel(img):
    """ 为jpg图像添加alpha通道 """
    r_channel, g_channel, b_channel = cv2.split(img)  # 剥离jpg图像通道
    alpha_channel = np.ones(b_channel.shape, dtype=b_channel.dtype) * 255  # 创建Alpha通道
    img_new = cv2.merge((r_channel, g_channel, b_channel, alpha_channel))  # 融合通道
    return img_new


def handel_img(img):  # 进行图片的边缘检测
    imgGray = cv2.cvtColor(img, cv2.COLOR_RGBA2GRAY)  # 转灰度图
    imgBlur = cv2.GaussianBlur(imgGray, (5, 5), 1)  # 高斯模糊
    imgCanny = cv2.Canny(imgBlur, 60, 60)  # Canny算子边缘检测
    return imgCanny


def match(img_jpg_path, img_png_path):  # 两张图片进行相似度匹配
    # 读取图像
    img_jpg = cv2.imread(img_jpg_path, cv2.IMREAD_UNCHANGED)
    img_png = cv2.imread(img_png_path, cv2.IMREAD_UNCHANGED)

    # 判断jpg图像是否已经为4通道
    if img_jpg.shape[2] == 3:
        img_jpg = add_alpha_channel(img_jpg)
    img = handel_img(img_jpg)
    small_img = handel_img(img_png)
    res_TM_CCOEFF_NORMED = cv2.matchTemplate(img, small_img, 3)
    value = cv2.minMaxLoc(res_TM_CCOEFF_NORMED)
    value = value[3][0]  # 获取到移动距离

    return value


# def check_balance_threshold():  # 检查话费余额是否低于某个阈值
#     if balance < low_threshold:
#         pass


# 1. 打开Firefoxdriver，试试下载图片
cs = CrackSlider()
cs.get_pic()

# 2. 对比图片，计算距离
img_jpg_path = 'target.jpg'
img_png_path = 'template.png'
distance = match(img_jpg_path, img_png_path)
distance = distance + 9.4

# 3. 移动
cs.crack_slider(distance)

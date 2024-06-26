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

# 登录的网站
url = 'https://gd.189.cn/common/newLogin.html'

# 登录的账号和密码
# username = "19047188667"  # 广东号码
# password = "109238"
username=sys.argv[1]
password=sys.argv[2]
# low_threshold = 10  # 低阈值为10元
# high_threshold = 500  # 高阈值为500元


class CrackSlider():

    def __init__(self):  # 通过浏览器截图，识别验证码中缺口位置，获取需要滑动距离，并破解滑动验证码
        super(CrackSlider, self).__init__()
        options = Options()
        options.headless = True    # 设置Firefox浏览器为无头模式（即不直接打开浏览器，只要对应页面的html代码）
        options.add_argument("--headless") #设置火狐为headless无界面模式
        self.opts = options
        self.driver = webdriver.Firefox(options=self.opts)
        self.url = url
        self.wait = WebDriverWait(self.driver, 10)  # 设置全局等待时间

    def get_pic(self):  # 获取滑块和背景图片
        self.driver.get(self.url)  # 登录网站
        time.sleep(5)
        tab_items = self.driver.find_elements(By.CLASS_NAME, 'tabItem')
        # 遍历所有找到的元素，根据文本内容选择特定元素
        for item in tab_items:
            if item.text == "密码登录":  # 替换为你感兴趣的文本内容
                item.click()
                # print("点击了包含 '密码登录' 文本的 'tabItem' 元素")
                # break
        # else:
        #     print("没有找到包含 '密码登录' 文本的 'tabItem' 元素")

        self.driver.find_element(By.XPATH, '//input[@class="inputItem" and @placeholder="广东电信手机/固话/宽带/IPTV账号"]').send_keys(username)

        self.driver.find_element(By.ID, 'accCode').send_keys(password)

        # # 勾选同意协议
        checkbox = self.driver.find_element(By.XPATH, '//input[@type="checkbox" and @autocomplete="off"]')

        # 点击复选框
        checkbox.click()
        time.sleep(1)

        # # 点击登录
        trigger_element = self.driver.find_element(By.CLASS_NAME, 'btn-toLogin')
        trigger_element.click()
        time.sleep(2)

        # # 下载相应的滑块和背景图片
        # element = driver.find_element_by_xpath("//img[@data-v-3e50952d]")

        target_link = self.driver.find_element(By.CLASS_NAME, "verify-img-panel").find_element(By.TAG_NAME,'img').get_attribute('src')
        template_link = self.driver.find_element(By.CLASS_NAME, "verify-sub-block").find_element(By.TAG_NAME,'img').get_attribute('src')
        # print(target_link)
        # print(template_link)
        # target_img = Image.open(BytesIO(base64.b64decode(target_link.split(',')[1])))
        # template_img = Image.open(BytesIO(base64.b64decode(template_link.split(',')[1])))
        with open('target.jpg', 'wb') as f:
            target_img_data = base64.b64decode(target_link.split(',')[1])
            f.write(target_img_data)

        with open('template.png', 'wb') as f:
            template_img_data = base64.b64decode(template_link.split(',')[1])
            f.write(template_img_data)

        # 加载图像文件为PIL图像对象
        target_img = Image.open('target.jpg')
        template_img = Image.open('template.png')

        # 检查图像模式并转换为RGB模式（如果需要）
        if target_img.mode in ('RGBA', 'LA'):
            target_img = target_img.convert('RGB')

        if template_img.mode in ('RGBA', 'LA'):
            template_img = template_img.convert('RGB')

        # 保存图像为JPEG格式
        target_img.save('target.jpg', format='JPEG')
        template_img.save('template.png', format='PNG')
        # target_img.save('target.jpg')
        # template_img.save('template.png')

    # 使用ActionChains来模拟人滑动滑块的操作

    def crack_slider(self, distance):
        slider = self.wait.until(EC.element_to_be_clickable((By.CLASS_NAME, 'verify-move-block')))
        ActionChains(self.driver).click_and_hold(slider).perform()
        ActionChains(self.driver).move_by_offset(xoffset=distance, yoffset=0).perform()
        time.sleep(3)
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
        self.driver.get(self.url)
        self.get_gd()

        # 停留10秒后关闭浏览器
        time.sleep(10)
        self.driver.quit()  # 关闭浏览器
        return 0

    def get_gd(self):
        global telephone_text  # 实时话费模块
        global query_result_text  # 实时话费模块
        global balance_info_text  # 余额查询模块
        global balance  # 余额查询模块中的余额
        global combo_info_text  # 套餐使用情况查询模块
        # 选择并点击按钮
        # button = self.driver.find_element(By.LINK_TEXT, '话费查询')
        # button_url = button.get_attribute("href")
        # # print('按钮的URL:', button_url)
        # self.driver.get(button_url)  # 如果按钮有href属性，则使用driver.get访问该链接
        # # 获取新页面的URL
        # new_page_url = self.driver.current_url
        new_page_url="https://gd.189.cn/service/myhome"
        self.driver.get(new_page_url)
        # print("新页面的URL:", new_page_url)
        '''实时话费板块'''
        # self.driver.find_element(By.ID, 'aa').click()
        time.sleep(2)
        telephone = self.driver.find_element(By.ID, 'loginNumberN')
        telephone_text = telephone.text
        # print('telephone_text: ', telephone.text)  # 电话号码

        balance = self.wait.until(EC.presence_of_element_located((By.CSS_SELECTOR, '.two.BalanceNum'))).text
        # print('balance: ', balance)  # 余额

        query_result_text=self.wait.until(EC.presence_of_element_located((By.CSS_SELECTOR, '.two.RealchargeNum'))).text

        # print('query_result: ', query_result_text)  # 话费

        #套餐查询
        # button=self.driver.find_element(By.XPATH,'//a[contains(text(), "套餐")]')
        # button_url = button.get_attribute("href")
        # print('按钮的URL:', button_url)
        # self.driver.get(button_url)  # 如果按钮有href属性，则使用driver.get访问该链接
        # time.sleep(2)
        # button=self.driver.find_element(By.XPATH,'//a[contains(text(), "套餐用量")]')
        # button_url = button.get_attribute("href")
        # print('按钮的URL:', button_url)
        # self.driver.get(button_url)  # 如果按钮有href属性，则使用driver.get访问该链接

        # combo_info_text=self.wait.until(EC.presence_of_element_located((By.CLASS_NAME,'tableContent'))).text
        button=self.driver.find_element(By.XPATH,'//a[contains(text(), "话费")]')
        button_url = button.get_attribute("href")
        # print('按钮的URL:', button_url)
        self.driver.get(button_url)  # 如果按钮有href属性，则使用driver.get访问该链接
        time.sleep(2)
        button=self.driver.find_element(By.XPATH,'//a[contains(text(), "话费账单")]')
        button_url = button.get_attribute("href")
        # print('按钮的URL:', button_url)
        self.driver.get(button_url)  # 如果按钮有href属性，则使用driver.get访问该链接
        time.sleep(3)
        # button=self.driver.find_element(By.XPATH,'//a[contains(text(), "查询")]')
        button = WebDriverWait(self.driver, 10).until(
            EC.element_to_be_clickable((By.XPATH, '//button[contains(text(), "查询")]'))
        )
        button.click()
        time.sleep(2)
        combo_info_text=self.driver.find_element(By.CLASS_NAME,"billItemFor").text
        # print('套餐使用情况: ' ,combo_info_text)
        txt_content = ("电话号码: " + "\n" + telephone_text + "\n\n"
                       + "实时话费: " + "\n" + query_result_text + "\n\n"
                       + "余额信息: " + "\n" + balance + "\n\n"
                       # + "具体余额: " + "\n" + balance + "\n\n"
                       + "套餐使用情况: " + combo_info_text)
        # with open("query_result.txt", "w", encoding="utf-8") as file:  # 将全局变量写入文件
        #     file.write("电话号码: " + "\n" + telephone_text + "\n\n")
        #     file.write("实时话费: " + "\n" + query_result_text + "\n\n")
        #     file.write("余额信息: " + "\n" + balance + "\n\n")
        #     # file.write("具体余额: " + "\n" + balance + "\n\n")
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
    txt_content = input_file.split('\n')

    # 解析文本内容，提取需要的信息
    data = {}
    consumption_amount = None
    subscription_set = []
    skip_line = False
    for idx, line in enumerate(txt_content):
        if "实时话费" in line:
            consumption_amount = float(txt_content[idx + 1].split()[-1])
        elif "套餐及固定费" in line:
            for i in range(idx + 1, len(txt_content),1):
                if "小计" in txt_content[i]:
                    break
                if not is_float(txt_content[i].strip()):
                    x=None
                    for j in range(i + 1, len(txt_content) - 1):
                        if is_float(txt_content[j].strip()) and not is_float(txt_content[j + 1].strip()):
                            x = j
                            break
                    if x is not None:
                        subscription_set.append({
                            "subscription_name": txt_content[i].strip(),
                            "subscription_amount": txt_content[x].strip()
                        })
        elif "合计" in line:
            break

    timestamp = datetime.now().strftime("%Y/%m/%d %H:%M:%S")  # 获取当前时间戳
    # 构造 JSON 数据
    data["balance"] = float(balance)
    data["timeStamp"] = timestamp
    data["consumption_condition"] = {
        "consumption_amount": consumption_amount,
        "consumption_set":
            subscription_set
    }
    print(data)
    # 写入 JSON 文件
    # with open(output_file, 'w', encoding='utf-8') as f:
    #     json.dump(data, f, ensure_ascii=False, indent=2)


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
distance = distance +3

# 3. 移动
cs.crack_slider(distance)

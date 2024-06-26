import sys

def crawl(username, password):
    # 你的爬虫代码
    print(f"Crawling with 用户名: {username} and 密码: {password}")
    # 爬取内容并返回结果
    return "爬取结果"

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python crawler.py <username> <password>")
        sys.exit(1)
    username = sys.argv[1]
    password = sys.argv[2]
    result = crawl(username, password)
    print(result)
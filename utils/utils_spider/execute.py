import subprocess
import sys


targetFile = sys.argv[1]
userID = sys.argv[2]
phoneNum = sys.argv[3]
pwd = sys.argv[4]

while True:
    try:
        # 执行slide.py文件
        # print(userID,targetFile,phoneNum,pwd)
        process = subprocess.Popen(['python',targetFile, userID, phoneNum, pwd])
        process.wait()  # 等待slide.py执行完成

        return_code = process.returncode

        # 如果slide.py正常退出，程序结束
        if return_code == 0:
            break
        else:
            print(f"slide.py returned a non-zero exit code ({return_code}). Restarting...")
    except Exception as e:
        print(f"An error occurred while running slide.py: {e}")
#!/usr/bin/python3
import os
import psycopg2
from dotenv import load_dotenv
from contextlib import closing
import requests
import sys
import time
import datetime


arg_name=sys.argv[1]
dotenv_path = os.path.join(os.path.dirname(__file__), '.env')
dotstatus_path = os.path.join(os.path.dirname(__file__), '.status')
R5m="select * from system.service_exec_log where message_receive_timestamp >= (now() - interval '60 day') and message_receive_timestamp <= (now() - interval '2 hour') and status = 'processing' and exec_result_file_path is null limit 10;" #Request every 5 minutes select * from films;
R1h="select * from system.service_exec_log where message_receive_timestamp >= (now() - interval '60 day') and message_receive_timestamp <= (now() - interval '24 hour') and status = 'processing'  limit 10;" #Request every 1 hour
R1d="select description from system.cron_job_process_log where status = 'error' and stop::date = now()::date limit 10;" #Request every 1 day
pwd = os.path.dirname(os.path.realpath(__file__))


if os.path.exists(dotenv_path):
    load_dotenv(dotenv_path)
    DB_HOST=os.environ.get("DB_HOST")
    DB_NAME=os.environ.get("DB_NAME")
    DB_USER=os.environ.get("DB_USER")
    DB_PASSWORD=os.environ.get("DB_PASSWORD")
    TELEGRAM_TOKEN=os.environ.get("TELEGRAM_TOKEN")
    CHANNEL_ID=os.environ.get("CHANNEL_ID")
else:
    write_log(".env don't exist")
    sys.exit(2)


if os.path.exists(dotstatus_path):
    load_dotenv(dotstatus_path)
    R5M=os.environ.get("R5M")
    R1H=os.environ.get("R1H")
    R1D=os.environ.get("R1D")
else:
    os.mknod(pwd + '/.status')
    with open(pwd + "/.status", 'w') as file:
        file.write("R5M=0\nR1H=0\nR1D=0")
    load_dotenv(dotstatus_path)    
    R5M=os.environ.get("R5M")
    R1H=os.environ.get("R1H")
    R1D=os.environ.get("R1D")

def open_status_problem(var):
    with open(pwd + "/.status", 'rt') as file:
        content = file.read()
        content = content.replace(var+'=0', var + '=1')

    with open(pwd + "/.status", 'wt') as file:
    
        file.write(content)

def close_status_problem(var):
    with open(pwd + "/.status", 'rt') as file:
        content = file.read()
        content = content.replace(var+'=1', var + '=0')

    with open(pwd + "/.status", 'wt') as file:
    
        file.write(content)

def write_log(text):
    now = datetime.datetime.now()
    with open(pwd + "/sql.log", "a") as file:
        file.write(str(now) + text + "\n")

def send_telegram(text: str):
    token = TELEGRAM_TOKEN
    url = "https://api.telegram.org/bot"
    channel_id = CHANNEL_ID
    url += token
    method = url + "/sendMessage"

    r = requests.post(method, data={
         "chat_id": channel_id,
         "text": text
          })

    if r.status_code != 200:
        raise Exception("post_text error")


def check_requests():
    try:
        with closing(psycopg2.connect(dbname=DB_NAME, user=DB_USER, password=DB_PASSWORD, host=DB_HOST)) as conn:
            with conn.cursor() as db:


                if arg_name == "R5m":
                    db.execute(R5m)
                    write_log(" ------------\n\t\t\t   Выполнил запрос")
                    if R5M == "0":
                        write_log(" Статус проблемы закрыт")
                        count_rows = len(db.fetchall())
                        write_log(" Посчитал строки")
                        if count_rows != 0:
                            write_log(" Строки не равны нулю")
                            send_telegram("Alarm!\nОчередь на загрузку первичных данных имеет: " + str(count_rows) + " строки")
                            write_log(" Отправил сообщение")
                            open_status_problem("R5M")
                            write_log(" Статус проблемы открыл")
                            while count_rows != 0:
                                time.sleep(60)
                                write_log(" Проверяю решилось ли")
                                db.execute(R5m)
                                count_rows = len(db.fetchall())
                            close_status_problem("R5M")
                            write_log(" Закрыл статус проблемы")
                            send_telegram("Good!\nОчередь на загрузку первичных данных имеет: " + str(count_rows) + " строк")
                            write_log(" Отправил сообщение что всё хорошо \n\t\t\t   ------------")
                        else:
                            write_log(" ------------")
                    else:
                        sys.exit(2)



                if arg_name == "R1h":
                    db.execute(R1h)
                    write_log(" ------------\n\t\t\t   Выполнил запрос")
                    if R1H == "0":
                        write_log(" Статус проблемы закрыт")
                        count_rows = len(db.fetchall())
                        write_log(" Посчитал строки")
                        if count_rows != 0:
                            write_log(" Строки не равны нулю")
                            send_telegram("Alarm!\nОтветы в СМЭВ о загрузке первичных данных имеет: " + str(count_rows) + " строки")
                            write_log(" Отправил сообщение")
                            open_status_problem("R1H")
                            write_log(" Статус проблемы открыл")
                            while count_rows != 0:
                                time.sleep(60)
                                write_log(" Проверяю решилось ли")
                                db.execute(R1h)
                                count_rows = len(db.fetchall())
                            close_status_problem("R1H")
                            write_log(" Закрыл статус проблемы")
                            send_telegram("Good!\nОтветы в СМЭВ о загрузке первичных данных имеет: " + str(count_rows) + " строки")
                            write_log(" Отправил сообщение что всё хорошо \n\t\t\t   ------------")
                        else:
                            write_log(" ------------")
                    else:
                        sys.exit(2)

                if arg_name == "R1d":
                    db.execute(R1d)
                    write_log(" ------------\n\t\t\t   Выполнил запрос")
                    if R1D == "0":
                        write_log(" Статус проблемы закрыт")
                        count_rows = len(db.fetchall())
                        write_log(" Посчитал строки")
                        if count_rows != 0:
                            write_log(" Строки не равны нулю")
                            send_telegram("Alarm!\nПроблемы с расчетом витрин отчетов данных имеют: " + str(count_rows) + " строки")
                            write_log(" Отправил сообщение")
                            open_status_problem("R1D")
                            write_log(" Статус проблемы открыл")
                            while count_rows != 0:
                                time.sleep(60)
                                write_log(" Проверяю решилось ли")
                                db.execute(R1d)
                                count_rows = len(db.fetchall())
                            close_status_problem("R1D")
                            write_log(" Закрыл статус проблемы")
                            send_telegram("Good!\nПроблемы с расчетом витрин отчетов данных имеют: " + str(count_rows) + " строки")
                            write_log(" Отправил сообщение что всё хорошо \n\t\t\t   ------------")
                        else:
                            write_log(" ------------")
                    else:
                        sys.exit(2)

    except psycopg2.OperationalError as err:
        send_telegram("При подключении к бд произошла ошибка:\n" + str(err))
        sys.exit(1)
    except psycopg2.errors.UndefinedTable as err:
        send_telegram("При запросе к бд произошла ошибка:\n" + str(err))
        sys.exit(1)


if __name__ == "__main__":
    check_requests()

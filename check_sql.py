import os
import psycopg2
from dotenv import load_dotenv
from contextlib import closing
import requests
import sys
import time

arg_name=sys.argv[1]
dotenv_path = os.path.join(os.path.dirname(__file__), '.env')



R5m="select * from system.service_exec_log where message_receive_timestamp >= (now() - interval '60 day') and message_receive_timestamp <= (now() - interval '2 hour') and status = 'processing' and exec_result_file_path is null limit 10;" #Request every 5 minutes
R1h="select * from system.service_exec_log where message_receive_timestamp >= (now() - interval '60 day') and message_receive_timestamp <= (now() - interval '24 hour') and status = 'processing'  limit 10;" #Request every 1 hour
R1d="select description from system.cron_job_process_log where status = 'error' and stop::date = now()::date limit 10;" #Request every 1 day


if os.path.exists(dotenv_path):
    load_dotenv(dotenv_path)
    DB_NAME=os.environ.get("DB_NAME")
    DB_USER=os.environ.get("DB_USER")
    DB_PASSWORD=os.environ.get("DB_PASSWORD")
    TELEGRAM_TOKEN=os.environ.get("TELEGRAM_TOKEN")
    CHANNEL_ID=os.environ.get("CHANNEL_ID")
else:
    print(".env don't exist")
    sys.exit(2)

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
        with closing(psycopg2.connect(dbname=DB_NAME, user=DB_USER, password=DB_PASSWORD, host='localhost')) as conn:
            with conn.cursor() as db:


                if arg_name == "R5m":
                    db.execute("select * from films")
                    count_rows = len(db.fetchall())
                    if count_rows != 0:
                        send_telegram("Alarm!\n" + str(count_rows) + " rows from (now() - interval '2 hour')")
                        while count_rows != 0:
                            time.sleep(10)
                            db.execute("select * from films")
                            count_rows = len(db.fetchall())
                        send_telegram("Good!\n" + str(count_rows) + " rows from (now() - interval '2 hour')")


                if arg_name == "R1h":
                    db.execute(R1h)
                    count_rows = len(db.fetchall())
                    if count_rows != 0:
                        send_telegram("alarm " + str(count_rows) + " rows from (now() - interval '24 hour')")
                        while count_rows != 0:
                            time.sleep(10)
                            db.execute(R1h)
                            count_rows = len(db.fetchall())
                        send_telegram("Good!\n " + str(count_rows) + " rows from (now() - interval '24 hour')")


                if arg_name == "R1d":
                    db.execute(R1d)
                    count_rows = len(db.fetchall())
                    if count_rows != 0:
                        send_telegram("alarm " + str(count_rows) + " rows from system.cron_job_process_log")
                        while count_rows != 0:
                            time.sleep(10)
                            db.execute(R1d)
                            count_rows = len(db.fetchall())
                        send_telegram("Good!\n " + str(count_rows) + " rows from system.cron_job_process_log")


    except psycopg2.OperationalError as err:
        send_telegram("При подключении к бд произошла ошибка:\n" + str(err))
        sys.exit(1)
    except psycopg2.errors.UndefinedTable as err:
        send_telegram("При запросе к бд произошла ошибка:\n" + str(err))
        sys.exit(1)


if __name__ == "__main__":
    check_requests()
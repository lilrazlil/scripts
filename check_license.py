import subprocess
import re
from datetime import datetime
from time import gmtime, strftime

data_now = strftime("%Y-%m-%d", gmtime())

def data_license_kesl():
    word="License expiration date:"
    f = subprocess.check_output(['kesl-control','--app-info']).decode("utf-8").split(sep="\n")
    for line in f:
        if str(line).find(word) != -1:
            return line

def data_license_sn():
    word="UpgrExp="
    f = subprocess.check_output(['cat','/opt/secretnet/etc/serial']).decode("utf-8").split(sep="\n")
    for line in f:
        if str(line).find(word) != -1:
            return line

def days_between(d1, d2):
    d1 = datetime.strptime(d1, "%Y-%m-%d")
    d2 = datetime.strptime(d2, "%Y-%m-%d")
    return (d2 - d1)

def diff_days_kesl():
    data_kesl=re.findall('\d+', str(data_license_kesl()))
    data_kesl= '-'.join(data_kesl)
    day_kesl = str(days_between(data_now,data_kesl)).split(sep=' ')
    return day_kesl[0]

def diff_days_sn():
    data_sn=re.findall('\d+', str(data_license_sn()))
    data_sn= '-'.join(data_sn)
    day_sn = str(days_between(data_now,data_sn)).split(sep=' ')
    return day_sn[0]

if __name__ == "__main__":
    kesl = diff_days_kesl()
    sn = diff_days_sn()
    print("[{\"kesl\" : \"", kesl,"\"},{\"sn\" : \"",sn ,"\"}]")

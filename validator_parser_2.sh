#!/bin/bash
number_block=$(cat global_block); #выбираем с какого блока начать
rpc_address="135.181.5.158";
PIDFILE="/tmp/validator_parser.pid"


if [ -f "$PIDFILE" ];
  then
    if pgrep -F $PIDFILE >/dev/null;
      then
        exit 1
      else echo $$ > $PIDFILE
    fi
  else echo $$ > $PIDFILE
fi


function accept_validators { # функция, чтобы добавить в файл валидаторы, которые приняли блок
for add in $(cat address)
do
sed -i "/${add}/d" ./valid_name;
done
}

function list_pass {
    echo -e "$number_block\t$validator_pass\t$validator_name\t$validator_ip\t$proposer_address" >> list_pass;
}

function write_page_json {
    echo $number_block | xargs -n1 -i curl http://${rpc_address}/rpc/commit?height={} > true.json; # записываем во временный файл json
}

function parser_validator_address {
    jq '.result.signed_header.commit.signatures[] | .["validator_address"]' ./true.json | sed -e 's/^"//' -e 's/"$//' > address; # собираем адреса валидаторов, которые есть в json
}

while true # бесконечный цикл
do
    write_page_json;
    canonical=$(grep canonical ./true.json | awk '{print $2}' | sed -e 's/^"//' -e 's/"$//'); # забираем canonical  
    if [[ "${canonical}" == "true" ]]; then # если он true, то переходим на следующую проверку, если false, то ждем 5 минут, чтобы успели набраться блоки
        parser_validator_address;
        count_address=$(sort address | uniq -c | wc -l); # смотрим их кол-во
        if (( ${count_address} <= 37 )); then # если больше или 5 пропусков, то true
            cp full_valids valid_name; # временный файл
            accept_validators; # функция
            v=$(awk '{print $1}' ./valid_name); # в переменную заносим имена пропустивших валидаторов
            validator_pass=$(echo $v | sed 's/ /,/g'); # убираем пробелы и ставим запятые, чтобы это была одна большая строка
            proposer_address=$(grep "proposer_address" ./true.json | awk '{print $2}' | sed -e 's/^"//' -e 's/"$//'); # забираем proposer_address
            validator_name=$(grep $proposer_address ./full_valids | awk '{print $1}'); # Забираем имя валидатора
            validator_ip=$(grep $proposer_address ./full_valids | awk '{print $2}');# Забираем адрес валидатора
            list_pass; # добавляем в файл указанные переменные в этом порядке
            number_block=$(($number_block+1)); # следующий блок
            echo $number_block > global_block; # добавляем в файл 
        else # если пропусков меньше 5, то переходим на следующий блок 
            number_block=$(($number_block+1));
            echo $number_block > global_block;
        fi 
    else # если canocal false, то ждем 
    sleep 5m; 
    fi
sleep 5; # время перед следующей итерацией
done
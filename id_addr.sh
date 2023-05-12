#!/bin/bash

function show_virsh_list(){
  count_old=$(virsh list --all | wc -l)
  count=$(expr $count_old - 2)
  show_list=$(virsh list --all | tail -n $count | awk '{print $2}')
  echo $show_list
}

function show_ip_addr_virsh(){
  list=$(show_virsh_list)
  for i in ${list[@]}; do
    vm=$i
    for mac in `virsh domiflist ${vm} |grep -o -E "([0-9a-f]{2}:){5}([0-9a-f]{2})"`; do
      ip=$(arp -e | grep $mac  | grep -o -P "^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}")
      echo $vm $ip
    done
  done
}

function main(){
 show_ip_addr_virsh
}

main

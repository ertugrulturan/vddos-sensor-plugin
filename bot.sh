#!/bin/bash
screen -X -S underattackon kill
vddos-switch domain1.com high
vddos-switch domain2.com high
vddos reload
while true
do
netstat -tn | awk '{print $5}' | sed -e 's/:.*//' | grep '\.'| sort | uniq -c | sort -nr | head -24 > logs.txt
cat logs.txt | base64 > antiddos.txt
php /root/bot.php
sleep 1.5
done
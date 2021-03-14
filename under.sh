#!/bin/bash
screen -X -S t13r kill
vddos-switch domain1.com captcha
vddos-switch domain2.com captcha
vddos reload
while true
do
netstat -tn | awk '{print $5}' | sed -e 's/:.*//' | grep '\.'| sort | uniq -c | sort -nr | head -24 > logs.txt
cat logs.txt | base64 > antiddos.txt
php /root/bot2.php
sleep 1.4
done

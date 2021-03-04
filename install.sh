#!/bin/bash
if ! [ $(id -u) = 0 ]; then
   echo "only root run installer"
   exit 1
fi
if [ ! -f /vddos/vddos ]; then
	echo "First VDDOS INSTALL"
fi
yum install net-tools
yum install screen
cd /root/
wget https://raw.githubusercontent.com/ertugrulturan/vddos-sensor-plugin/main/bot.php
wget https://raw.githubusercontent.com/ertugrulturan/vddos-sensor-plugin/main/bot.sh
wget https://raw.githubusercontent.com/ertugrulturan/vddos-sensor-plugin/main/bot2.php
wget https://github.com/ertugrulturan/vddos-sensor-plugin/blob/main/under.sh
clear
echo "Only bot.sh And under.sh Domain Name Edit And Run Script!"
echo 'Start (Run Command) : $ screen -S "t13r" bash /root/bot.sh'

#!/bin/bash
if ! [ $(id -u) = 0 ]; then
   echo "only root run installer"
   exit 1
fi
if [ ! -f /vddos/vddos ]; then
	echo "First VDDOS INSTALL"
fi
wget https://raw.githubusercontent.com/ertugrulturan/vddos-sensor-plugin/main/lwguardian
wget https://github.com/ertugrulturan/vddos-sensor-plugin/releases/download/v1.0/layerweb-vddosensor
mv layerweb-vddosensor /usr/bin/
mv lwguardian /usr/bin/
vddosservicedir=$(echo "W1VuaXRdCkRlc2NyaXB0aW9uPWxheWVyd2ViCkFmdGVyPW5ldHdvcmsudGFyZ2V0CgpbU2VydmljZV0KVHlwZT1zaW1wbGUKRXhlY1N0YXJ0aW9uPS91c3IvYmluL2xheWVyd2ViLXZkZG9zZW5zb3IKRXhlY1N0b3A9a2lsbGFsbCBsYXllcndlYi12ZGRvc2Vuc29yClJlc3RhcnQ9b24tZmFpbHVyZQpSZXN0YXJ0U2VjPTVzCgpbSW5zdGFsbF0KV2FudGVkQnk9bXVsdGktdXNlci50YXJnZXQ=" | base64 --decode)
echo "$vddosservicedir" > /etc/systemd/system/layerwebvddos.service
systemctl daemon-reload
systemctl enable layerwebvddos.service
systemctl start layerwebvddos.service
echo "VDDOS Sensor mode active! - LAYERWEB Systems"


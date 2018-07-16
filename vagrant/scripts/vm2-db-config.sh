#!/bin/bash
#
# 

echo "[*] Updating CentOS via yum update:"
yum update -y
echo "[*] Setting mongodb-org-3.6 to yum.repos.d:" 
sh -c "echo '[mongodb-org-3.6]' > /etc/yum.repos.d/mongodb-org-3.6.repo"
sh -c "echo 'name=MongoDB Repository' >> /etc/yum.repos.d/mongodb-org-3.6.repo"
sh -c "echo 'baseurl=https://repo.mongodb.org/yum/redhat/7/mongodb-org/3.6/x86_64/' >> /etc/yum.repos.d/mongodb-org-3.6.repo"
sh -c "echo 'gpgcheck=1' >> /etc/yum.repos.d/mongodb-org-3.6.repo"
sh -c "echo 'enabled=1' >> /etc/yum.repos.d/mongodb-org-3.6.repo"
sh -c "echo 'gpgkey=https://www.mongodb.org/static/pgp/server-3.6.asc' >> /etc/yum.repos.d/mongodb-org-3.6.repo"
echo "[*] Installing mongodb-org-3.6 via yum:"
yum install -y mongodb-org
echo "[*] Starting mongodb-org-3.6 via systemctl:"
systemctl start mongod
# echo "[*] Exposing mongo-db by edting /etc/mongod.conf (bindIP: 127.0.0.1,192.168.50.5)"

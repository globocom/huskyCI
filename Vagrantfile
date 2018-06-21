# -*- mode: ruby -*-
# vi: set ft=ruby :


# Setting scripts to configure each VM 
$vm2DBConfig = <<-SCRIPT
echo "[*] Updating CentOS via yum update:"
sudo yum update -y
echo "[*] Setting mongodb-org-3.6 to yum.repos.d:" 
sudo sh -c "echo '[mongodb-org-3.6]' > /etc/yum.repos.d/mongodb-org-3.6.repo"
sudo sh -c "echo 'name=MongoDB Repository' >> /etc/yum.repos.d/mongodb-org-3.6.repo"
sudo sh -c "echo 'baseurl=https://repo.mongodb.org/yum/redhat/7/mongodb-org/3.6/x86_64/' >> /etc/yum.repos.d/mongodb-org-3.6.repo"
sudo sh -c "echo 'gpgcheck=1' >> /etc/yum.repos.d/mongodb-org-3.6.repo"
sudo sh -c "echo 'enabled=1' >> /etc/yum.repos.d/mongodb-org-3.6.repo"
sudo sh -c "echo 'gpgkey=https://www.mongodb.org/static/pgp/server-3.6.asc' >> /etc/yum.repos.d/mongodb-org-3.6.repo"
echo "[*] Installing mongodb-org-3.6 via yum:"
sudo yum install -y mongodb-org
echo "[*] Starting mongodb-org-3.6 via systemctl:"
sudo systemctl start mongod
SCRIPT

$vm3DockerConfig = <<-SCRIPT
echo "[*] Updating CentOS via yum update:"
sudo yum update -y
echo "[*] Exposing docker daemon port:" 
sudo mkdir /etc/systemd/system/docker.service.d
sudo sh -c "echo '[Service]' >> /etc/systemd/system/docker.service.d/override.conf"
sudo sh -c "echo 'ExecStart=' >> /etc/systemd/system/docker.service.d/override.conf"
sudo sh -c "echo 'ExecStart=/usr/bin/dockerd -H tcp://0.0.0.0:2376' >> /etc/systemd/system/docker.service.d/override.conf"   
sudo systemctl enable --now docker.service
echo "[*] Restarting docker:" 
sudo systemctl stop docker
sudo systemctl start docker
SCRIPT

$vm4RegisterConfig = <<-SCRIPT
echo "[*] Updating CentOS via yum update:"
sudo yum update -y
# echo "[*] Exposing docker daemon port:" 
# sudo mkdir /etc/systemd/system/docker.service.d
# sudo sh -c "echo '[Service]' >> /etc/systemd/system/docker.service.d/override.conf"
# sudo sh -c "echo 'ExecStart=' >> /etc/systemd/system/docker.service.d/override.conf"
# sudo sh -c "echo 'ExecStart=/usr/bin/dockerd -H tcp://0.0.0.0:2376' >> /etc/systemd/system/docker.service.d/override.conf"   
# sudo systemctl enable --now docker.service
# echo "[*] Restarting docker:" 
# sudo systemctl stop docker
# sudo systemctl start docker
SCRIPT

Vagrant.configure("2") do |config|
  
  config.vm.box = "centos/7"
  
  config.vm.provider "virtualbox" do |v|
    v.memory = 512
    v.cpus = 1
  end

  config.vm.define "vm1-api" do |vm1|
    vm1.vm.network "private_network", ip: "192.168.50.4" 
  end

  config.vm.define "vm2-db" do |vm2|
    vm2.vm.network "private_network", ip: "192.168.50.5"
    vm2.vm.provision :shell, inline: $vm2DBConfig
  end

  config.vm.define "vm3-docker" do |vm3|
    vm3.vm.provision "docker" do |d| end
    vm3.vm.network "private_network", ip: "192.168.50.6"
    vm3.vm.provision :shell, inline: $vm3DockerConfig
  end

  config.vm.define "vm4-register" do |vm4|
    vm4.vm.provision "docker" do |d| end
    vm4.vm.network "private_network", ip: "192.168.50.7"
    vm4.vm.provision :shell, inline: $vm4RegisterConfig
  end


end

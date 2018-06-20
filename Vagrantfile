# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure("2") do |config|
  
  config.vm.box = "centos/7"
  
  config.vm.provider "virtualbox" do |v|
    v.memory = 2048
    v.cpus = 1
  end

  config.vm.provision "docker" do |d|
  end

  config.vm.provision 'shell', inline: <<-EOF
      sudo yum update -y

      echo '[*] Exposing docker daemon port' 
      sudo mkdir /etc/systemd/system/docker.service.d
      sudo sh -c "echo '[Service]' >> /etc/systemd/system/docker.service.d/override.conf"
      sudo sh -c "echo 'ExecStart=' >> /etc/systemd/system/docker.service.d/override.conf"
      sudo sh -c "echo 'ExecStart=/usr/bin/dockerd -H tcp://0.0.0.0:2376' >> /etc/systemd/system/docker.service.d/override.conf"
     
      sudo systemctl enable --now docker.service
      sudo systemctl stop docker
      sudo systemctl start docker
     
    EOF
  
  config.vm.define "vm1-api" do |vm1|
    vm1.vm.network "private_network", ip: "192.168.50.4"
  end

  config.vm.define "vm2-bd" do |vm2|
    vm2.vm.network "private_network", ip: "192.168.50.5"
  end

  config.vm.define "vm3-docker" do |vm3|
    vm2.vm.network "private_network", ip: "192.168.50.6"
  end

  config.vm.define "vm4-register" do |vm4|
    vm2.vm.network "private_network", ip: "192.168.50.7"
  end


end

Vagrant.configure("2") do |config|

  config.vm.synced_folder ".", "/vagrant", type: "rsync", rsync__exclude: [".git/", "vendor"]

  config.vm.box = "centos/7"
  
  config.vm.provider "virtualbox" do |v|
    v.memory = 1024
    v.cpus = 1
  end

  # config.vm.define "vm1-api" do |vm1|
  #   vm1.vm.network "private_network", ip: "192.168.50.4" 
  # end


  config.vm.define "vm2-db" do |vm2|

    vm2.vm.network "private_network", ip: "192.168.50.5"
    vm2.vm.provision :shell, path: "vagrant/scripts/vm2-db-config.sh", privileged: true
  end

  config.vm.define "vm3-docker" do |vm3|
    vm3.vm.provision "docker" do |d| end
    vm3.vm.network "private_network", ip: "192.168.50.6"
    vm3.vm.provision :shell, path: "vagrant/scripts/vm3-docker-config.sh", privileged: true
  end

  # config.vm.define "vm4-register" do |vm4|
  #   vm4.vm.provision "docker" do |d| end
  #   vm4.vm.network "private_network", ip: "192.168.50.7"
  # end

end

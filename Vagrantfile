# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure("2") do |config|
  (1..3).each do |i|
    config.vm.define "sftp_#{i}" do |sftp|
      sftp.vm.box = "ubuntu/bionic64"
      sftp.vm.hostname = "sftp-#{i}"
      sftp.vm.network "private_network", ip: "10.0.0.20#{i}"

      sftp.vm.provider "virtualbox" do |vb|
        vb.name = "sftp-#{i}"
        vb.memory = "1024"
        vb.cpus = 2
      end
      sftp.vm.provision "file", source: "scripts", destination: "/tmp/scripts/"
      sftp.vm.provision "shell", path: "provision/provision.sh"
    
	end
  end
  
  config.vm.define "central_server" do |central|
    central.vm.box = "ubuntu/bionic64"
    central.vm.hostname = "central-server"
    central.vm.network "private_network", ip: "10.0.0.204"

    central.vm.provider "virtualbox" do |vb|
      vb.name = "central-server"
      vb.memory = "1024"
      vb.cpus = 2
    end
	central.vm.provision "file", source: "app", destination: "/home/vagrant/app"
    central.vm.provision "shell", path: "provision/central_server_provision.sh"
  end
end
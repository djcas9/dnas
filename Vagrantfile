# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
Vagrant.require_version '>= 1.5.0'
Vagrant.configure("2") do |config|
  # Every Vagrant virtual environment requires a box to build off of.
  config.vm.box = "DNAS"
  config.vm.box_url = "http://files.vagrantup.com/precise64.box"

  # Fix docker not being able to resolve private registry in VirtualBox
  config.vm.provider :virtualbox do |vb, override|
    vb.customize ["modifyvm", :id, "--natdnshostresolver1", "on"]
    vb.customize ["modifyvm", :id, "--natdnsproxy1", "on"]
  end

  config.vm.provision "docker" do |d|
    d.build_image "/vagrant", args: '-t dnas'
    d.run "sudo dnas -i eth0 -w ~/dnas.txt -d ~/dnas.db -H", demonize: false
    # , args: "-i eth0 -H -w dnas.txt -d dnas.db"
  end

  # plugin conflict
  if Vagrant.has_plugin?("vagrant-vbguest")
    config.vbguest.auto_update = false
  end
end

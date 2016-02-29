# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|

  config.vm.box = "debian/contrib-jessie64"
  config.vm.network :forwarded_port, guest: 80, host: 8080
  config.vm.provider "virtualbox" do |v|
    v.memory = 256
  end

  config.vm.provision "shell", inline: "apt-get update && apt-get install openssh-server sudo zerofree python ca-certificates -y"

  # Install the contents of dist. probably won't run without a db but at least you will know the
  # package is valid.
  config.vm.provision "shell", inline: "dpkg -i /vagrant/dist/*"

end

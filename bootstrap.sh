#/bin/bash

# install apex
curl https://raw.githubusercontent.com/apex/apex/master/install.sh | sh

# install terraform
sudo curl -L https://releases.hashicorp.com/terraform/0.6.16/terraform_0.6.16_linux_amd64.zip -o /usr/local/bin/tf.zip
cd /usr/local/bin && sudo unzip tf.zip
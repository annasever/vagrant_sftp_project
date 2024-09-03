export DEBIAN_FRONTEND=noninteractive

echo "-- Running apt update -- "

apt-get update

echo "-- Installing rkhunter & sshpass -- "
apt-get install -y rkhunter sshpass

echo "-- Creating sftp user & group -- "
groupadd sftp
useradd -m sftp -g sftp
echo "sftp:pass" | chpasswd

echo "-- Configuring sshd -- "
cat <<EOF > /etc/ssh/sshd_config
Subsystem sftp internal-sftp

Match group sftp
ChrootDirectory /home/sftp
X11Forwarding no
AllowTcpForwarding no
RSAAuthentication yes
PasswordAuthentication yes
AuthorizedKeysFile %h/.ssh/authorized_keys
EOF

# "PasswordAuthentication yes" must be blown away once ssh-pub-key-exchange-script.sh is executed

echo "-- Giving ownership of /home/sftp to root -- "
chown root:root /home/sftp

echo "-- Copying bin, lib and lib64 to chroot -- "
sudo cp -r /bin /home/sftp/
sudo cp -r /lib /home/sftp/
sudo cp -r /lib64 /home/sftp/
sudo cp -r /dev /home/sftp/

echo "-- Creating uploads directory for sftp -- "
mkdir -p /home/sftp/uploads
chown sftp:sftp /home/sftp/uploads
chmod 755 /home/sftp # TODO: Maybe move it up top
chmod 755 /home/sftp/uploads

echo "-- Creating .ssh directory for sftp -- "
mkdir -p /home/sftp/.ssh
touch /home/sftp/.ssh/authorized_keys
chown sftp:sftp /home/sftp/.ssh/authorized_keys
chmod 600 /home/sftp/.ssh/authorized_keys
chown sftp:sftp /home/sftp/.ssh
chmod 700 /home/sftp/.ssh

echo "-- Moving scripts to /home/sftp/ --"
mv /tmp/scripts/* /home/sftp/scripts    
chown -R sftp:sftp /home/sftp/scripts
chmod -R 755 /home/sftp/scripts
echo "-- Verifying move --"
ls -l /home/sftp/scripts/

echo "-- Adding cron job for sftp-exchange.sh -- "
sudo -u sftp bash -c 'crontab -l | { cat; echo "*/5 * * * * /home/sftp/scripts/sftp-exchange.sh"; } | crontab -'

sudo -u sftp /home/sftp/scripts/generate-ssh-key.sh


echo "-- Restarting ssh service -- "
systemctl restart ssh
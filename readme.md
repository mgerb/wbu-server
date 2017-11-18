# Where Bar You Server

This is the back end for an app I worked on to learn React Native app development.
It is no longer being maintained.

[Check out a blog post I did here.](https://mitchellgerber.com/post/Projects/2017-11-18-my-first-app)

## flags

> -p production mode

# TODO

- set up a testing framework

## Setup

### SSH
Add SSH public key to authorized_keys

`vim /root/.ssh/authorized_keys`

Enable SSH without password - only SSH key

`vim /etc/ssh/sshd_config`

> `PasswordAuthentication no`

### Redis
`sudo yum update`

`sudo yum install redis`

`sudo systemctl enable redis`

### Enable port 22/443 with Firewalld
`sudo systemctl enable firewalld`

`sudo firewall-cmd --permanent --add-service=ssh`

`sudo firewall-cmd --permanent --add-service=http`

`sudo firewall-cmd --permanent --add-service=https`

`sudo firewall-cmd --permanent --add-port=123/tcp`

### Install Go
Download the latest version of Go from the website

`tar -xvf go1.8.3.linux-amd64.tar.gz`

`mv go /usr/local`

`mkdir /home/go /home/go/bin /home/go/src /home/go/pkg`

#### Edit .bashrc

```
export PATH=$PATH:/usr/local/go/bin
export GOPATH=/home/go
export GOBIN=/home/go/bin
export PATH=$PATH:$GOBIN
```

`source ~/.bashrc`


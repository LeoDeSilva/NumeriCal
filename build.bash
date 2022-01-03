#!/bin/bash

sudo rm /usr/local/bin/cal /usr/local/bin/numerical
go build main.go
sudo cp main /usr/local/bin/numerical
sudo cp /usr/local/bin/numerical /usr/local/bin/cal
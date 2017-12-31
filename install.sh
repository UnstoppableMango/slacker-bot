#!/usr/bin/bash

sudo		/bin/systemctl slacker-bot stop
unzip		/srv/slacker-bot/drop/SlackerBot.zip
mv		-v	/srv/slacker-bot/drop/SlackerBot/* /srv/www/slacker-bot/
sudo		/bin/systemctl slacker-bot start
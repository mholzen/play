# install play-go as a service

sudo ln -sf /home/marc/play-go/conf/play-go.service /etc/systemd/system
sudo ln -sf /home/marc/play-go/conf/play-go.path /etc/systemd/system


# restart service

sudo systemctl daemon-reload
sudo systemctl enable play-go.service play-go.path
sudo systemctl start play-go.path
# Optionally, start your service directly, if it's not set to start with the path unit
sudo systemctl start play-go.service
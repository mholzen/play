# install play-go as a service

sudo ln -sf /home/marc/play-go/conf/play-go.service /etc/systemd/system
sudo ln -sf /home/marc/play-go/conf/play-go-watcher.service /etc/systemd/system
sudo ln -sf /home/marc/play-go/conf/play-go-watcher.path /etc/systemd/system


# restart service

sudo systemctl daemon-reload
sudo systemctl enable play-go.service play-go-watcher.service play-go-watcher.path
sudo systemctl start play-go-watcher.path

# Optionally, start your service directly, if it's not set to start with the path unit
sudo systemctl start play-go.service
# serve the React app

(cd ../play-editor; make push-build)


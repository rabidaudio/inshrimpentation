#!/bin/bash
# /boot/firstrun.sh

# Move binary to final location
mv /boot/inshrimp /usr/local/bin/inshrimp
chmod +x /usr/local/bin/inshrimp

# Create a systemd service file dynamically
cat << 'EOF' > /etc/systemd/system/inshrimp.service
[Unit]
Description=Shrimp monitoring server
After=network.target

[Service]
ExecStart=/usr/local/bin/inshrimp
Restart=always
User=pi

[Install]
WantedBy=multi-user.target
EOF

# TODO: also configure and start the mediamtx service

# Enable and start the service
systemctl daemon-reload
systemctl enable inshrimp.service
systemctl start inshrimp.service

# Clean up or self-disable
rm /boot/firstrun.sh

build:
	go build -o cloudflareddns .

run:
	go build -o cloudflareddns .
	./cloudflareddns

clean:
	go clean
	rm cloudflareddns

install:
	install cloudflare-ddns /usr/local/bin/cloudflareddns
	install -m 664 cloudflareddns.service /etc/systemd/system/cloudflareddns.service
	install -m 664 cloudflareddns.timer /etc/systemd/system/cloudflareddns.timer
	systemctl enable cloudflareddns.service
	systemctl enable cloudflareddns.timer
	systemctl daemon-reload
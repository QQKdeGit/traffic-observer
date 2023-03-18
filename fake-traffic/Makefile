all: nginx fake-traffic

nginx:
	docker run --rm -it -p 8054:80 nginx

fake-traffic:
	bash fake-traffic.sh
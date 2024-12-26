


fmt:
	find . -name '*.go' -exec gofumpt -w -s -extra {} \;

monerod:
	go build -o monerod cmd/monerod/main.go 
	./moneroger \
		--datadir=$(PWD)/data \
		--daemon-port=18081 \

run:
	# Create wallet and data directories
	mkdir -p $(PWD)/data
	# Run on mainnet
	go build -o moneroger cmd/moneroger/main.go 
	./moneroger \
		--datadir=$(PWD)/data \
		--wallet=$(PWD)/data/wallet.keys \
		--daemon-port=18081 \
		--wallet-port=18083

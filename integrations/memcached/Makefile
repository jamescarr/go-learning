file_size := 10
file_name := bigold

generate_file:
	dd if=/dev/urandom of=${file_name}file.dat bs=1048576 count=${file_size}

install:
	glide install

test:
	go test

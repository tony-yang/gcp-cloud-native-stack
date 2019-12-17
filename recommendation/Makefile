PROTODIR=../proto

update-proto:
	python -m grpc_tools.protoc -I$(PROTODIR) --python_out=. --grpc_python_out=. ${PROTODIR}/app.proto

update-dep:
	pip-compile --output-file requirements.txt requirements.in
	pip install -r requirements.txt

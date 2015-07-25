.PHONY: test

test:
	captain build
	docker run --privileged=true \
		-v "${HOME}/.ssh:/root/.ssh:ro" \
		-v "${HOME}/.dockercfg:/root/.dockercfg:ro" \
		-e REPOSITORY=ejholmes/captain-test \
		-e BRANCH=master \
		-e SHA=2e4edf57db00d55051c64d1568e2214858a0897d \
		harbure/captain-builder

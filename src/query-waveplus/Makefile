SRC=airquality.go build.go convert.go main.go must.go serialnum.go type.go util.go waveplus.go yaml.go
EXECNAME=waveplus_mqtt
GOLANG=/usr/bin/go
CURL=/usr/bin/curl
GIT=/usr/bin/git
REPONAME=$(shell basename `git rev-parse --show-toplevel`)
DOCKERREPO=${REPONAME}
VERSION=v1.0-beta
SHA1=$(shell git rev-parse --short HEAD)
NOW=$(shell date +%Y-%m-%d_%T)

${EXECNAME}: ${SRC}
	echo ${REPONAME}
	${GOLANG} build \
		-o ${EXECNAME} \
		-ldflags "-X main.sha1ver=${SHA1} \
		-X main.buildTime=${NOW} \
		-X main.repoName=${REPONAME}"

Docker: pws_exporter
	docker build -t ${DOCKERREPO}:${VERSION} .

update-go:
	${GOLANG} get tinygo.org/x/bluetooth
	${GOLANG} get gopkg.in/yaml.v3
	${GOLANG} mod tidy

clean:
	rm -f ${EXECNAME}



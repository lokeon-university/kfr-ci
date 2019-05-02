BOTIMAGE=registry.heroku.com/kfr-cibot/web
SERVERIMAGE=gcr.io/kfr-ci/kfr-server
BOTENV=`cat tg-bot/.inlinenv`
SERVERENV=`cat server/.inlinenv`

bot:
	docker build -t ${BOTIMAGE} -f tg-bot/Dockerfile .

server:
	docker build -t ${SERVERIMAGE} -f server/Dockerfile .

push: bot server
	docker push ${BOTIMAGE}
	docker push ${SERVERIMAGE}

release-bot:
	heroku config:set -a kfr-cibot ${BOTENV}
	heroku container:release web -a kfr-cibot

release-server:
	gcloud beta run deploy --image ${SERVERIMAGE} --update-env-vars ${SERVERENV}

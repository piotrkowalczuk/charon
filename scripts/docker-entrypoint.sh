#!/bin/sh
set -e

: ${CHAROND_PORT:=8080}
: ${CHAROND_HOST:=0.0.0.0}
: ${MNEMOSYNED_LOG_ENVIRONMENT:=production}
: ${MNEMOSYNED_LOG_LEVEL:=info}
: ${CHAROND_STORAGE:=postgres}
: ${CHAROND_MONITORING:=false}
: ${CHAROND_PASSWORD_BCRYPT_COST:=10}
: ${CHAROND_MNEMOSYNED_ADDRESS:=mnemosyned:8080}
: ${CHAROND_MNEMOSYNED_TLS_ENABLED:=false}
: ${CHAROND_POSTGRES_ADDRESS:=postgres://postgres:postgres@postgres/postgres?sslmode=disable}
: ${CHAROND_POSTGRES_DEBUG:=false}
: ${CHAROND_TLS_ENABLED:=false}

if [ "$1" = 'charond' ]; then
	exec charond \
		-host=${CHAROND_HOST} \
		-port=${CHAROND_PORT} \
		-log.environment=${CHAROND_LOG_ENVIRONMENT} \
		-log.level=${CHAROND_LOG_LEVEL} \
		-mnemosyned.address=${CHAROND_MNEMOSYNED_ADDRESS} \
		-mnemosyned.tls=${CHAROND_MNEMOSYNED_TLS_ENABLED} \
		-mnemosyned.tls.crt=${CHAROND_MNEMOSYNED_TLS_CRT} \
		-password.strategy=${CHAROND_PASSWORD_STRATEGY} \
		-password.bcryptcost=${CHAROND_PASSWORD_BCRYPT_COST} \
		-monitoring=${CHAROND_MONITORING} \
		-postgres.address=${CHAROND_POSTGRES_ADDRESS} \
		-postgres.debug=${CHAROND_POSTGRES_DEBUG} \
		-tls=${CHAROND_TLS_ENABLED} \
		-tls.crt=${CHAROND_TLS_CRT} \
		-tls.key=${CHAROND_TLS_KEY}
fi

: ${CHARONCTL_CHAROND_HOST:=charond}
: ${CHARONCTL_AUTH_ENABLED:=true}
: ${CHARONCTL_REGISTER_SUPERUSER:=false}
: ${CHARONCTL_REGISTER_CONFIRMED:=false}
: ${CHARONCTL_REGISTER_STAFF:=false}
: ${CHARONCTL_REGISTER_ACTIVE:=false}
: ${CHARONCTL_REGISTER_IF_NOT_EXISTS:=false}
: ${CHARONCTL_REGISTER_PERMISSIONS:=""}
: ${CHARONCTL_FIXTURES_PATH:="/data/fixtures.json"}

if [ "$1" = 'charonctl' ]; then
	if [ "$2" = 'register' ]; then
		eval charonctl register \
			-address='${CHARONCTL_CHAROND_HOST}:${CHAROND_PORT}' \
			-auth=${CHARONCTL_AUTH_ENABLED} \
			-auth.username=\"${CHARONCTL_AUTH_USERNAME}\" \
			-auth.password=\"${CHARONCTL_AUTH_PASSWORD}\" \
			-register.ifnotexists=${CHARONCTL_REGISTER_IF_NOT_EXISTS} \
			-register.username=\"${CHARONCTL_REGISTER_USERNAME}\" \
			-register.password=\"${CHARONCTL_REGISTER_PASSWORD}\" \
			-register.firstname=\"${CHARONCTL_REGISTER_FIRSTNAME}\" \
			-register.lastname=\"${CHARONCTL_REGISTER_LASTNAME}\" \
			-register.superuser=${CHARONCTL_REGISTER_SUPERUSER} \
			-register.confirmed=${CHARONCTL_REGISTER_CONFIRMED} \
			-register.staff=${CHARONCTL_REGISTER_STAFF} \
			-register.active=${CHARONCTL_REGISTER_ACTIVE} \
			-register.permission=\"${CHARONCTL_REGISTER_PERMISSIONS}\"
		exit $?
	fi
	if [ "$2" = 'load' ]; then
		eval charonctl load \
			-address='${CHARONCTL_CHAROND_HOST}:${CHAROND_PORT}' \
			-auth=${CHARONCTL_AUTH_ENABLED} \
			-auth.username=\"${CHARONCTL_AUTH_USERNAME}\" \
			-auth.password=\"${CHARONCTL_AUTH_PASSWORD}\" \
			-fixtures.path=\"${CHARONCTL_FIXTURES_PATH}\"
		exit $?
	fi
fi

exec "$@"

#!/bin/sh
set -e

: ${CHAROND_PORT:=8080}
: ${CHAROND_HOST:=0.0.0.0}
: ${CHAROND_LOG_ADAPTER:=stdout}
: ${CHAROND_LOG_FORMAT:=json}
: ${CHAROND_LOG_LEVEL:=6}
: ${CHAROND_STORAGE:=postgres}
: ${CHAROND_MONITORING:=false}
: ${CHAROND_PASSWORD_BCRYPT_COST:=10}
: ${CHAROND_MNEMOSYNED_ADDRESS:=mnemosyned:8080}
: ${CHAROND_POSTGRES_ADDRESS:=postgres://postgres:postgres@postgres/postgres?sslmode=disable}
: ${CHAROND_POSTGRES_DEBUG:=false}
: ${CHAROND_TLS_ENABLED:=false}
: ${CHAROND_LDAP_ENABLED:=false}

if [ "$1" = 'charond' ]; then
exec charond \
	-host=${CHAROND_HOST} \
	-port=${CHAROND_PORT} \
	-log.adapter=${CHAROND_LOG_ADAPTER} \
	-log.format=${CHAROND_LOG_FORMAT} \
	-log.level=${CHAROND_LOG_LEVEL} \
	-mnemosyned.address=${CHAROND_MNEMOSYNED_ADDRESS} \
	-password.strategy=$CHAROND_PASSWORD_STRATEGY \
	-password.bcryptcost=$CHAROND_PASSWORD_BCRYPT_COST \
	-monitoring=$CHAROND_MONITORING \
	-postgres.address=$CHAROND_POSTGRES_ADDRESS \
	-postgres.debug=$CHAROND_POSTGRES_DEBUG \
	-tls=$CHAROND_TLS_ENABLED \
	-tls.certfile=$CHAROND_TLS_CERT_FILE \
	-tls.keyfile=$CHAROND_TLS_KEY_FILE \
	-ldap=$CHAROND_LDAP_ENABLED \
	-ldap.address=$CHAROND_LDAP_ADDRESS \
	-ldap.dn=$CHAROND_LDAP_DN \
	-ldap.password=$CHAROND_LDAP_PASSWORD
fi

: ${CHARONCTL_CHAROND_HOST:=charond}
: ${CHARONCTL_AUTH_ENABLED:=true}
: ${CHARONCTL_REGISTER_SUPERUSER:=false}
: ${CHARONCTL_REGISTER_CONFIRMED:=false}
: ${CHARONCTL_REGISTER_STAFF:=false}
: ${CHARONCTL_REGISTER_ACTIVE:=false}

if [ "$1" = 'charonctl' ]; then
	if [ "$2" = 'register' ]; then
		IFS=$'\n'
		permissions=""
		for i in $(echo $CHARONCTL_REGISTER_PERMISSIONS | tr "," "\n")
		do
		  permissions=$permissions" -register.permission='$i' "
		done
		eval charonctl register \
			-address="${CHARONCTL_CHAROND_HOST}:${CHAROND_PORT}" \
			-auth=${CHARONCTL_AUTH_ENABLED} \
			-auth.username="${CHARONCTL_AUTH_USERNAME}" \
			-auth.password="${CHARONCTL_AUTH_PASSWORD}" \
			-register.username="${CHARONCTL_REGISTER_USERNAME}" \
			-register.password="${CHARONCTL_REGISTER_PASSWORD}" \
			-register.firstname="${CHARONCTL_REGISTER_FIRSTNAME}" \
			-register.lastname="${CHARONCTL_REGISTER_LASTNAME}" \
			-register.superuser=${CHARONCTL_REGISTER_SUPERUSER} \
			-register.confirmed=${CHARONCTL_REGISTER_CONFIRMED} \
			-register.staff=${CHARONCTL_REGISTER_STAFF} \
			-register.active=${CHARONCTL_REGISTER_ACTIVE} \
			${permissions} # if not last then it can break rest of the script
		exit $?
	fi
fi

exec "$@"
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
	-host=$CHAROND_HOST \
	-port=$CHAROND_PORT \
	-log.adapter=$CHAROND_LOG_ADAPTER \
	-log.format=$CHAROND_LOG_FORMAT \
	-log.level=$CHAROND_LOG_LEVEL \
	-mnemosyned.address=$CHAROND_MNEMOSYNED_ADDRESS \
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

: ${CHARONCTL_AUTH_DISABLED:=false}

if [ "$1" = 'charonctl register' ]; then
exec charonctl register \
	-address="$CHAROND_HOST:$CHAROND_PORT" \
	-auth.disabled=$CHARONCTL_AUTH_DISABLED \
	-auth.username=$CHARONCTL_AUTH_USERNAME \
	-auth.password=$CHARONCTL_AUTH_PASSWORD \
	-register.username=$CHARONCTL_REGISTER_USERNAME \
	-register.username=$CHARONCTL_REGISTER_PASSWORD \
	-register.username=$CHARONCTL_REGISTER_FIRSTNAME \
	-register.username=$CHARONCTL_REGISTER_LASTNAME \
	-register.username=$CHARONCTL_REGISTER_PERMISSIONS \
	-register.username=$CHARONCTL_REGISTER_SUPERUSER \
	-register.username=$CHARONCTL_REGISTER_CONFIRMED \
	-register.username=$CHARONCTL_REGISTER_STAFF \
	-register.username=$CHARONCTL_REGISTER_ACTIVE
fi

exec "$@"
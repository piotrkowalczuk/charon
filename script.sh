#!/usr/bin/env bash

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
			-auth.username="admin@travelaudience.com" \
			-auth.password="admin" \
			-register.username="${CHARONCTL_REGISTER_USERNAME}" \
			-register.password="${CHARONCTL_REGISTER_PASSWORD}" \
			-register.firstname=${CHARONCTL_REGISTER_FIRSTNAME} \
			-register.lastname=${CHARONCTL_REGISTER_LASTNAME} \
			-register.superuser=${CHARONCTL_REGISTER_SUPERUSER} \
			-register.confirmed=${CHARONCTL_REGISTER_CONFIRMED} \
			-register.staff=${CHARONCTL_REGISTER_STAFF} \
			-register.active=${CHARONCTL_REGISTER_ACTIVE} \
			$permissions # if not last then it can break rest of the script
	fi
fi
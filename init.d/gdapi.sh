#! /bin/sh

### BEGIN INIT INFO
# Provides:		gdapi
# Required-Start:	$remote_fs $syslog
# Required-Stop:	$remote_fs $syslog
# Default-Start:	2 3 4 5
# Default-Stop:		
# Short-Description:	GDO Homemade garage door opener
### END INIT INFO

set -e

# /etc/init.d/gdapi: start and stop the GDO Homemade garage door opener daemon

GDAPI_PATH=/home/djerman/sbin
GDAPI_CMD=gdapi

test -x $GDAPI_PATH || exit 0
( $GDAPI_PATH/$GDAPI_CMD -\? 2>&1 | grep -q gdapi ) 2>/dev/null || exit 0

umask 022

### no config
# if test -f /etc/default/ssh; then
#     . /etc/default/ssh
# fi
#

. /lib/lsb/init-functions
#
if [ -n "$2" ]; then
    GDAPI_OPTS="$GDAPI_OPTS $2"
fi

# Are we running from init?
run_by_init() {
    ([ "$previous" ] && [ "$runlevel" ]) || [ "$runlevel" = S ]
}

#check_for_no_start() {
    # forget it if we're trying to start, and /etc/ssh/sshd_not_to_be_run exists
    # if [ -e /etc/ssh/sshd_not_to_be_run ]; then 
	# if [ "$1" = log_end_msg ]; then
	#     log_end_msg 0 || true
	# fi
	# if ! run_by_init; then
	#     log_action_msg "OpenBSD Secure Shell server not in use (/etc/ssh/sshd_not_to_be_run)" || true
	# fi
	# exit 0
    # fi
#}

check_dev_null() {
    if [ ! -c /dev/null ]; then
	if [ "$1" = log_end_msg ]; then
	    log_end_msg 1 || true
	fi
	if ! run_by_init; then
	    log_action_msg "/dev/null is not a character device!" || true
	fi
	exit 1
    fi
}

#check_privsep_dir() {
    # Create the PrivSep empty dir if necessary
    # if [ ! -d /run/sshd ]; then
	# mkdir /run/sshd
	# chmod 0755 /run/sshd
    # fi
#}

#check_config() {
    # if [ ! -e /etc/ssh/sshd_not_to_be_run ]; then
	# /usr/sbin/sshd $SSHD_OPTS -t || exit 1
    # fi
#}

export PATH="${PATH:+$PATH:}$GDAPI_PATH"

case "$1" in
  start)
#	check_privsep_dir
#	check_for_no_start
#	check_dev_nulls
	log_daemon_msg "Starting gdapi server" "$GDAPI_CMD" || true
	if start-stop-daemon --start --quiet --oknodo --chuid 0:0 --pidfile /run/sshd.pid --exec /usr/sbin/sshd -- $SSHD_OPTS; then
	    log_end_msg 0 || true
	else
	    log_end_msg 1 || true
	fi
	;;
  stop)
	log_daemon_msg "Stopping gdapi server" "$GDAPI_CMD" || true
	if start-stop-daemon --stop --quiet --oknodo --pidfile /run/$GDAPI_CMD.pid --exec $GDAPI_PATH/$GDAPI_CMD; then
	    log_end_msg 0 || true
	else
	    log_end_msg 1 || true
	fi
	;;

  reload|force-reload)
#	check_for_no_start
#	check_config
	log_daemon_msg "Reloading GDAPI server's configuration" "sshd" || true
	if start-stop-daemon --stop --signal 1 --quiet --oknodo --pidfile /run/$GDAPI_CMD.pid --exec $GDAPI_PATH/$GDAPI_CMD; then
	    log_end_msg 0 || true
	else
	    log_end_msg 1 || true
	fi
	;;

  restart)
	check_privsep_dir
	check_config
	log_daemon_msg "Restarting OpenBSD Secure Shell server" "sshd" || true
	start-stop-daemon --stop --quiet --oknodo --retry 30 --pidfile /run/$GDAPI_CMD.pid --exec $GDAPI_PATH/$GDAPI_CMD;
	check_for_no_start log_end_msg
	check_dev_null log_end_msg
	if start-stop-daemon --start --quiet --oknodo --chuid 0:0 --pidfile /run/$GDAPI_CMD.pid --exec $GDAPI_PATH/$GDAPI_CMD; -- $GDAPI_OPTS; then
	    log_end_msg 0 || true
	else
	    log_end_msg 1 || true
	fi
	;;

  try-restart)
	check_privsep_dir
	check_config
	log_daemon_msg "Restarting OpenBSD Secure Shell server" "sshd" || true
	RET=0
	start-stop-daemon --stop --quiet --retry 30 --pidfile /run/$GDAPI_CMD.pid --exec $GDAPI_PATH/$GDAPI_CMD; || RET="$?"
	case $RET in
	    0)
		# old daemon stopped
		check_for_no_start log_end_msg
		check_dev_null log_end_msg
		if start-stop-daemon --start --quiet --oknodo --chuid 0:0 --pidfile /run/$GDAPI_CMD.pid --exec $GDAPI_PATH/$GDAPI_CMD; -- $GDAPI_OPTS; then
		    log_end_msg 0 || true
		else
		    log_end_msg 1 || true
		fi
		;;
	    1)
		# daemon not running
		log_progress_msg "(not running)" || true
		log_end_msg 0 || true
		;;
	    *)
		# failed to stop
		log_progress_msg "(failed to stop)" || true
		log_end_msg 1 || true
		;;
	esac
	;;

  status)
	status_of_proc -p /run/$GDAPI_CMD.pid $GDAPI_PATH/$GDAPI_CMD $GDAPI_CMD && exit 0 || exit $?
	;;

  *)
	log_action_msg "Usage: /etc/init.d/$GDAPI_CMD.sh {start|stop|reload|force-reload|restart|try-restart|status}" || true
	exit 1
esac

exit 0

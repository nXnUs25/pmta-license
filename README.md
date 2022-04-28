
# pmta-license
Simple check for PowerMTA license expiration check

## Usage

Check supports email and OS signals which mean it can be pluged into other monitoring systems such Incing / Sendu / Nagios or it can by executed via cronjobs. 
It can be also run just as a command line tool to quickly verify the license key hosted on server.

If check is integrated with some other monitoring tools i.e. Sensu it mean alerts will be handled by that systems base on returned exit code

0 - [OK]
2 - [CRITICAL]

Executing the check from command line with flag `-q|-quiet` will not trigger any alerts, It will just log all to the console. 

```shell
❯ ./check-pmta-license -h
2022/04/28 10:25:02 check-pmta-license: [ERROR]  mkdir /var/log/check-pmta-license: permission denied


Usage: 
        check-pmta-license [-h|-help] [-e|-email addresses] [-t|-threshold num] [-r|-remainder num] [-f|-file pathToFile] [-s|-server smtp host] [-p|-port port] [-m|-message text] [-q|-quiet] [-v|-verbose]

Description:

  -c    when flag is used, script can be plugged into sensu / icinga monitoring tool. Handles two types of status: [0 - OK] [2 - ERROR]
  -chkmode
        when flag is used, script can be plugged into sensu / icinga monitoring tool. Handles two types of status: [0 - OK] [2 - ERROR]
  -e value
        email addresses to which we are going to send notifications. Comma separated values.
  -email value
        email addresses to which we are going to send notifications. Comma separated values.
  -f string
        path to the license file (default "/opt/pmta/license")
  -file string
        path to the license file (default "/opt/pmta/license")
  -h    display the usage info
  -help
        display the usage info
  -m string
        Email Subject (default "PMTA License is going to expire in few days")
  -message string
        Email Subject (default "PMTA License is going to expire in few days")
  -p int
        port number of smtp server (default 25)
  -port int
        port number of smtp server (default 25)
  -q    when flag is used, script will only print messages to stdout/log. No emails will be send.
  -quiet
        when flag is used, script will only print messages to stdout/log. No emails will be send.
  -r int
        the number of days before expired date since first alert (default 2)
  -remainder int
        the number of days before expiration date since first alert (default 2)
  -s string
        SMTP server name (hostname) (default "smtp-server.something.com")
  -server string
        SMTP server name (hostname) (default "smtp-server.something.com"")
  -t int
        the number of days before expired date (default 5)
  -threshold int
        the number of days before expiration date (default 5)
  -v    when flag is used, script will print debug messages to stdout/log file.
  -verbose
        when flag is used, script will print debug messages to stdout/log file.

The script will check the license for PMTA Server, and if it is below the threshold
Mail will be send for the notification
There are going to be 2 alerts before the expire date, and continuous alerts after expire date.
        - 1st alert will be triggered at X - [5] days equal threshold
        - 2nd alert will be triggered at X - Y (5 - 2 = 3) days before expire date

**Note: 
        X is the threshold number
        Y is the remainder number
        where Y < X
```

Example command to execute chech and send email alerts.
`./check-pmta-license -e nonus25@other.com,chmm@gmai.com,some.address@gmal.com,some.addr@gmai.com -t 44 -r 20 -f ./license`
Above example will worn 



### Build

To build the binares simply just execute `make` command.

```shell
❯ make
go fmt
go build -o check-pmta-license *.go
ls -l check-pmta-license
-rwxrwxr-x  1 nonus25  staff  3655568 28 Apr 10:20 check-pmta-license
```

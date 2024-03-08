# Neorgify:

A go server for reading in your markdown and neorg files and finding tasks that
are incomplete. These tasks are logged and remembered by neorgify so that it can
remind you via any ntfy.sh server.

## Assumptions:
1. Your ntfy.sh server has authentication, unknown if this works with servers
   without.
1. Tasks you want notifications for are in files that end in .md, .neorg, .org,
   or .txt.
1. Incomplete tasks begin with some number of spaces or dashes then a "( )" or
   "[ ]".
1. You only want to be reminded of incomplete tasks (some formats have other
   todo options like repeat).
1. You have a computer that can run this at all times or just when you want
   notifications.

## Configuration:
Neorgify searches for files within the folder that it is running in for a couple of pieces of configuration.
Files must contain just one line.

Files:
- server: https://ntfy.yourserver.com/YourNtfy.shTopic
- login: base64 hash of \<ntfy-username\>:\<password\>

you can create the login file with this bash script:
```bash
echo -n '<ntfy-username>:<password>' | base64 > login
```
make sure you replace the values with your login credentials or this will be disappointing.

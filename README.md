# Neorgify

A go server for reading in your markdown, org, and norg files and finding
tasks that are incomplete. Neorgify logs and remembers these tasks so that it
can remind you via any ntfy.sh server.

## Assumptions

1. Your ntfy.sh server has authentication, unknown if this works with servers
   without.
1. Tasks you want notifications for are in files that end in .md, .norg, .org,
   or .txt.
1. Incomplete tasks begin with some number of spaces or dashes then a "( )" or
   "[ ]".
1. You want to know of incomplete tasks. Some formats have other todo options
   like repeat.
1. You have a computer that can run this at all times or when you want
   notifications.
1. You don't already use "notifications" or "timers" as a topic in your <ntfy.sh>

## Configuration

Neorgify searches for files within the folder that it's running in for a couple
of pieces of configuration.

**Files must contain one line.**

Files:

- server: https://\<your-url>
- login: base-64 hash of \<ntfy-username\>:\<password\>
- folder: full path to your folder, example: /home/username/Notes

You can create the login file with this bash script:

```bash
echo -n '<ntfy-username>:<password>' | base64 > login
```

**Make sure you replace the values with your login credentials or this will be disappointing.**

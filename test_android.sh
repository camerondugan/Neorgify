#!/usr/bin/env bash

gogio -target android neorgify
adb uninstall localhost.neorgify
adb install neorgify.apk

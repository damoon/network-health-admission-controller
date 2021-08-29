#!/bin/sh
echo `which network-health-admission-controller` | entr -nr `which network-health-admission-controller` $@

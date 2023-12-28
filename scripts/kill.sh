#!/bin/bash

netstat -vancp tcp | grep 22123 | head -1 | awk '{ print $9 }' | xargs kill -9

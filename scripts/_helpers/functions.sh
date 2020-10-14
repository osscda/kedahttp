#!/bin/bash

function shoutln {
  echo
  printf '%s\n' "$1" | awk '{ print toupper($0) }'
  echo
}

function shout {
  printf '%s' "$1" | awk '{ print toupper($0) }'
}

function log {
  echo -n "-> $1"
}

function logln () {
  echo "-> $1"
}

commandExists () {
    type "$1" &> /dev/null ;
}

function pause () {
  read -s -n 1 -p "$*"
}

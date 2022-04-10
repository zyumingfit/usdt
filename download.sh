#!/bin/bash

while [[ -n "$1" ]] ;do
    case "$1" in
        -date)
            Date=$2
            shift 2
            ;;
        *)
            exit 1
            ;;
    esac
done

Key=TUhrAdFsYTAzfiKC53VLWjpD5T8gEgtxgq
Ip=34.222.151.32
Port=5000

wget http://${Ip}:${Port}/report/${Date}.txt --header keypass:${Key}

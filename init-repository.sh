#!/bin/bash

find ~/Desktop/ -type d -name .git | xargs rm -rf
git init
git add .
git commit -m "first commit"
git branch -M main
git remote add origin git@github.com:xizhanbo/gin.git
git push -u origin main
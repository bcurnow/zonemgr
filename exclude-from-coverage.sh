#!/bin/bash

while read p || [ -n "$p" ]
do
sed -i "/${p//\//\\/}/d" ./coverage.out 
done < ./exclude-from-coverage.txt
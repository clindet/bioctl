#!/bin/bash

mkdir a
mkdir b

zip -q -r -o myfile.zip a
tar -czvf myfile.tar.gz b
rmdir a
rmdir b

bioctl -u 'myfile.zip myfile.tar.gz'

ls -l

rmdir a
rmdir b
rm myfile.zip myfile.tar.gz


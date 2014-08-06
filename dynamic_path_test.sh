#!/bin/bash
this_file=$(basename $0)
base_dir_len=`expr ${#this_file} - 3`
echo $base_dir_len
test_dir=${this_file: 0: `expr $base_dir_len`}
echo $test_dir
echo "cd to $test_dir"
cd $test_dir
./$this_file
cd -

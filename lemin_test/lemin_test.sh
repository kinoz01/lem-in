#!/bin/bash

# Check if a file argument is provided
if [ $# -ne 1 ]; then
  echo "Usage: $0 <input_file>"
  exit 1
fi

# Store the input file argument
input_file=$1

# Run leminTest and lemin, saving their outputs to respective files
./leminTest "$input_file" > leminTest.txt
./lemin "$input_file" > lemi.txt

grep '^L' lemi.txt > lemin.txt

# Count number of lines and words for both files
leminTest_lines=$(wc -l < leminTest.txt)
leminTest_words=$(wc -w < leminTest.txt)
lemin_lines=$(wc -l < lemin.txt)
lemin_words=$(wc -w < lemin.txt)

# Print the results
echo "leminTest: $leminTest_lines frames, $leminTest_words moves"
echo "lemin    : $lemin_lines frames, $lemin_words moves"


if [ "$leminTest_lines" -eq "$lemin_lines" ] && [ "$leminTest_words" -eq "$lemin_words" ]; then
  echo -e "\e[32mOK\e[0m"  
else
  echo -e "\e[31mKO\e[0m"  
fi

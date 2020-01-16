#!/bin/bash
mkdir test_par
cd test_par

echo 'touch /tmp/$1 /tmp/$2; sleep ${1}' > job.sh && bioctl par --cmd "sh job.sh" -t 2 --index 1,2,5-10
bioctl par --cmd 'sh job.sh {{index}} {{key2}}' -t 4 --index 1,2,5-10 --env "key2:123"
bioctl par --cmd 'echo {{index}} {{key2}}; sleep {{index}}' -t 4 --index 1,2,5-10 --env "key2:123" --verbose 2 --save-log
bioctl par --cmd 'echo {{key2}}; sleep {{index}}' -t 4 --index 1,2,5-10 --env "key2:123" --force-idx false --save-log
echo 'sh job.sh' | bioctl par -t 4 --index 1,2,5-10 -

cd ..
rm -rf test_par

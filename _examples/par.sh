echo 'echo $1 $2; sleep ${1}' > job.sh && bioctl par --cmd "sh job.sh" -t 2 --index 1,2,5-10
bioctl par --cmd 'sh job.sh {{index}} {{key2}}' -t 4 --index 1,2,5-10 --env "key2:123"
bioctl par --cmd 'echo {{index}} {{key2}}; sleep {{index}}' -t 4 --index 1,2,5-10 --env "key2:123"
bioctl par --cmd 'echo {{key2}}; sleep {{index}}' -t 4 --index 1,2,5-10 --env "key2:123" --force-idx false
echo 'sh job.sh' | bioctl par -t 4 --index 1,2,5-10 -
rm job.sh
rm -rf _log

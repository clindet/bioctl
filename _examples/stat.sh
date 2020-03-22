bioctl stat --min 1 2 3 4 5 100
bioctl stat --max 1 2 3 4 5
bioctl stat --mean 1 3 5 7 9 26 100
bioctl stat --median 1 10 3 4 100 143 123 12 22.2
bioctl stat --sum 1 1 2 3 2 14234 12 12 1331 23 12 12

# returns the most frequent value
bioctl stat --mfreq 2 10 2 2 10 10 11 12 14

# returns the amount of variation in the dataset.
bioctl stat --var 2 10 2 2 10

# returns the pearson product-moment correlation coefficient between two group variables
# vector1: 1-6 vector2: 10-51
bioctl stat --pearson 1 2 3 4 5 6 10 20 30 40 50 51

# returns the relative standing in a slice of floats.
bioctl stat --percentile 30 1 2 3 4 5 6 7 8 9 10

bioctl stat --freq 1 2 3 2 1 1 1 1 1 1 1 1 123 124 12 | sort

bioctl statdf --pearson _examples/test.csv --at ",0:1" --print

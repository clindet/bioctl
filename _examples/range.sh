#!/bin/bash
# range number sequence
bioctl range 1 100 2.5
bioctl range --start 2 --end 1000 --step 15

# char mode
bioctl range --mode char --step 5
bioctl range --mode char --step 3 --start-char s --end-char z --step 2 --ref-str qrstuvwxy12442z
bioctl range --mode char --step 3 --start-char s --end-char z --step 2 --ref-str qrstuvwxy12442z --sep ''

#!/bin/bash
# convert Pubmed XML to clean JSON string
bget api ncbi -q "Galectins control MTOR and AMPK in response to lysosomal damage to induce autophagy OR MTOR-independent autophagy induced by interrupted endoplasmic reticulum-mitochondrial Ca2+ communication: a dead end in cancer cells. OR The PARK10 gene USP24 is a negative regulator of autophagy and ULK1 protein stability OR Coordinate regulation of autophagy and the ubiquitin proteasome system by MTOR." | bioctl cvrt --xml2json pubmed

# convert SRA XML to clean JSON string
bget api ncbi -d 'sra' -q PRJNA527715 | bioctl cvrt --xml2json sra -

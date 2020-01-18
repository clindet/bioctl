#!/bin/bash
# json pretty
echo '{"a": {"a": 123}, "b": {"b": 567}}' | bioctl fmt --json-pretty --indent 2 -
# key:value => slice
echo '{"a": {"a": 123}, "b": {"b": 567}}' | bioctl fmt --json-to-slice --indent 4 -


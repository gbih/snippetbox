# horizontal: enscript -1rG
# https://www.ibm.com/support/knowledgecenter/ssw_aix_72/e_commands/enscript.html
# */*.tmpl \
# */**/*.tmpl \
# */**/**/*.tmpl \
# enscript --tabsize=4 -1G --line-numbers -p - --word-wrap  --color=0 \
enscript --tabsize=4 -1G  -p - --word-wrap  --color=0 \
*.go \
*/*.go \
*/**/*.go \
*/**/**/*.go \
| pstopdf -i -o SOURCE.pdf && open SOURCE.pdf
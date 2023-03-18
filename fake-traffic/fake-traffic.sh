METHODS=(GET POST PUT DELETE)

while true
do
    curl 172.20.62.49:8054      \
        --proxy localhost:8079  \
        > /dev/null 2>&1        \
        -X ${METHODS[$RANDOM % ${#METHODS[@]} ]}

    sleep $( awk -v min=0 -v max=1 'BEGIN{srand(); print min+rand()*(max-min)}')
done
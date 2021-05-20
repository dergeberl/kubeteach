#!/bin/bash
TASKS=$(kubectl get task --no-headers | cut -d " " -f1)
for t in $TASKS
do
    kubectl apply -f exercises/set1/solutions/$t.yaml
    TIMEOUT=15
    i=0
    while [ $i -lt $TIMEOUT ]; do 
        i=$[$i+1];
        kubectl get task $t -o jsonpath='{.status.state}' | grep successful
        if [ $? -eq 0 ]; then
            echo "$t successful"
            i=$[$i+$TIMEOUT];
        else
          sleep 1
        fi
        if [ $i -eq $TIMEOUT ]; then
           echo "$t did not get successful in time"
           exit 1
        fi
    done
done

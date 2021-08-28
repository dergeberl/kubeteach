#!/bin/bash
TASKS=$(kubectl get task -n kubeteach-system --no-headers | cut -d " " -f1)
for t in $TASKS
do
    kubectl apply -f https://raw.githubusercontent.com/dergeberl/kubeteach-charts/main/solutions/exerciseset1/$t.yaml
    TIMEOUT=15
    i=0
    while [ $i -lt $TIMEOUT ]; do 
        i=$[$i+1];
        kubectl get task $t -n kubeteach-system -o jsonpath='{.status.state}' | grep successful
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

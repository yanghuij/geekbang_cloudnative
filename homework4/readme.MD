
kubectl create ns yhjhttpserver 
kubectl create -f httpserver.yaml -n yhjhttpserver


1. https
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=yhj Inc./CN=*.yhj.io' -keyout yhj.io.key -out yhj.io.crt
kubectl create -n istio-system secret tls yhj-credential --key=yhj.io.key --cert=yhj.io.crt
kubectl apply -f istio-specs-https.yaml -n yhjhttpserver

kubectl get svc -nistio-system
#istio-ingressgateway   LoadBalancer   $INGRESS_IP

curl --resolve httpsserver.yhj.io:443:$INGRESS_IP https://httpsserver.yhj.io/healthz -v -k


2. L7
kubectl apply -f istio-specs-l7.yaml -n yhjhttpserver

kubectl get svc -nistio-system
#istio-ingressgateway   LoadBalancer   $INGRESS_IP

curl -H "Host: httpserver.yhj.io" $INGRESS_IP/yhj/healthz
curl -H "Host: httpserver.yhj.io" $INGRESS_IP/

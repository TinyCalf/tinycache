ps aux | grep '[t]inycache.test' | awk '{print $2}' | xargs kill
rm nohup.out
nohup go test -run ^TestServer0$ & 
nohup go test -run ^TestServer1$ &
nohup go test -run ^TestServer2$ &